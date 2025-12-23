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
