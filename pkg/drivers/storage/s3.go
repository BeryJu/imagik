package storage

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
	log "github.com/sirupsen/logrus"

	"beryju.io/imagik/pkg/config"
	"beryju.io/imagik/pkg/schema"
	"github.com/getsentry/sentry-go"
)

type S3StorageDriver struct {
	s3     *minio.Client
	bucket string
	log    *log.Entry
}

func NewS3StorageDriver() (*S3StorageDriver, error) {
	endpoint, err := url.Parse(config.C.StorageS3Config.Endpoint)
	if err != nil {
		return nil, err
	}

	opts := &minio.Options{
		Secure: strings.EqualFold(endpoint.Scheme, "https"),
	}
	if config.C.StorageS3Config.AccessKey != "" {
		opts.Creds = credentials.NewStaticV4(
			config.C.StorageS3Config.AccessKey,
			config.C.StorageS3Config.SecretKey,
			"",
		)
	} else {
		opts.Creds = credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvAWS{},
				&credentials.EnvMinio{},
				&credentials.FileAWSCredentials{},
				&credentials.FileMinioClient{},
				&credentials.IAM{},
			},
		)
	}

	minioClient, err := minio.New(endpoint.Host, opts)
	if err != nil {
		log.Fatalln(err)
	}

	return &S3StorageDriver{
		s3:     minioClient,
		bucket: config.C.StorageS3Config.Bucket,
		log:    log.WithField("component", "imagik.drivers.storage.s3"),
	}, nil
}

func (sd *S3StorageDriver) getTagsMap(ctx context.Context, key string) map[string]string {
	tagsM := make(map[string]string, 0)
	tags, err := sd.s3.GetObjectTagging(ctx, sd.bucket, key, minio.GetObjectTaggingOptions{})
	if err != nil {
		sd.log.WithError(err).WithField("key", key).Warning("failed to get tags for object")
		return tagsM
	}
	for key, value := range tags.ToMap() {
		if key != "" && value != "" {
			tagsM[key] = value
		}
	}
	return tagsM
}

func (sd *S3StorageDriver) Walk(ctx context.Context, handler func(path string, info ObjectInfo)) error {
	objects := sd.s3.ListObjects(ctx, sd.bucket, minio.ListObjectsOptions{
		Recursive:    true,
		WithMetadata: true,
	})

	// Log the objects found
	for obj := range objects {
		handler(obj.Key, ObjectInfo{
			Tags: sd.getTagsMap(ctx, obj.Key),
			ETag: obj.ETag,
		})
	}

	return nil
}

func (sd *S3StorageDriver) Serve(rw http.ResponseWriter, r *http.Request, path string) {
	if config.C.StorageS3Config.UsePresign {
		sd.servePresign(rw, r, path)
		return
	}
	sd.serveDownload(rw, r, path)
}

func (sd *S3StorageDriver) serveDownload(rw http.ResponseWriter, r *http.Request, p string) {
	localPath := path.Join(os.TempDir(), "imagik/", p)
	// check if we already have the key locally
	_, err := os.Stat(localPath)
	if err != nil {
		err := sd.downloadFile(r.Context(), p, localPath)
		if err != nil {
			sd.log.WithError(err).Warning("failed to download file")
			return
		}
	}
	http.ServeFile(rw, r, localPath)
}

func (sd *S3StorageDriver) downloadFile(ctx context.Context, key string, localPath string) error {
	f, err := os.Create(localPath)
	if err != nil {
		sd.log.WithError(err).Warning("failed to open local file")
		return err
	}
	obj, err := sd.s3.GetObject(ctx, sd.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	n, err := io.Copy(f, obj)
	if err != nil {
		return err
	}
	sd.log.WithField("size_bytes", n).Trace("streamed object into file")
	return nil
}

func (sd *S3StorageDriver) servePresign(rw http.ResponseWriter, r *http.Request, path string) {
	span := sentry.StartSpan(r.Context(), "imagik.server.serve_file")
	span.Description = path
	span.SetTag("imagik.path", path)
	defer span.Finish()

	req, err := sd.s3.PresignedGetObject(span.Context(), sd.bucket, path, time.Hour*24, url.Values{})
	if err != nil {
		sd.log.WithError(err).Warning("failed to pre-sign request")
	}
	http.Redirect(rw, r, req.String(), http.StatusTemporaryRedirect)
}

func (sd *S3StorageDriver) needsHashUpdate(path string, info ObjectInfo) bool {
	fh := &FileHash{}
	sd.log.WithField("tags", info.Tags).Trace("tags")
	// Check if any tag is missing
	for key := range fh.Map() {
		if _, ok := info.Tags[key]; !ok {
			sd.log.WithField("key", path).WithField("tag", key).Trace("object needs updated tag")
			return true
		}
	}
	if etag, ok := info.Tags["ETag"]; !ok || etag != info.ETag {
		sd.log.WithField("key", path).WithField("etag", etag).Trace("object etag has changed")
		return true
	}
	return false
}

func (sd *S3StorageDriver) HashesForFile(path string, info ObjectInfo, ctx context.Context) (*FileHash, error) {
	span := sentry.StartSpan(ctx, "imagik.hash.file_hashes")
	span.Description = path
	span.SetTag("imagik.path", path)
	defer span.Finish()

	needsHashing := sd.needsHashUpdate(path, info)
	if !needsHashing {
		sd.log.WithField("key", path).Trace("object doesn't need hashing")
		return info.Hash(), nil
	}

	sd.log.WithField("key", path).Trace("[hash] downloading object")
	obj, err := sd.s3.GetObject(ctx, sd.bucket, path, minio.GetObjectOptions{})
	sd.log.WithField("key", path).Trace("[hash] object downloaded")
	if err != nil {
		return nil, err
	}

	sha512hasher := sha512.New()
	sha256hasher := sha256.New()
	sha128hasher := sha1.New()
	md5hasher := md5.New()
	mw := io.MultiWriter(sha512hasher, sha256hasher, sha128hasher, md5hasher)

	mime, err := mimetype.DetectReader(obj)
	if err != nil {
		sd.log.WithError(err).Warning("failed to detect mime type")
		return nil, err
	}
	obj.Seek(0, io.SeekStart)

	if _, err := io.Copy(mw, obj); err != nil {
		sd.log.WithError(err).Warning("failed to stream to hasher")
		return nil, err
	}
	sha512sum := hex.EncodeToString(sha512hasher.Sum(nil))
	fh := &FileHash{
		SHA128:      hex.EncodeToString(sha128hasher.Sum(nil)),
		SHA256:      hex.EncodeToString(sha256hasher.Sum(nil)),
		SHA512:      sha512sum,
		SHA512Short: sha512sum[:16],
		MD5:         hex.EncodeToString(md5hasher.Sum(nil)),
		ETag:        info.ETag,
		Mime:        mime.String(),
	}

	tset, err := tags.NewTags(fh.Map(), false)
	if err != nil {
		sd.log.WithError(err).WithField("key", path).Warning("failed to create tags")
		return nil, err
	}
	sd.log.WithField("key", path).Trace("[hash] updating tags")
	err = sd.s3.PutObjectTagging(ctx, sd.bucket, path, tset, minio.PutObjectTaggingOptions{})
	if err != nil {
		sd.log.WithError(err).WithField("key", path).Warning("failed to update tags for object")
	}
	return fh, err
}

func (sd *S3StorageDriver) Upload(ctx context.Context, src io.Reader, p string) (*FileHash, error) {
	res, err := sd.s3.PutObject(ctx, sd.bucket, p, src, -1, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}
	sd.log.WithField("path", p).Debug("Successfully wrote file.")
	return sd.HashesForFile(p, ObjectInfo{
		Tags: make(map[string]string),
		ETag: res.ETag,
	}, ctx)
}

func (sd *S3StorageDriver) List(ctx context.Context, offset string) ([]schema.ListFile, error) {
	if !strings.HasSuffix(offset, "/") {
		offset += "/"
	}
	objects := sd.s3.ListObjects(ctx, sd.bucket, minio.ListObjectsOptions{
		Prefix: offset,
	})
	files := make([]schema.ListFile, 0)

	for obj := range objects {
		tags := sd.getTagsMap(ctx, obj.Key)
		files = append(files, schema.ListFile{
			Name:     strings.ReplaceAll(obj.Key, offset, ""),
			Type:     "file",
			FullPath: fmt.Sprintf("/%s", obj.Key),
			Mime:     tags[formatHashLabel("Mime")],
		})
	}

	return files, nil
}

func (sd *S3StorageDriver) Rename(ctx context.Context, from string, to string) error {
	_, err := sd.s3.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: sd.bucket,
		Object: to,
	}, minio.CopySrcOptions{
		Bucket: sd.bucket,
		Object: from,
	})
	if err != nil {
		sd.log.WithError(err).Warning("failed to copy object")
		return err
	}
	err = sd.s3.RemoveObject(ctx, sd.bucket, from, minio.RemoveObjectOptions{})
	if err != nil {
		sd.log.WithError(err).Warning("failed to delete source object")
		return err
	}
	return nil
}

func (sd *S3StorageDriver) CleanURL(raw string) string {
	raw = strings.TrimPrefix(raw, "/")
	return raw
}

func (sd *S3StorageDriver) Stat(path string, ctx context.Context) *schema.MetaResponse {
	fullPath := sd.CleanURL(path)
	response := schema.MetaResponse{
		GenericResponse: schema.GenericResponse{
			Successful: true,
		},
	}
	// Get stat for common file stats
	stat, err := sd.s3.StatObject(ctx, sd.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		log.Warn(err)
	} else {
		response.Name = stat.Key
		response.CreationDate = stat.LastModified
		response.Size = stat.Size
	}
	// Get hashes for linking
	hashes, err := sd.HashesForFile(fullPath, ObjectInfo{}, ctx)
	if err != nil {
		sd.log.WithError(err).Warning("failed to get hashes for file")
		return nil
	}
	response.Mime = hashes.Mime
	response.Hashes = hashes.Map()
	return &response
}
