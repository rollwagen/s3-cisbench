package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"

	"github.com/aws/smithy-go"

	"github.com/rollwagen/s3-cisbench/internal/output"

	"github.com/briandowns/spinner"
	"github.com/rollwagen/s3-cisbench/internal/audit"
	"github.com/rollwagen/s3-cisbench/internal/aws"
	"github.com/spf13/cobra"
)

var outputFormat string

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit S3 buckets against applicable CIS benchmark items",
	Long:  `Audit S3 buckets against applicable CIS benchmark items`,
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[11], 80*time.Millisecond)
		s.Suffix = " Getting S3 buckets..."
		s.Start()

		var reports []audit.BucketReport
		bucketAuditor := audit.New()
		buckets, err := aws.GetBuckets()
		if err != nil {
			s.Stop()
			var e smithy.APIError
			if errors.As(err, &e) {
				fmt.Printf("Error listing S3 buckets: %v: %v", e.ErrorCode(), e.ErrorMessage())
			} else {
				fmt.Printf(color.RedString("Unexpected error: ")+"%v", err)
			}
			os.Exit(1)
		}
		s.Suffix = " Auditing buckets..."
		for _, b := range buckets {
			reports = append(reports, bucketAuditor.Report(b.Name, b.AccountID, b.Region))
		}
		s.Stop()

		switch {
		case outputFormat == "txt":
			output.PrintReport(reports)
			//fmt.Println("...skipping output...")
		case outputFormat == "json":
			b, _ := json.MarshalIndent(reports, "", "  ")
			fmt.Println(string(b))
		case outputFormat == "csv":
			fmt.Println("TODO csv output")
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// auditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// auditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	auditCmd.Flags().StringVarP(&outputFormat, "output", "o", "txt", "Define outputFormat report (txt, csv, json)")
}
