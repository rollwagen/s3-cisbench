package printers

import (
	"io"

	"github.com/rollwagen/s3-cisbench/internal/audit"
)

// BucketReportPrinter is an interface that knows how to print BucketReport objects.
type BucketReportPrinter interface {
	// PrintReport receives a report, formats it and prints it to a writer.
	PrintReport([]audit.BucketReport, io.Writer) error
}
