package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type FileHash struct {
	SHA128      string
	SHA256      string
	SHA512      string
	SHA512Short string
	MD5         string
}

func (fh *FileHash) Map() map[string]string {
	m := make(map[string]string, 5)
	m["SHA128"] = fh.SHA128
	m["SHA256"] = fh.SHA256
	m["SHA512"] = fh.SHA512
	m["SHA512Short"] = fh.SHA512Short
	m["MD5"] = fh.MD5
	return m
}

func HashesForFile(path string) (*FileHash, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to stat path")
	}
	if stat.IsDir() {
		return nil, errors.New("path is directory")
	}

	f, err := os.Open(path)
	if err != nil {
		log.Warning(err)
	}
	defer f.Close()
	sha512hasher := sha512.New()
	sha256hasher := sha256.New()
	sha128hasher := sha1.New()
	md5hasher := md5.New()
	mw := io.MultiWriter(sha512hasher, sha256hasher, sha128hasher, md5hasher)

	if _, err := io.Copy(mw, f); err != nil {
		log.Warning(err)
	}
	sha512sum := hex.EncodeToString(sha512hasher.Sum(nil))
	return &FileHash{
		SHA128:      hex.EncodeToString(sha128hasher.Sum(nil)),
		SHA256:      hex.EncodeToString(sha256hasher.Sum(nil)),
		SHA512:      sha512sum,
		SHA512Short: sha512sum[:16],
		MD5:         hex.EncodeToString(md5hasher.Sum(nil)),
	}, nil
}
