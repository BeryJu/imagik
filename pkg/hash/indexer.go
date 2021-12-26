package hash

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"beryju.org/imagik/pkg/config"
	"github.com/cornelk/hashmap"

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

func (hm *HashMap) Populated() bool {
	return hm.hashMap.Len() > 0
}

// RunIndexer Run full indexing
func (hm *HashMap) RunIndexer() {
	dir := config.C.RootDir
	hm.logger.WithField("dir", dir).Info("Started indexing...")
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hm.hashFile(path, info, err)
		return nil
	})
	hm.logger.WithField("hashes", hm.hashMap.Len()).Info("Finished indexing...")
}

func (hm *HashMap) Get(hash string) (string, bool) {
	val, exists := hm.hashMap.Get(hash)
	if val == nil {
		return "", exists
	}
	return val.(string), exists
}

func (hm *HashMap) UpdateSingle(path string) *FileHash {
	stat, err := os.Stat(path)
	return hm.hashFile(path, stat, err)
}

func (hm *HashMap) hashFile(p string, info os.FileInfo, err error) *FileHash {
	if err != nil {
		hm.logger.Warning(err)
	}

	if info.IsDir() {
		return nil
	}

	hashes, err := HashesForFile(p)
	if err != nil {
		// Don't return the error to not stop the walking
		hm.logger.Warning(err)
	}
	base := path.Base(p)
	ext := path.Ext(base)
	filename := strings.Replace(base, ext, "", 1)

	hm.writeM.Lock()
	defer hm.writeM.Unlock()
	hm.hashMap.Set(filename, p)
	hm.hashMap.Set(hashes.MD5, p)
	hm.hashMap.Set(hashes.SHA128, p)
	hm.hashMap.Set(hashes.SHA256, p)
	hm.hashMap.Set(hashes.SHA512, p)
	hm.hashMap.Set(hashes.SHA512Short, p)
	return hashes
}
