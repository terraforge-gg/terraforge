package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/terraforge-gg/terraforge/internal/config"
)

func NewAwsConfig(cfg *config.Config) (*aws.Config, error) {
	ctx := context.Background()
	aws_cfg, err := aws_config.LoadDefaultConfig(ctx,
		aws_config.WithRegion(cfg.AwsRegion),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AwsAccessKeyId,
			cfg.AwsSecretAccessKey,
			cfg.AwsSessionToken,
		)),
	)

	return &aws_cfg, err
}
