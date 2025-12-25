package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Storage struct {
	client *s3.Client
	bucket string
}

type Object struct {
	Body        io.ReadCloser
	ContentType string
	Size        int64
}

type Storage interface {
	Get(ctx context.Context, key string) (*Object, error)
}

func NewR2Storage(client *s3.Client, bucket string) *R2Storage {
	return &R2Storage{
		client: client,
		bucket: bucket,
	}
}

func (r *R2Storage) UploadSignature(ctx context.Context, estimateID int, body io.Reader, contentType string) error {
	key := fmt.Sprintf("signatures/%d.png", estimateID)

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})

	return err

}

func (r *R2Storage) Get(ctx context.Context, key string) (*Object, error) {

	out, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return &Object{
		Body:        out.Body,
		ContentType: aws.ToString(out.ContentType),
		Size:        *out.ContentLength,
	}, nil

}
