package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/terraforge-gg/terraforge/internal/config"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(cfg *config.Config, aws_config *aws.Config) *s3.Client {
	endpoint := ""
	if cfg.Env == "dev" {
		endpoint = "http://localhost:4566"
	} else {
		endpoint = "https://" + cfg.CfAccountId + ".r2.cloudflarestorage.com"
	}

	return s3.NewFromConfig(*aws_config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
}

func ExtractS3Key(path string) string {
	path = strings.TrimPrefix(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
