package printers

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/rollwagen/s3-cisbench/internal/audit"
)

type JSONPrinter struct{}

func (r *JSONPrinter) PrintReport(reports []audit.BucketReport, w io.Writer) error {
	b, _ := json.MarshalIndent(reports, "", "  ")
	_, err := fmt.Fprintln(w, string(b))

	return err
}
