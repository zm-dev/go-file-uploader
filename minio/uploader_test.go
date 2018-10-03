package minio

import (
	"testing"
	"github.com/minio/minio-go"
	"log"
	. "github.com/zm-dev/go-file-uploader"
	"os"
)

var uploader Uploader

func TestMain(m *testing.M) {

	minioClient, err := minio.New(
		"59.111.58.150:9000",
		"zm2018",
		"zhiming2018",
		false,
	)

	if err != nil {
		log.Fatalf("minio client 创建失败! error: %+v", err)
	}
	uploader = NewMinioUploader(HashFunc(MD5HashFunc), minioClient, "test", Hash2StorageNameFunc(TwoCharsPrefixHash2StorageNameFunc))
	m.Run()
}

func TestMinioUploader_Upload(t *testing.T) {
	filename := "/Users/taoyu/Desktop/2-4 章节总结.mp4"
	fi, err := os.Stat(filename)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	err = uploader.Upload(FileHeader{file.Name(), fi.Size(), file})
	if err != nil {
		log.Fatalln(err)
	}
}
