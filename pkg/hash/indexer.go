package hash

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/BeryJu/imagik/pkg/config"
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
	hm.logger.WithField("dir", dir).Debug("Started indexing...")
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return hm.walk(path, info, err)
	})
	hm.logger.WithField("hashes", hm.hashMap.Len()).Debug("Finished indexing...")
}

func (hm *HashMap) Get(hash string) (string, bool) {
	val, exists := hm.hashMap.Get(hash)
	if val == nil {
		return "", exists
	}
	return val.(string), exists
}

func (hm *HashMap) UpdateSingle(path string) error {
	stat, err := os.Stat(path)
	return hm.walk(path, stat, err)
}

func (hm *HashMap) walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	hashes, err := HashesForFile(path)
	if err != nil {
		// Don't return the error to not stop the walking
		hm.logger.Warning(err)
	}

	hm.writeM.Lock()
	hm.hashMap.Set(hashes.MD5, path)
	hm.hashMap.Set(hashes.SHA128, path)
	hm.hashMap.Set(hashes.SHA256, path)
	hm.hashMap.Set(hashes.SHA512, path)
	hm.hashMap.Set(hashes.SHA512Short, path)
	hm.writeM.Unlock()
	if err != nil {
		hm.logger.Warning(err)
	}
	return nil
}
