package output

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rollwagen/s3-cisbench/internal/audit"
)

func PrintReport(report []audit.BucketReport) {

	type Glyph string
	const (
		GlyphVLine   Glyph = "│" // "\uf6d7"
		GlyphHLine   Glyph = "\u2015"
		GlyphVDotted Glyph = "\uf6d7"
		GlyphHDotted Glyph = "\uE621"
	)

	for _, b := range report {
		fmt.Println()

		// Bucket name
		cBucket := color.New(color.FgHiBlue).Add(color.Bold)
		_, _ = cBucket.Println(" \uF5A7 " + b.Name)

		// CIS 2.1.1 - ServerSideEncryptionEnabled
		cBucket.Print(" " + GlyphHDotted)
		cCIS := color.New(color.FgHiCyan)
		_, _ = cCIS.Println("\tEnsure all S3 buckets employ encryption-at-rest [CIS 2.1.1]")

		cBucket.Print(" " + GlyphHDotted)
		if b.ServerSideEncryptionEnabled == true {
			//for _, rule := range encryptionOutput.ServerSideEncryptionConfiguration.Rules {
			//	if rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == "AES256" {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t\ufc98 ")
			//_, _ = c.Print("\t\t\uf046 ")
			c = color.New(color.FgGreen)
			//_, _ = cBucket.Println(" Server side encryption with AES256 is enabled")
			_, _ = c.Println(" Server side encryption is enabled")
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t\uf73f") //cBucket.Print("\t\t\uf071")
			c = color.New(color.FgRed)
			_, _ = c.Println(" No server side encryption found")
		}

		// Non CIS - Versioning enabled
		cBucket.Print(" " + GlyphHDotted)
		_, _ = cCIS.Println("\tS3 bucket versioning enabled (non-CIS)")

		cBucket.Print(" " + GlyphHDotted)
		if b.VersioningEnabled == true {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t\uf454 ")
			c = color.New(color.FgGreen)
			_, _ = c.Println(" S3 bucket has versioning enabled")
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t\uf73f") //cBucket.Print("\t\t\uf071")
			c = color.New(color.FgRed)
			_, _ = c.Println(" Versioning is not enabled")
		}

		cBucket.Println(" " + GlyphVDotted)
		//color.Green("Ξ" + "⚠⚠" + "✗✗" + "☡☡" + "∆∆" + "≈≈")
	}
}
