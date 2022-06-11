package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List AWS S3 buckets.",
	Run: func(cmd *cobra.Command, args []string) {
		PrintAllBuckets()
	},
}

func PrintAllBuckets() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS SDK configuration: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)
	log.Debug("Listing buckets")
	result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		var e smithy.APIError
		if errors.As(err, &e) {
			fmt.Printf(color.RedString("Error listing S3 buckets: ")+"%v - %v", e.ErrorCode(), e.ErrorMessage())
		} else {
			fmt.Printf(color.RedString("Unexpected error: ")+"%v", err)
		}
		os.Exit(1)
	}

	log.Infof("Received %d buckets", len(result.Buckets))

	c := color.New(color.FgYellow).Add(color.Underline)
	_, _ = c.Println("Creation date  Bucket name")
	for _, bucket := range result.Buckets {
		fmt.Println("   " + color.BlueString(bucket.CreationDate.Format("2006-01-02")) + "  " + color.CyanString(*bucket.Name))
	}

	c = color.New(color.FgYellow)
	_, _ = c.Println(strconv.Itoa(len(result.Buckets)) + " \uF5A7 overall")
}

func init() {
	rootCmd.AddCommand(listCmd)
}
