package printers

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/rollwagen/s3-cisbench/internal/audit"
)

type TextPrinter struct{}

func (r *TextPrinter) PrintReport(report []audit.BucketReport, w io.Writer) error {
	if w != os.Stdout {
		return errors.New("this printer only support writing to stdout")
	}

	type Glyph string
	const (
		// GlyphVLine   Glyph = "│" // "\uf6d7"
		// GlyphHLine   Glyph = "\u2015"

		GlyphVDotted Glyph = "\uf6d7"
		GlyphHDotted Glyph = "\uE621"
	)

	for _, b := range report {
		fmt.Println()

		cBucket := color.New(color.FgHiBlue).Add(color.Bold)
		colorBucketPrintln := func(a any) {
			_, _ = cBucket.Println(a)
		}
		colorBucketPrint := func(a any) {
			_, _ = cBucket.Print(a)
		}

		// Bucket name
		colorBucketPrintln(" \uF5A7 " + b.Name)

		// CIS 2.1.1 - ServerSideEncryptionEnabled
		colorBucketPrint(" " + GlyphHDotted)
		cCIS := color.New(color.FgHiCyan)
		_, _ = cCIS.Println("\tEnsure all S3 buckets employ encryption-at-rest [CIS 2.1.1]")

		colorBucketPrint(" " + GlyphHDotted)
		if b.ServerSideEncryptionEnabled {
			// for _, rule := range encryptionOutput.ServerSideEncryptionConfiguration.Rules {
			//	if rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == "AES256" {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t\ufc98 ")
			// _, _ = c.Print("\t\t\uf046 ")
			c = color.New(color.FgGreen)
			// _, _ = cBucket.Println(" Server side encryption with AES256 is enabled")
			if b.CustomerManagedKey {
				_, _ = c.Println(" Server side encryption is enabled with customer managed \uf80a")
			} else {
				_, _ = c.Println(" Server side encryption is enabled")
			}
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t\uf73f") // cBucket.Print("\t\t\uf071")
			c = color.New(color.FgRed)
			_, _ = c.Println(" No server side encryption found")
		}

		// CIS 2.1.2 - Ensure S3 Bucket Policy is set to deny HTTP requests
		colorBucketPrint(" " + GlyphHDotted)
		cCIS = color.New(color.FgHiCyan)
		_, _ = cCIS.Println("\tEnsure S3 Bucket Policy is set to deny HTTP requests [CIS 2.1.2]")
		colorBucketPrint(" " + GlyphHDotted)
		if b.PolicyDenyHTTP {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t\uf023 ")
			c = color.New(color.FgGreen)
			_, _ = c.Println(" Bucket policy to deny HTTP requests is present")
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t\uf09c") // cBucket.Print("\t\t\uf071")
			c = color.New(color.FgRed)
			_, _ = c.Println(" No Bucket policy to deny HTTP requests found")
		}

		// Non CIS - Versioning enabled
		colorBucketPrint(" " + GlyphHDotted)
		_, _ = cCIS.Println("\tS3 bucket versioning enabled (non-CIS)")

		colorBucketPrint(" " + GlyphHDotted)
		if b.VersioningEnabled {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t\uf454 ")
			c = color.New(color.FgGreen)
			_, _ = c.Println(" S3 bucket has versioning enabled")
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t\uf73f") // cBucket.Print("\t\t\uf071")
			c = color.New(color.FgRed)
			_, _ = c.Println(" Versioning is not enabled")
		}

		// CIS 2.1.5 Ensure that S3 Buckets are configured with 'Block public access'
		/*
		 * ✖ ✔ BlockPublicAcls (BPA)
		 * ✖ ✔ BlockPublicPolicy (BPP)
		 * ✖ ✔ IgnorePublicAcls (IPA)
		 * ✖ ✔ RestrictPublicBuckets (RPB)
		 */
		colorBucketPrint(" " + GlyphHDotted)
		_, _ = cCIS.Println("\tEnsure that S3 Buckets are configured with 'Block public access' [CIS 2.1.5]")
		colorBucketPrint(" " + GlyphHDotted)
		if b.BlockPublicAccess.BlockPublicAcls {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t✔")
			c = color.New(color.FgGreen)
			_, _ = c.Println(" Block Public ACLs is enabled")
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t✖")
			c = color.New(color.FgRed)
			_, _ = c.Println(" Block Public ACLs is disabled")
		}

		colorBucketPrint(" " + GlyphHDotted)
		if b.BlockPublicAccess.BlockPublicPolicy {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t✔")
			c = color.New(color.FgGreen)
			_, _ = c.Println(" Block Public Policy is enabled")
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t✖")
			c = color.New(color.FgRed)
			_, _ = c.Println(" Block Public Policy is disabled")
		}

		colorBucketPrint(" " + GlyphHDotted)
		if b.BlockPublicAccess.IgnorePublicAcls {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t✔")
			c = color.New(color.FgGreen)
			_, _ = c.Println(" Ignore Public ACLs is enabled")
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t✖")
			c = color.New(color.FgRed)
			_, _ = c.Println(" Ignore Public ACLs is disabled")
		}

		colorBucketPrint(" " + GlyphHDotted)
		if b.BlockPublicAccess.RestrictPublicBuckets {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t✔")
			c = color.New(color.FgGreen)
			_, _ = c.Println(" Restrict Public Access is enabled")
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t✖")
			c = color.New(color.FgRed)
			_, _ = c.Println(" Restrict Public Access is disabled")
		}

		// bucket report END
		colorBucketPrintln(" " + GlyphVDotted + "\n")

		// color.Green("Ξ" + "⚠⚠" + "✗✗" + "☡☡" + "∆∆" + "≈≈")
	}
	return nil
}
