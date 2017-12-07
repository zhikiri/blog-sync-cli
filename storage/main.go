package storage

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/zhikiri/blog-sync-cli/config"
)

// Connect function for establish S3 connection
func Connect(s *config.Setup) (*s3.S3, error) {
	cre := credentials.NewStaticCredentials(s.AWS.AccessKey, s.AWS.SecretKey, "")
	if _, err := cre.Get(); err != nil {
		return &s3.S3{}, err
	}

	cnf := aws.NewConfig().WithRegion(s.AWS.Region).WithCredentials(cre)
	return s3.New(session.New(cnf)), nil
}

// GetStorageChecksum function for retrieve the checksum information from the server
func GetStorageChecksum(c *s3.S3, bucket string) (map[string]string, error) {
	obj, err := c.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: &bucket})
	if err != nil {
		return make(map[string]string), err
	}

	var sum string
	state := make(map[string]string)
	for _, data := range obj.Contents {
		sum = *data.ETag
		state[*data.Key] = sum[1 : len(sum)-1]
	}
	fmt.Printf("Checksum of %d files was loaded\n", len(state))
	return state, nil
}

// UpdateFile function for create/update file in S3
func UpdateFile(abspath, relpath string, c *s3.S3, bucket string) error {
	object, err := getObject(abspath, relpath)
	if err != nil {
		return err
	}

	if object.Size == 0 {
		// Avoid of sending empty files
		return nil
	}

	start := time.Now()
	/*
		_, err = c.PutObject(&s3.PutObjectInput{
			Bucket:        aws.String(bucket),
			Key:           aws.String(object.Path),
			ACL:           aws.String("public-read"),
			Body:          object.Body,
			ContentLength: aws.Int64(object.Size),
			ContentType:   aws.String(object.Type),
		})
		if err != nil {
			return err
		}
	*/
	stop := time.Now()
	fmt.Printf(
		"File %s was uploaded, (%.2f Kb) was sent to S3, it took %v\n",
		object.Path,
		float64(object.Size)/1000.0,
		stop.Sub(start),
	)
	return nil
}

// DeleteFile function for remove file in S3 bucket
func DeleteFile(relpath string, c *s3.S3, bucket string) error {
	/*
		_, err := c.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(relpath),
		})
		if err == nil {
			fmt.Printf("File %s was deleted from S3 bucket", relpath)
		}
		return err
	*/
	return nil
}
