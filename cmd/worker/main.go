package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"worker/internal/bucket"
	"worker/internal/queue"
)

func main() {
	// TODO: rabbitmq config
	qcfg := queue.RabbitMQConfig{
		URL:       "amqp://guest:guest@localhost:5672/",
		TopicName: "test",
		Timeout:   time.Second * 30,
	}

	// TODO: create new queue
	qc, err := queue.New(queue.RabbitMQ, qcfg)

	if err != nil {
		panic(err)
	}

	// TODO: create channel to consume messages
	c := make(chan queue.QueueDto)
	qc.Consume(c)

	// TODO: bucket config
	bcfg := bucket.AWSConfig{
		Config: &aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				"access_key_id",
				"secret_access_key",
				"",
			),
		},
		BucketDownload: "bucket-download-raw",
		BucketUpload:   "bucket-upload-compact",
	}

	// TODO: create new bucket session
	b, err := bucket.New(bucket.AWSProvider, bcfg)
	if err != nil {
		panic(err)
	}

	for msg := range c {
		source := fmt.Sprint("%s/%s", msg.Path, msg.Filename)
		destine := fmt.Sprint("%d_%s", msg.ID, msg.Filename)

		file, err := b.Download(source, destine)
		if err != nil {
			log.Printf("error to download file: %s", err.Error())
			continue
		}

		body, err := io.ReadAll(file)
		if err != nil {
			log.Printf("error to read file: %s", err.Error())
			continue
		}

		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err = zw.Write(body)
		if err != nil {
			log.Printf("error to compress file: %s", err.Error())
			continue
		}

		if err := zw.Close(); err != nil {
			log.Printf("error to close compress file: %s", err.Error())
			continue
		}

		zr, err := gzip.NewReader(&buf)
		if err != nil {
			log.Printf("error to decompress file: %s", err.Error())
			continue
		}

		err = b.Upload(zr, source)
		if err != nil {
			log.Printf("error to upload file: %s", err.Error())
			continue
		}

		err = os.Remove(destine)
		if err != nil {
			log.Printf("error to remove file: %s", err.Error())
			continue
		}

	}

}
