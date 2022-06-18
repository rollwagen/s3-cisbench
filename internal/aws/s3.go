package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

type Bucket struct {
	Name      string `json:"bucketName"`
	AccountID string `json:"accountId"`
	Region    string `json:"region"`
}

var ctx = context.Background()

func newS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Errorf("Failed to load AWS SDK configuration: %v", err)
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}

func GetBucketNamesWithPrefix(prefix string) ([]string, error) {
	s3Client, _ := newS3Client()
	result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Errorf("Error listing S3 buckets: %v", err)
		return nil, err
	}

	var bucketNames []string
	for _, b := range result.Buckets {
		name := b.Name
		if strings.HasPrefix(*name, prefix) {
			bucketNames = append(bucketNames, *name)
		}
	}
	return bucketNames, nil
}

func GetBucketByName(name string) (Bucket, error) {
	s3Client, _ := newS3Client()
	cfg, _ := config.LoadDefaultConfig(ctx)
	region, err := manager.GetBucketRegion(ctx, s3Client, name)
	if err != nil {
		return Bucket{}, err
	}
	accountID, _ := GetAccountID(&cfg)
	bucket := Bucket{name, accountID, region}

	return bucket, nil
}

func GetBuckets() ([]Bucket, error) {
	var buckets []Bucket

	log.Debug("Listing buckets")
	s3Client, _ := newS3Client()
	result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Errorf("Error listing S3 buckets: %v", err)
		return nil, err
	}

	for _, b := range result.Buckets {
		name := b.Name
		region, _ := manager.GetBucketRegion(ctx, s3Client, *name)
		cfg, _ := config.LoadDefaultConfig(ctx)
		accountID, _ := GetAccountID(&cfg)
		buckets = append(buckets, Bucket{*name, accountID, region})
	}
	return buckets, nil
}
