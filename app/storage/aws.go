package storage

import (
	"bytes"
	"encoding/hex"
	"log"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/zhikiri/blog-sync-cli/app/config"
)

type AWS struct {
	connection *s3.S3
	settings   *config.Settings
}

func NewAWS(settings config.Settings, creds *credentials.Credentials) (AWS, error) {
	awsConfig := aws.NewConfig().WithRegion(settings.Region)
	if creds != nil {
		awsConfig = awsConfig.WithCredentials(creds)
	}
	storage := AWS{
		connection: s3.New(session.New(awsConfig)),
		settings:   &settings,
	}
	return storage, nil
}

func NewAWSAuth(settings config.Settings, access, secret string) (AWS, error) {
	awsCreds := credentials.NewStaticCredentials(access, secret, "")
	if _, err := awsCreds.Get(); err != nil {
		return AWS{}, err
	}
	return NewAWS(settings, awsCreds)
}

func (aws AWS) GetFiles() ([]File, error) {
	bucket := aws.settings.Bucket
	objects, err := aws.connection.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: &bucket})
	if err != nil {
		return make([]File, 0), err
	}

	var checksum []byte
	var eTag string
	files := make([]File, len(objects.Contents))
	for _, data := range objects.Contents {
		eTag = *data.ETag
		checksum = make([]byte, 16) // *data.ETag
		if _, err := hex.Decode(checksum, bytes.TrimSpace([]byte(eTag[1:len(eTag)-1]))); err != nil {
			return make([]File, 0), err
		}
		files = append(files, File{
			Path:     path.Join("/", *data.Key),
			Checksum: checksum,
		})
	}
	return files, nil
}

func (a AWS) PutFile(file File, path string) error {
	if err := file.loadFileDataFrom(path); err != nil {
		return err
	}
	if file.Size == 0 {
		// Prevent empty file upload
		return nil
	}

	start := time.Now()
	_, err := a.connection.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(a.settings.Bucket),
		Key:           aws.String(file.Path),
		ACL:           aws.String("public-read"),
		Body:          file.Body,
		ContentLength: aws.Int64(file.Size),
		ContentType:   aws.String(file.Type),
	})
	if err != nil {
		return err
	}
	stop := time.Now()
	log.Printf(
		"File %s was uploaded, (%.2f Kb) was sent to S3, it took %v\n",
		file.Path,
		float64(file.Size)/1000.0,
		stop.Sub(start),
	)

	return nil
}

func (a AWS) DelFile(file File) error {
	_, err := a.connection.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(a.settings.Bucket),
		Key:    aws.String(file.Path),
	})
	if err == nil {
		log.Printf("File %s was deleted from S3 bucket\n", file.Path)
	}
	return err
}
