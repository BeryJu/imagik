package storage

import (
	"context"
	"io"
	"net/http"
	"regexp"

	"beryju.io/imagik/pkg/config"
	"beryju.io/imagik/pkg/schema"
)

func formatHashLabel(val string) string {
	return val
}

type ObjectInfo struct {
	Tags map[string]string
	ETag string
}

func (oi *ObjectInfo) Hash() *FileHash {
	fh := &FileHash{
		SHA128:      oi.Tags[formatHashLabel("SHA128")],
		SHA256:      oi.Tags[formatHashLabel("SHA256")],
		SHA512:      oi.Tags[formatHashLabel("SHA512")],
		SHA512Short: oi.Tags[formatHashLabel("SHA512Short")],
		MD5:         oi.Tags[formatHashLabel("MD5")],
		Mime:        oi.Tags[formatHashLabel("Mime")],
		ETag:        oi.ETag,
	}
	return fh
}

type FileHash struct {
	SHA128      string
	SHA256      string
	SHA512      string
	SHA512Short string
	MD5         string
	ETag        string
	Mime        string
}

var escapedChars = regexp.MustCompile(`(\s|;|=)`)

func (fh *FileHash) Map() map[string]string {
	m := make(map[string]string, 5)
	m[formatHashLabel("SHA128")] = fh.SHA128
	m[formatHashLabel("SHA256")] = fh.SHA256
	m[formatHashLabel("SHA512")] = fh.SHA512
	m[formatHashLabel("SHA512Short")] = fh.SHA512Short
	m[formatHashLabel("MD5")] = fh.MD5
	m[formatHashLabel("ETag")] = fh.ETag
	m[formatHashLabel("Mime")] = escapedChars.ReplaceAllString(fh.Mime, "_")
	return m
}

type Driver interface {
	Walk(context.Context, func(path string, info ObjectInfo)) error
	HashesForFile(path string, info ObjectInfo, ctx context.Context) (*FileHash, error)
	Serve(w http.ResponseWriter, r *http.Request, path string)
	Upload(ctx context.Context, src io.Reader, p string) (*FileHash, error)
	List(ctx context.Context, offset string) ([]schema.ListFile, error)
	Rename(ctx context.Context, from string, to string) error
	CleanURL(raw string) string
	Stat(path string, ctx context.Context) *schema.MetaResponse
}

func FromConfig() Driver {
	switch config.C.StorageDriver {
	case "local":
		return NewLocalStorageDriver()
	case "s3":
		sd, err := NewS3StorageDriver()
		if err != nil {
			panic(err)
		}
		return sd
	default:
		panic("invalid storage driver")
	}
}
