package hash

import (
	"context"
	"path"
	"strings"
	"sync"

	"beryju.io/imagik/pkg/drivers/storage"
	"github.com/cornelk/hashmap"
	"github.com/getsentry/sentry-go"

	log "github.com/sirupsen/logrus"
)

type HashMap struct {
	logger  *log.Entry
	hashMap *hashmap.Map[string, string]
	writeM  sync.Mutex
	ctx     context.Context
	sd      storage.Driver
}

func New() *HashMap {
	m := &HashMap{
		logger:  log.WithField("component", "imagik.hash-map"),
		hashMap: hashmap.New[string, string](),
		writeM:  sync.Mutex{},
		ctx:     context.Background(),
		sd:      storage.FromConfig(),
	}
	return m
}

func (hm *HashMap) Populated() bool {
	return hm.hashMap.Len() > 0
}

// RunIndexer Run full indexing
func (hm *HashMap) RunIndexer() {
	hm.logger.Info("Started indexing...")
	span := sentry.StartSpan(hm.ctx, "imagik.hash.indexer", sentry.WithTransactionName("File Hasher"))
	defer span.Finish()
	err := hm.sd.Walk(span.Context(), func(path string, info storage.ObjectInfo) {
		hm.hashFile(path, info, span.Context())
	})
	if err != nil {
		hm.logger.WithError(err).Warning("failed to walk storage")
	}
	hm.logger.WithField("hashes", hm.hashMap.Len()).Info("Finished indexing...")
}

func (hm *HashMap) Get(hash string, ctx context.Context) (string, bool) {
	span := sentry.StartSpan(ctx, "imagik.hashmap.lookup")
	span.Description = hash
	span.SetTag("imagik.hash", hash)
	defer span.Finish()
	val, exists := hm.hashMap.Get(hash)
	if !exists {
		return "", exists
	}
	return val, exists
}

func (hm *HashMap) UpdateSingle(path string) *storage.FileHash {
	return hm.hashFile(path, storage.ObjectInfo{}, hm.ctx)
}

func (hm *HashMap) hashFile(p string, info storage.ObjectInfo, ctx context.Context) *storage.FileHash {
	hashes, err := hm.sd.HashesForFile(p, info, ctx)
	if err != nil {
		// Don't return the error to not stop the walking
		hm.logger.Warning(err)
		return nil
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
