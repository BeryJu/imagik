package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"beryju.org/imagik/pkg/config"
	"beryju.org/imagik/pkg/schema"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/getsentry/sentry-go"
)

type S3StorageDriver struct {
	s3         *s3.Client
	bucket     string
	log        *log.Entry
	preSigned  *s3.PresignClient
	downloader *manager.Downloader
	uploader   *manager.Uploader
}

func NewS3StorageDriver() (*S3StorageDriver, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           config.C.StorageS3Config.Endpoint,
			SigningRegion: region,
		}, nil
	})

	awsCfg, err := awsCfg.LoadDefaultConfig(
		context.Background(),
		awsCfg.WithRegion(config.C.StorageS3Config.Region),
		awsCfg.WithEndpointResolverWithOptions(customResolver),
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			config.C.StorageS3Config.AccessKey,
			config.C.StorageS3Config.SecretKey,
			"",
		)),
	)
	if err != nil {
		log.Warning("unable to load SDK config, %v", err)
		return nil, err
	}

	// Create S3 service client
	svc := s3.NewFromConfig(awsCfg)
	preSigned := s3.NewPresignClient(svc)
	downloader := manager.NewDownloader(svc)
	uploader := manager.NewUploader(svc)
	return &S3StorageDriver{
		s3:         svc,
		bucket:     config.C.StorageS3Config.Bucket,
		log:        log.WithField("component", "imagik.drivers.storage.s3"),
		preSigned:  preSigned,
		downloader: downloader,
		uploader:   uploader,
	}, nil
}

func (sd *S3StorageDriver) Walk(ctx context.Context, handler func(path string, info ObjectInfo)) error {
	p := s3.NewListObjectsV2Paginator(sd.s3, &s3.ListObjectsV2Input{
		Bucket: aws.String(sd.bucket),
	})

	var i int
	for p.HasMorePages() {
		i++

		page, err := p.NextPage(ctx)
		if err != nil {
			return err
		}

		// Log the objects found
		for _, obj := range page.Contents {
			tags, err := sd.s3.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
				Bucket: aws.String(sd.bucket),
				Key:    obj.Key,
			})
			if err != nil {
				sd.log.WithError(err).WithField("key", *obj.Key).Warning("failed to get tags for object")
				continue
			}
			tagsM := make(map[string]string, 0)
			for _, tag := range tags.TagSet {
				if tag.Key != nil && tag.Value != nil {
					tagsM[*tag.Key] = *tag.Value
				}
			}
			handler(*obj.Key, ObjectInfo{
				Tags: tagsM,
				ETag: *obj.ETag,
			})
		}
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
	_, err = sd.downloader.Download(ctx, f, &s3.GetObjectInput{
		Bucket: aws.String(sd.bucket),
		Key:    aws.String(key),
	})
	sd.log.WithField("key", key).Trace("object downloaded")
	if err != nil {
		return err
	}
	return nil
}

func (sd *S3StorageDriver) servePresign(rw http.ResponseWriter, r *http.Request, path string) {
	span := sentry.StartSpan(r.Context(), "imagik.server.serve_file")
	span.Description = path
	span.SetTag("imagik.path", path)
	defer span.Finish()

	req, err := sd.preSigned.PresignGetObject(span.Context(), &s3.GetObjectInput{
		Bucket: aws.String(sd.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		sd.log.WithError(err).Warning("failed to pre-sign request")
	}
	http.Redirect(rw, r, req.URL, http.StatusTemporaryRedirect)
}

func (sd *S3StorageDriver) needsHashUpdate(path string, info ObjectInfo) bool {
	fh := &FileHash{}
	sd.log.WithField("tags", info.Tags).Trace("tags")
	// Check if any tag is missing
	for key, _ := range fh.Map() {
		if _, ok := info.Tags[formatHashLabel(key)]; !ok {
			sd.log.WithField("key", path).WithField("tag", key).Trace("object needs updated tag")
			return true
		}
	}
	if etag, ok := info.Tags[formatHashLabel("ETag")]; !ok || etag != info.ETag {
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

	buffer := manager.NewWriteAtBuffer([]byte{})
	sd.log.WithField("key", path).Trace("[hash] downloading object")
	_, err := sd.downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(sd.bucket),
		Key:    aws.String(path),
	})
	sd.log.WithField("key", path).Trace("[hash] object downloaded")
	if err != nil {
		return nil, err
	}

	sha512hasher := sha512.New()
	sha256hasher := sha256.New()
	sha128hasher := sha1.New()
	md5hasher := md5.New()
	mw := io.MultiWriter(sha512hasher, sha256hasher, sha128hasher, md5hasher)

	if _, err := io.Copy(mw, bytes.NewReader(buffer.Bytes())); err != nil {
		sd.log.Warning(err)
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
	}

	tset := make([]types.Tag, 0)
	for key, value := range fh.Map() {
		tset = append(tset, types.Tag{
			Key:   aws.String(formatHashLabel(key)),
			Value: aws.String(value),
		})
	}
	sd.log.WithField("key", path).Trace("[hash] updating tags")
	_, err = sd.s3.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
		Bucket: aws.String(sd.bucket),
		Key:    aws.String(path),
		Tagging: &types.Tagging{
			TagSet: tset,
		},
	})
	if err != nil {
		sd.log.WithError(err).WithField("key", path).Warning("failed to update tags for object")
	}
	return fh, err
}

func (sd *S3StorageDriver) Upload(ctx context.Context, src io.Reader, p string) (*FileHash, error) {
	res, err := sd.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(sd.bucket),
		Key:    aws.String(p),
		Body:   src,
	})
	if err != nil {
		return nil, err
	}
	sd.log.WithField("path", p).Debug("Successfully wrote file.")
	return sd.HashesForFile(p, ObjectInfo{
		Tags: make(map[string]string),
		ETag: *res.ETag,
	}, ctx)
}

func (sd *S3StorageDriver) List(ctx context.Context, offset string) ([]schema.ListFile, error) {
	if !strings.HasSuffix(offset, "/") {
		offset += "/"
	}
	p := s3.NewListObjectsV2Paginator(sd.s3, &s3.ListObjectsV2Input{
		Bucket: aws.String(sd.bucket),
		Prefix: aws.String(offset),
	})

	var i int
	files := make([]schema.ListFile, 0)
	for p.HasMorePages() {
		i++

		page, err := p.NextPage(ctx)
		if err != nil {
			return files, err
		}

		// Log the objects found
		for _, obj := range page.Contents {
			files = append(files, schema.ListFile{
				Name:     strings.ReplaceAll(*obj.Key, offset, ""),
				Type:     "file",
				FullPath: fmt.Sprintf("/%s", *obj.Key),
				Mime:     "-",
			})
		}
	}

	return files, nil
}

func (sd *S3StorageDriver) Rename(ctx context.Context, from string, to string) error {
	_, err := sd.s3.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(sd.bucket),
		CopySource: aws.String(from),
		Key:        aws.String(to),
	})
	if err != nil {
		sd.log.WithError(err).Warning("failed to copy object")
		return err
	}
	_, err = sd.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(sd.bucket),
		Key:    aws.String(from),
	})
	if err != nil {
		sd.log.WithError(err).Warning("failed to delete source object")
		return err
	}
	return nil
}
