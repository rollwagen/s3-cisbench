package printers

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/rollwagen/s3-cisbench/internal/audit"
)

type CSVPrinter struct{}

func (r *CSVPrinter) PrintReport(reports []audit.BucketReport, w io.Writer) error {
	b := func(b bool) string {
		return strconv.FormatBool(b)
	}

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
			b(r.ServerSideEncryptionEnabled),
			b(r.VersioningEnabled),
			b(r.MFADelete),
			b(r.PolicyDenyHTTP),
			b(bpa.BlockPublicAcls),
			b(bpa.IgnorePublicAcls),
			b(bpa.BlockPublicPolicy),
			b(bpa.RestrictPublicBuckets),
		}
		data = append(data, row)
	}
	_ = csvWriter.WriteAll(data)

	return nil
}
