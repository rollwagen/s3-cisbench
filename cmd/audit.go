package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rollwagen/s3-cisbench/internal/printers"

	"github.com/fatih/color"

	"github.com/aws/smithy-go"

	"github.com/briandowns/spinner"
	"github.com/rollwagen/s3-cisbench/internal/audit"
	"github.com/rollwagen/s3-cisbench/internal/aws"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var outputFormat string

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit S3 buckets against applicable CIS benchmark items",
	Long:  `Audit S3 buckets against applicable CIS benchmark items`,
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[11], 60*time.Millisecond)
		if !debug { // no spinner when debug output enabled
			s.Start()
		}
		s.Suffix = " Getting S3 buckets..."

		var reports []audit.BucketReport
		bucketAuditor := audit.New()
		buckets, err := aws.GetBuckets()
		if err != nil {
			s.Stop()
			var e smithy.APIError
			if errors.As(err, &e) {
				log.Errorf("Error listing S3 buckets: %v: %v", e.ErrorCode(), e.ErrorMessage())
			} else {
				log.Errorf(color.RedString("Unexpected error: ")+"%v", err)
			}
			os.Exit(1)
		}
		s.Suffix = " Auditing buckets..."
		for i, b := range buckets {
			s.Suffix = fmt.Sprintf(" Auditing buckets: [%d/%d] %s...", i, len(buckets), b.Name)
			reports = append(reports, bucketAuditor.Report(b.Name, b.AccountID, b.Region))
		}
		s.Suffix = " Printing report..."
		s.Stop()

		writer := os.Stdout
		var p printers.BucketReportPrinter
		switch {
		case outputFormat == "txt":
			p = &printers.TextPrinter{}
		case outputFormat == "json":
			p = &printers.JSONPrinter{}
		case outputFormat == "csv":
			p = &printers.CSVPrinter{}
		case outputFormat == "noout":
			p = &printers.NooutPrinter{}
		}

		_ = p.PrintReport(reports, writer)
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.Flags().StringVarP(&outputFormat, "output", "o", "txt", "Define outputFormat report (txt, csv, json, noout)")
}
