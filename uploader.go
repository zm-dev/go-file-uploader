package go_file_uploader

import (
	"io"
	"context"
	"errors"
	"time"
	"net/url"
)

type FileHeader struct {
	Filename string
	Size     int64
	File     io.ReadSeeker
}

type Uploader interface {
	Upload(fh FileHeader, extra string) (f *FileModel, err error)
	PresignedGetObject(hashValue string, expires time.Duration, reqParams url.Values) (u *url.URL, err error)
	ReadFile(hashValue string) (rf ReadFile, size int64, err error)
	Store() Store
}

type ReadFile interface {
	io.Reader
	io.Closer
	io.Seeker
	io.ReaderAt
}

func Upload(ctx context.Context, fh FileHeader, extra string) (f *FileModel, err error) {
	u, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("uploader不存在")
	}
	return u.Upload(fh, extra)
}
