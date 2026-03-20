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
		aws_config.WithRegion(cfg.CfRegion),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.CfAccessKeyId,
			cfg.CfSecretAccessKey,
			"",
		)),
	)

	if err != nil {
		return nil, err
	}

	return &aws_cfg, err
}
