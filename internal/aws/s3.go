package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

type Bucket struct {
	Name      string
	AccountID string
	Region    string
}

func GetBuckets() ([]Bucket, error) {

	var buckets []Bucket

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Errorf("Failed to load AWS SDK configuration: %v", err)
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Debugf("Error listing S3 buckets: %v", err)
		return nil, err
	}

	for _, b := range result.Buckets {
		name := b.Name
		region, _ := manager.GetBucketRegion(ctx, s3Client, *name)
		accountID, _ := GetAccountID(&cfg)
		buckets = append(buckets, Bucket{*name, accountID, region})
	}
	return buckets, nil
}
