package bucket

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

const (
	AWSProvider BucketType = iota
)

func New(bt BucketType, cfg any) (*Bucket, error) {
	b := new(Bucket)
	rt := reflect.TypeOf(cfg)
	switch bt {
	case AWSProvider:
		if rt != reflect.TypeOf(AWSConfig{}) {
			return nil, fmt.Errorf("invalid config type")
		}

		b.provide = newAWSSession(cfg.(AWSConfig))
		return b, nil
	default:
		return nil, fmt.Errorf("invalid bucket type")
	}
}

type BucketType int

type BucketInterface interface {
	Upload(file io.Reader, key string) error
	Download(source string, destine string) (*os.File, error)
	Delete(source string) error
}

type Bucket struct {
	provide BucketInterface
}

func (b *Bucket) Upload(file io.Reader, key string) error {
	return b.provide.Upload(file, key)
}

func (b *Bucket) Download(source string, destine string) (*os.File, error) {
	return b.provide.Download(source, destine)
}

func (b *Bucket) Delete(source string) error {
	return b.provide.Delete(source)
}
