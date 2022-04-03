package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/rollwagen/s3-cisbench/internal/aws"

	"github.com/rollwagen/s3-cisbench/internal/audit"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// if debug logging is on or off
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "s3-cisbench",
	Short: "s3-csibench is a tool that analyses S3 bucket against CIS benchmark rules",
	Long:  `s3-csibench is a tool that analyses S3 bucket against CIS benchmark rules. Full details can be found at https://github.com/rollwagen/s3-cisbench`,

	// Uncomment the following line if your bare application has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("Failed to load AWS SDK configuration: %v", err)
		}
		s3Client := s3.NewFromConfig(cfg)
		result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			color.Red("Error listing S3 buckets.")
			os.Exit(1)
		}

		accountID, _ := aws.GetAccountID(&cfg)
		bucketAuditor := audit.New()
		var reports []audit.BucketReport
		for _, bucket := range result.Buckets {
			region, _ := manager.GetBucketRegion(ctx, s3Client, *bucket.Name)
			reports = append(reports, bucketAuditor.Report(*bucket.Name, accountID, region))
			//auditBucket(*bucket.Name, *accountID, region)
		}
		if debug {
			b, _ := json.Marshal(reports)
			log.Debug(string(b))
		}
	},
}

/*
func auditBucket(bucketName string, accountID string, region string) {
	log.SetLevel(log.InfoLevel)
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	s3Client := s3.NewFromConfig(cfg)

	//input := &s3.GetBucketVersioningInput{Bucket: bucketName, ExpectedBucketOwner: accountID}
	//versioningOutput, err := s3Client.GetBucketVersioning(context.TODO(), input)
	//if err != nil {
	//	color.Red("Error getting bucket versioning status: %v", err)
	//} else {
	//	fmt.Println("Versioning status: " + versioningOutput.Status)
	//	fmt.Println("MFA Delete: " + versioningOutput.MFADelete)
	//}

	c := color.New(color.FgHiBlue)
	_, _ = c.Println(" \uF5A7 " + bucketName)
	c = color.New(color.FgHiCyan)
	_, _ = c.Println("\tCIS 2.1.1 Ensure all S3 buckets employ encryption-at-rest")
	encryptionInput := &s3.GetBucketEncryptionInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
	encryptionOutput, err := s3Client.GetBucketEncryption(context.TODO(), encryptionInput)
	if err != nil {
		// api error ServerSideEncryptionConfigurationNotFoundError: The server side encryption configuration was not found
		log.Debug("Error getting bucket encryption status.")
		c := color.New(color.FgHiRed).Add(color.Bold)
		_, _ = c.Print("\t\t\uf73f")
		//c.Print("\t\t\uf071")
		c = color.New(color.FgRed)
		_, _ = c.Println(" No server side encryption found")

	} else {
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
}
*/

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.s3-cisbench.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable verbose logging")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// Example rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
