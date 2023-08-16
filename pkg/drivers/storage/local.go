package storage

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"

	"beryju.io/imagik/pkg/config"
	"beryju.io/imagik/pkg/schema"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

type LocalStorageDriver struct {
	root string
	log  *log.Entry
}

func NewLocalStorageDriver() *LocalStorageDriver {
	return &LocalStorageDriver{
		root: config.C.StorageLocalConfig.Root,
		log:  log.WithField("component", "imagik.drivers.storage.local"),
	}
}

func (lsd *LocalStorageDriver) Walk(ctx context.Context, walker func(path string, info ObjectInfo)) error {
	return filepath.Walk(lsd.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		walker(path, ObjectInfo{})
		return nil
	})
}

func (lsd *LocalStorageDriver) Serve(rw http.ResponseWriter, r *http.Request, path string) {
	span := sentry.StartSpan(r.Context(), "imagik.server.serve_file")
	span.Description = path
	span.SetTag("imagik.path", path)
	defer span.Finish()
	http.ServeFile(rw, r, path)
}

func (lsd *LocalStorageDriver) HashesForFile(path string, info ObjectInfo, ctx context.Context) (*FileHash, error) {
	span := sentry.StartSpan(ctx, "imagik.hash.file_hashes")
	span.Description = path
	span.SetTag("imagik.path", path)
	defer span.Finish()
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to stat path")
	}
	if stat.IsDir() {
		return nil, errors.New("path is directory")
	}

	f, err := os.Open(path)
	if err != nil {
		lsd.log.Warning(err)
		return nil, err
	}
	defer f.Close()
	sha512hasher := sha512.New()
	sha256hasher := sha256.New()
	sha128hasher := sha1.New()
	md5hasher := md5.New()
	mw := io.MultiWriter(sha512hasher, sha256hasher, sha128hasher, md5hasher)
	mime, err := mimetype.DetectReader(f)
	if err != nil {
		lsd.log.WithError(err).Warning("failed to detect mime type")
		return nil, err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		lsd.log.WithError(err).Warning("failed to re-read file")
		return nil, err
	}

	if _, err := io.Copy(mw, f); err != nil {
		lsd.log.WithError(err).Warning("failed to stream to hasher")
		return nil, err
	}
	sha512sum := hex.EncodeToString(sha512hasher.Sum(nil))
	return &FileHash{
		SHA128:      hex.EncodeToString(sha128hasher.Sum(nil)),
		SHA256:      hex.EncodeToString(sha256hasher.Sum(nil)),
		SHA512:      sha512sum,
		SHA512Short: sha512sum[:16],
		MD5:         hex.EncodeToString(md5hasher.Sum(nil)),
		Mime:        mime.String(),
	}, nil
}

func (lsd *LocalStorageDriver) Upload(ctx context.Context, src io.Reader, p string) (*FileHash, error) {
	filePath := lsd.CleanURL(p)
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		lsd.log.Warning(err)
		return nil, err
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		lsd.log.Warning(err)
		return nil, err
	}
	n, err := io.Copy(f, src)
	if err != nil {
		lsd.log.Warning(err)
		return nil, err
	}
	lsd.log.WithField("n", n).WithField("path", filePath).Debug("Successfully wrote file.")
	return lsd.HashesForFile(filePath, ObjectInfo{}, ctx)
}

func getElementsForDirectory(path string) int {
	files, err := os.ReadDir(path)
	if err != nil {
		return 0
	}
	return len(files)
}

func (lsd *LocalStorageDriver) List(ctx context.Context, offset string) ([]schema.ListFile, error) {
	fullDir := lsd.CleanURL(offset)
	files, err := os.ReadDir(fullDir)
	if err != nil {
		return nil, err
	}
	lfiles := make([]schema.ListFile, 0)
	for _, f := range files {
		fullName := path.Join(fullDir, f.Name())
		file := schema.ListFile{
			Name:     f.Name(),
			FullPath: filepath.Join(filepath.FromSlash(path.Clean("/"+offset)), f.Name()),
		}
		if f.IsDir() {
			file.Type = "directory"
			file.ChildElements = getElementsForDirectory(fullName)
		} else {
			file.Type = "file"
			mime, err := mimetype.DetectFile(fullName)
			if err == nil {
				file.Mime = mime.String()
			}
		}
		lfiles = append(lfiles, file)
	}
	return lfiles, nil
}

func (lsd *LocalStorageDriver) Rename(ctx context.Context, from string, to string) error {
	if _, err := os.Stat(from); err != nil {
		return err
	}
	err := os.Rename(from, to)
	if err != nil {
		return err
	}
	return nil
}

func (lsd *LocalStorageDriver) CleanURL(raw string) string {
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	raw = strings.TrimPrefix(raw, config.C.StorageLocalConfig.Root)
	return filepath.Join(config.C.StorageLocalConfig.Root, filepath.FromSlash(path.Clean("/"+raw)))
}

func (lsd *LocalStorageDriver) Stat(path string, ctx context.Context) *schema.MetaResponse {
	fullPath := lsd.CleanURL(path)
	// Get stat for common file stats
	stat, err := os.Stat(fullPath)
	if err != nil {
		log.Warn(err)
	}
	// Get hashes for linking
	hashes, err := lsd.HashesForFile(fullPath, ObjectInfo{}, ctx)
	if err != nil {
		lsd.log.WithError(err).Warning("failed to get hashes for file")
		return nil
	}
	response := schema.MetaResponse{
		GenericResponse: schema.GenericResponse{
			Successful: true,
		},
		Mime:   hashes.Mime,
		Hashes: hashes.Map(),
	}
	if stat != nil {
		response.Name = stat.Name()
		response.CreationDate = stat.ModTime()
		response.Size = stat.Size()
	}
	return &response
}
