package printers

import (
	"io"

	"github.com/rollwagen/s3-cisbench/internal/audit"
	log "github.com/sirupsen/logrus"
)

// NooutPrinter that does nothin i.e. swallows all output
type NooutPrinter struct{}

func (r *NooutPrinter) PrintReport(_ []audit.BucketReport, _ io.Writer) error {
	log.Debug("Omitting output because set to 'noout'")
	return nil
}
