package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestS3BucketCreated(t *testing.T) {
	bucketName := "test-bucket-name"

	s3Client, _, err := NewTestLocalStackS3Client(t, bucketName)
	if err != nil {
		t.Fatalf("failed to create S3 client: %v", err)
	}

	ctx := context.Background()
	resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		t.Fatalf("failed to list buckets: %v", err)
	}

	assert.NotZero(t, len(resp.Buckets))
	found := false

	for _, bucket := range resp.Buckets {
		if *bucket.Name == bucketName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected bucket %q to exist, but it was not found", bucketName)
	}
}
