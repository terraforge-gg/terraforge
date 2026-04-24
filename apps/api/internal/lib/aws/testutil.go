package aws

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

type TestLocalStack struct {
	Container   *localstack.LocalStackContainer
	DatabaseUrl string
}

func NewTestLocalStackS3Client(t *testing.T, bucketName string) (*s3.Client, string, error) {
	t.Helper()
	ctx := context.Background()

	localstackContainer, err := localstack.Run(ctx, "localstack/localstack:1.4.0")
	if err != nil {
		return nil, "", fmt.Errorf("failed to start container: %w", err)
	}

	host, err := localstackContainer.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get host: %w", err)
	}

	port, err := localstackContainer.MappedPort(ctx, "4566/tcp")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get port: %w", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"test",
			"test",
			"",
		)),
		config.WithBaseEndpoint(endpoint),
	)

	if err != nil {
		return nil, "", fmt.Errorf("failed to create aws config: %w", err)
	}

	localStackClient := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	if localStackClient == nil {
		return nil, "", fmt.Errorf("failed to initialize localstack client")
	}

	setupTestBucket(t, context.Background(), localStackClient, bucketName)
	return localStackClient, endpoint, nil
}

func setupTestBucket(t *testing.T, ctx context.Context, s3Client *s3.Client, bucketName string) error {
	t.Helper()
	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		if !strings.Contains(err.Error(), "BucketAlreadyOwned") {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}
