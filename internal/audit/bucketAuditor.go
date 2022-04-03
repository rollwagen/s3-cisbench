package audit

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

type BucketReport struct {
	Name                        string `json:"Name"`
	AccountID                   string `json:"AccountId"`
	Region                      string `json:"Region"`
	ServerSideEncryptionEnabled bool   `json:"ServerSideEncryptionEnabled"`
}

type BucketAuditor struct{}

func New() *BucketAuditor {
	return &BucketAuditor{}
}

func (auditor *BucketAuditor) Report(bucketName string, accountID string, region string) BucketReport {

	log.SetLevel(log.DebugLevel)

	bucketReport := BucketReport{Name: bucketName, AccountID: accountID, Region: region}

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	s3Client := s3.NewFromConfig(cfg)

	/*
		input := &s3.GetBucketVersioningInput{Bucket: bucketName, ExpectedBucketOwner: AccountID}
		versioningOutput, err := s3Client.GetBucketVersioning(context.TODO(), input)
		if err != nil {
			color.Red("Error getting bucket versioning status: %v", err)
		} else {
			fmt.Println("Versioning status: " + versioningOutput.Status)
			fmt.Println("MFA Delete: " + versioningOutput.MFADelete)
		}
	*/

	c := color.New(color.FgHiBlue)
	_, _ = c.Println(" \uF5A7 " + bucketName)
	c = color.New(color.FgHiCyan)
	_, _ = c.Println("\tCIS 2.1.1 Ensure all S3 buckets employ encryption-at-rest")
	encryptionInput := &s3.GetBucketEncryptionInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
	encryptionOutput, err := s3Client.GetBucketEncryption(context.TODO(), encryptionInput)
	if err != nil {
		// api error ServerSideEncryptionConfigurationNotFoundError: The server side encryption configuration was not found
		//log.Debug("Error getting bucket encryption status.")
		bucketReport.ServerSideEncryptionEnabled = false

		c := color.New(color.FgHiRed).Add(color.Bold)
		_, _ = c.Print("\t\t\uf73f") //c.Print("\t\t\uf071")
		c = color.New(color.FgRed)
		_, _ = c.Println(" No server side encryption found")

	} else {
		bucketReport.ServerSideEncryptionEnabled = true
		//color.Green("Ξ" + "⚠⚠" + "✗✗" + "☡☡" + "∆∆" + "≈≈")
		for _, rule := range encryptionOutput.ServerSideEncryptionConfiguration.Rules {
			if rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == "AES256" {
				c := color.New(color.FgHiGreen).Add(color.Bold)
				_, _ = c.Print("\t\t\uf046 ")
				c = color.New(color.FgWhite)
				_, _ = c.Println(" Server side encryption with AES256 is enabled")
			}
			//c.Println(rule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID)
			//c.Println(rule.BucketKeyEnabled)
		}
	}
	return bucketReport
}
