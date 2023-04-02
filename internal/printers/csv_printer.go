package printers

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/rollwagen/s3-cisbench/internal/audit"
)

type CSVPrinter struct{}

func (r *CSVPrinter) PrintReport(reports []audit.BucketReport, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	var data [][]string
	data = append(data, []string{
		"Account Id",
		"Region",
		"Bucket Name",
		"Server Side Encryption",
		"Versioning enabled",
		"MFA delete",
		"Deny HTTP only",
		"Block Public ACLs",
		"Ignore Public ACLs",
		"Block Public Policy",
		"Restrict Public Buckets",
	})
	for _, r := range reports {
		bpa := r.BlockPublicAccess
		row := []string{
			r.AccountID,
			r.Region,
			r.Name,
			strconv.FormatBool(r.ServerSideEncryptionEnabled),
			strconv.FormatBool(r.VersioningEnabled),
			strconv.FormatBool(r.MFADelete),
			strconv.FormatBool(r.PolicyDenyHTTP),
			strconv.FormatBool(bpa.BlockPublicAcls),
			strconv.FormatBool(bpa.IgnorePublicAcls),
			strconv.FormatBool(bpa.BlockPublicPolicy),
			strconv.FormatBool(bpa.RestrictPublicBuckets),
		}
		data = append(data, row)
	}
	_ = csvWriter.WriteAll(data)

	return nil
}
