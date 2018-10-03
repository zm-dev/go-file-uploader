package go_file_uploader

import (
	"io"
	"context"
	"errors"
	"time"
	"net/url"
)

type File interface {
	io.Reader
	io.Seeker
}

type FileHeader struct {
	Filename string
	Size     int64
	File     File
}

type Uploader interface {
	Upload(fh FileHeader) error
	PresignedGetObject(hashValue string, expires time.Duration, reqParams url.Values) (u *url.URL, err error)
}

func Upload(ctx context.Context, fh FileHeader) error {
	u, ok := FromContext(ctx)
	if !ok {
		return errors.New("uploader不存在")
	}
	return u.Upload(fh)
}
