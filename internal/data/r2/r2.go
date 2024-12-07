package r2

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client interface {
	Put(ctx context.Context, bucket, objectKey string, content io.ReadCloser, options func(*s3.PutObjectInput)) error
	Delete(ctx context.Context, bucket, objectKey string) error
}

type client struct {
	s3      *s3.Client
	presign *s3.PresignClient
}

func NewClient(accountID string, accessKeyID string, accessKeySecret string) Client {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})
	presignClient := s3.NewPresignClient(s3Client)

	return &client{
		s3:      s3Client,
		presign: presignClient,
	}
}

func (c *client) Put(
	ctx context.Context,
	bucket, objectKey string,
	content io.ReadCloser,
	options func(*s3.PutObjectInput),
) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
		Body:   content,
	}

	if options != nil {
		options(input)
	}

	_, err := c.s3.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Delete(ctx context.Context, bucket, objectKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	}

	_, err := c.s3.DeleteObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
