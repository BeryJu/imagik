package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/cornelk/hashmap"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type HashMap struct {
	logger  *log.Entry
	hashMap hashmap.HashMap
	writeM  sync.Mutex
}

func New() *HashMap {
	m := &HashMap{
		logger:  log.WithField("component", "hash-map"),
		hashMap: hashmap.HashMap{},
		writeM:  sync.Mutex{},
	}
	return m
}

// RunIndexer Run full indexing
func (hm *HashMap) RunIndexer() {
	hm.logger.Debug("Started indexing...")
	filepath.Walk(viper.GetString(config.ConfigRootDir), func(path string, info os.FileInfo, err error) error {
		return hm.walk(path, info, err)
	})
	hm.logger.WithField("hashes", hm.hashMap.Len()).Debug("Finished indexing...")
}

func (hm *HashMap) Get(hash string) (string, bool) {
	val, exists := hm.hashMap.Get(hash)
	if val == nil {
		return "", exists
	} else {
		return val.(string), exists
	}
}

func (hm *HashMap) UpdateSingle(path string) error {
	stat, err := os.Stat(path)
	return hm.walk(path, stat, err)
}

func (hm *HashMap) walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		log.Warning(err)
	}
	defer f.Close()
	sha256hasher := sha256.New()
	sha1hasher := sha1.New()
	md5hasher := md5.New()
	mw := io.MultiWriter(sha256hasher, sha1hasher, md5hasher)

	if _, err := io.Copy(mw, f); err != nil {
		log.Warning(err)
	}
	sha256sum := hex.EncodeToString(sha256hasher.Sum(nil))
	hm.writeM.Lock()
	defer hm.writeM.Unlock()
	hm.hashMap.Set(sha256sum, path)
	hm.hashMap.Set(sha256sum[:16], path)
	hm.hashMap.Set(hex.EncodeToString(sha1hasher.Sum(nil)), path)
	hm.hashMap.Set(hex.EncodeToString(md5hasher.Sum(nil)), path)
	if err != nil {
		log.Warning(err)
	}
	return nil
}
