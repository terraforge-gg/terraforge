package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ObjectStoreService interface {
	GeneratePresignedPutUrl(ctx context.Context, key string, contentType string, fileSize int64) (string, error)
	GetFileMetadate(ctx context.Context, key string) (*metadata, error)
	MoveFile(ctx context.Context, sourceKey string, destinationKey string) (string, error)
}

type objectStoreService struct {
	client           *s3.Client
	assetsBucketName string
}

func NewObjectStoreService(client *s3.Client, assetsBucketName string) ObjectStoreService {
	return &objectStoreService{
		client:           client,
		assetsBucketName: assetsBucketName,
	}
}

func (s *objectStoreService) GeneratePresignedPutUrl(ctx context.Context, key string, contentType string, fileSize int64) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	put, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.assetsBucketName),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fileSize),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	return put.URL, err
}

type metadata struct {
	ContentLength int64
	ETag          string
}

var (
	ErrFileNotFound     = errors.New("s3 file not found")
	ErrFailedToMoveFile = errors.New("failed to move file")
)

func (s *objectStoreService) GetFileMetadate(ctx context.Context, key string) (*metadata, error) {
	response, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.assetsBucketName,
		Key:    &key,
	})

	if err != nil {
		return nil, ErrFileNotFound
	}

	etag, err := strconv.Unquote(*response.ETag)

	if err != nil {
		return nil, err
	}

	return &metadata{
		ContentLength: *response.ContentLength,
		ETag:          etag,
	}, nil
}

func (s *objectStoreService) MoveFile(ctx context.Context, sourceKey string, destinationKey string) (string, error) {
	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     &s.assetsBucketName,
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.assetsBucketName, sourceKey)),
		Key:        &destinationKey,
	})

	if err != nil {
		return "", ErrFailedToMoveFile
	}

	return "/" + s.assetsBucketName + "/" + destinationKey, nil
}
