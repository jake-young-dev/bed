package file

import (
	"context"
	"fmt"
	"time"

	mn "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

/*
A simple Minio handler for uploading and deleting files
*/

const (
	CONTENT_TYPE = "application/tar+gzip"
)

type Cloud struct {
	cli        *mn.Client
	bucketName string
}

type ICloud interface {
	Upload(name string) error
	Delete(name string) error
}

// create a new minio file handler
func NewCloudHandler(url, id, key, bucket string) (*Cloud, error) {
	mcli, err := mn.New(url, &mn.Options{
		Creds: credentials.NewStaticV4(id, key, ""),
	})
	if err != nil {
		return nil, err
	}

	return &Cloud{
		cli:        mcli,
		bucketName: bucket,
	}, nil
}

// upload file to the "/data/" directory, minio bucket name is supplied during the struct creation and all file
// uploads have a 5 minute timeout
func (c *Cloud) Upload(name string) error {
	path := fmt.Sprintf("/data/%s", name)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err := c.cli.FPutObject(ctx, c.bucketName, name, path, mn.PutObjectOptions{ContentType: CONTENT_TYPE})
	if err != nil {
		return err
	}

	return nil
}

// delete file from minio, minio bucket name is supplied during the struct creation and all file
// uploads have a 5 minute timeout
func (c *Cloud) Delete(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	return c.cli.RemoveObject(ctx, c.bucketName, name, mn.RemoveObjectOptions{})
}
