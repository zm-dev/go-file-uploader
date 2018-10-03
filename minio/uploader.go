package minio

import (
	"github.com/minio/minio-go"
	"mime"
	"io"
	"path/filepath"
	. "github.com/zm-dev/go-file-uploader"
	"time"
	"net/url"
)

type minioUploader struct {
	h           Hasher
	minioClient *minio.Client
	bucketName  string
	h2sn        Hash2StorageName
}

func (mu *minioUploader) saveToMinio(hashValue string, fh FileHeader) error {
	name, err := mu.h2sn.Convent(hashValue)
	if err != nil {
		return err
	}
	_, err = fh.File.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	obj, _ := mu.minioClient.GetObject(mu.bucketName, name, minio.GetObjectOptions{})
	_, err = obj.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code != "NoSuchKey" {
			return err
		}
	} else {
		// 文件已经存在
		return nil
	}
	ext := filepath.Ext(fh.Filename)
	// 在 apline 镜像中 mime.TypeByExtension 只能用 jpg
	if ext == "jpeg" {
		ext = "jpg"
	}

	_, err = mu.minioClient.PutObject(
		mu.bucketName,
		name,
		fh.File,
		fh.Size,
		minio.PutObjectOptions{ContentType: mime.TypeByExtension(ext)},
	)
	return err
}

func (mu *minioUploader) Upload(fh FileHeader) error {

	hashValue, err := mu.h.Hash(fh.File)
	if err != nil {
		return err
	}

	if err := mu.saveToMinio(hashValue, fh); err != nil {
		return err
	}

	return nil
}

func (mu *minioUploader) PresignedGetObject(hashValue string, expires time.Duration, reqParams url.Values) (u *url.URL, err error) {
	name, err := mu.h2sn.Convent(hashValue)
	if err != nil {
		return nil, err
	}
	return mu.minioClient.PresignedGetObject(mu.bucketName, name, expires, reqParams)
}

func NewMinioUploader(h Hasher, minioClient *minio.Client, bucketName string, h2sn Hash2StorageName) Uploader {
	if h2sn == nil {
		h2sn = Hash2StorageNameFunc(DefaultHash2StorageNameFunc)
	}
	return &minioUploader{h: h, minioClient: minioClient, bucketName: bucketName, h2sn: h2sn}
}
