package audit

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

type BucketReport struct {
	Name                        string `json:"name"`
	AccountID                   string `json:"accountId"`
	Region                      string `json:"region"`
	ServerSideEncryptionEnabled bool   `json:"serverSideEncryptionEnabled"`
	VersioningEnabled           bool   `json:"versioningEnabled"`
	MFADelete                   bool   `json:"mfaDelete"`

	BlockPublicAccess struct {
		BlockPublicAcls       bool `json:"blockPublicAcls"`
		BlockPublicPolicy     bool `json:"blockPublicPolicy"`
		IgnorePublicAcls      bool `json:"ignorePublicAcls"`
		RestrictPublicBuckets bool `json:"restrictPublicBuckets"`
	}
}

type BucketAuditor struct{}

func New() *BucketAuditor {
	return &BucketAuditor{}
}

func (auditor *BucketAuditor) Report(bucketName string, accountID string, region string) BucketReport {

	//log.SetLevel(log.DebugLevel)

	bucketReport := BucketReport{Name: bucketName, AccountID: accountID, Region: region}

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	s3Client := s3.NewFromConfig(cfg)

	bucketReport.VersioningEnabled = false
	input := &s3.GetBucketVersioningInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
	versioningOutput, err := s3Client.GetBucketVersioning(context.TODO(), input)
	if err != nil {
		log.Errorf("Error getting bucket versioning status: %v", err)
	} else {

		versioningStatus := versioningOutput.Status
		log.Debug("Versioning status: " + versioningStatus)
		if versioningStatus == "Enabled" {
			bucketReport.VersioningEnabled = true
		}

		mfaDelete := versioningOutput.MFADelete
		log.Debug("MFA Delete: " + mfaDelete)
		if mfaDelete == "Enabled" {
			bucketReport.MFADelete = true
		}
	}

	encryptionInput := &s3.GetBucketEncryptionInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
	encryptionOutput, err := s3Client.GetBucketEncryption(context.TODO(), encryptionInput)
	if err != nil {
		// api error ServerSideEncryptionConfigurationNotFoundError:
		// The server side encryption configuration was not found
		log.Debug("Error getting bucket encryption status.")
		bucketReport.ServerSideEncryptionEnabled = false

	} else {
		bucketReport.ServerSideEncryptionEnabled = true

		for _, rule := range encryptionOutput.ServerSideEncryptionConfiguration.Rules {
			if rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == "AES256" {
				log.Debug("ApplyServerSideEncryptionByDefault.SSEAlgorithm is 'AES256'")
			}
			//c.Println(rule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID)
			//c.Println(rule.BucketKeyEnabled)
		}
	}

	publicAccessBlockInput := &s3.GetPublicAccessBlockInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
	publicAccessBlockOutput, err := s3Client.GetPublicAccessBlock(context.TODO(), publicAccessBlockInput)
	if err != nil {
		log.Debug("Error getting public access block info.")
	} else {
		conf := publicAccessBlockOutput.PublicAccessBlockConfiguration
		bucketReport.BlockPublicAccess.BlockPublicAcls = conf.BlockPublicAcls
		bucketReport.BlockPublicAccess.BlockPublicPolicy = conf.BlockPublicPolicy
		bucketReport.BlockPublicAccess.IgnorePublicAcls = conf.IgnorePublicAcls
		bucketReport.BlockPublicAccess.RestrictPublicBuckets = conf.RestrictPublicBuckets

	}

	/*
		bucketPolicyInput := &s3.GetBucketPolicyInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
		bucketPolicyOutput, err := s3Client.GetBucketPolicy(context.TODO(), bucketPolicyInput)
		if err != nil {
			log.Debug("Error getting bucket policy.")
		} else {
			fmt.Println("---")
			fmt.Println(*bucketPolicyOutput.Policy)
			fmt.Println("---")
		}
	*/

	return bucketReport
}
