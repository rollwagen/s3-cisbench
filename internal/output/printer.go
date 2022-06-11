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
		if b.ServerSideEncryptionEnabled {
			// for _, rule := range encryptionOutput.ServerSideEncryptionConfiguration.Rules {
			//	if rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == "AES256" {
			c := color.New(color.FgHiGreen).Add(color.Bold)
			_, _ = c.Print("\t\t\ufc98 ")
			//_, _ = c.Print("\t\t\uf046 ")
			c = color.New(color.FgGreen)
			//_, _ = cBucket.Println(" Server side encryption with AES256 is enabled")
			_, _ = c.Println(" Server side encryption is enabled")
		} else {
			c := color.New(color.FgHiRed).Add(color.Bold)
			_, _ = c.Print("\t\t\uf73f") // cBucket.Print("\t\t\uf071")
			c = color.New(color.FgRed)
			_, _ = c.Println(" No server side encryption found")
		}

		// CIS 2.1.2 - Ensure S3 Bucket Policy is set to deny HTTP requests
		cBucket.Print(" " + GlyphHDotted)
		cCIS = color.New(color.FgHiCyan)
		_, _ = cCIS.Println("\tEnsure S3 Bucket Policy is set to deny HTTP requests [CIS 2.1.2]")
		cBucket.Print(" " + GlyphHDotted)
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
		cBucket.Print(" " + GlyphHDotted)
		_, _ = cCIS.Println("\tS3 bucket versioning enabled (non-CIS)")

		cBucket.Print(" " + GlyphHDotted)
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
		cBucket.Print(" " + GlyphHDotted)
		_, _ = cCIS.Println("\tEnsure that S3 Buckets are configured with 'Block public access' [CIS 2.1.5]")
		cBucket.Print(" " + GlyphHDotted)
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

		cBucket.Print(" " + GlyphHDotted)
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

		cBucket.Print(" " + GlyphHDotted)
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

		cBucket.Print(" " + GlyphHDotted)
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
		cBucket.Println(" " + GlyphVDotted)

		// color.Green("Ξ" + "⚠⚠" + "✗✗" + "☡☡" + "∆∆" + "≈≈")
	}
}
