package bucket

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSConfig struct {
	Config         *aws.Config
	BucketDownload string
	BucketUpload   string
}

type awsSession struct {
	sess           *session.Session
	bucketDownload string
	bucketUpload   string
}

func newAWSSession(cfg AWSConfig) *awsSession {
	sess, err := session.NewSession(cfg.Config)
	if err != nil {
		return nil
	}

	return &awsSession{
		sess:           sess,
		bucketDownload: cfg.BucketDownload,
		bucketUpload:   cfg.BucketUpload,
	}
}

func (a *awsSession) Download(source string, destine string) (*os.File, error) {
	file, err := os.Create(destine)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(a.sess)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(a.bucketDownload),
			Key:    aws.String(source),
		})

	return file, nil
}

func (a *awsSession) Upload(file io.Reader, key string) error {
	uploader := s3manager.NewUploader(a.sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.bucketUpload),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return err
	}
	return nil
}

func (a *awsSession) Delete(source string) error {
	svc := s3.New(a.sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(a.bucketDownload),
		Key:    aws.String(source),
	})
	if err != nil {
		return err
	}
	return svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(a.bucketDownload),
		Key:    aws.String(source),
	})
}
