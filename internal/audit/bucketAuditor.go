package audit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	log "github.com/sirupsen/logrus"
)

type statement struct {
	Sid    string `json:"Sid,omitempty"` // statement ID
	Effect string `json:"Effect"`        // Allow or Deny
	// Principal    map[string]Value `json:"Principal,omitempty"`    // principal
	Principal    json.RawMessage  `json:"Principal,omitempty"`    // principal
	NotPrincipal map[string]Value `json:"NotPrincipal,omitempty"` // exception to a list of principals
	Action       Value            `json:"Action"`                 // allowed or denied action
	NotAction    Value            `json:"NotAction,omitempty"`    // matches everything except
	Resource     Value            `json:"Resource,omitempty"`     // object or objects that the statement covers
	NotResource  Value            `json:"NotResource,omitempty"`  // matches everything except
	Condition    condition        `json:"condition,omitempty"`
}

type condition map[string]json.RawMessage

type policyDocument struct {
	Version    string      `json:"Version"`
	ID         string      `json:"ID,omitempty"`
	Statements []statement `json:"statement"`
}

// AWS allows string or []string as value
// convert everything to []string to avoid casting
type Value []string

func (value *Value) UnmarshalJSON(b []byte) error {
	var raw interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	var p []string
	//  value could be string or []string -> convert everything to []string
	switch v := raw.(type) {
	case string:
		p = []string{v}
	case []interface{}:
		var items []string
		for _, item := range v {
			items = append(items, fmt.Sprintf("%v", item))
		}
		p = items
	default:
		return fmt.Errorf("invalid %s value element: allowed is only string or []string", value)
	}

	*value = p
	return nil
}

type BucketReport struct {
	Name                        string `json:"name"`
	AccountID                   string `json:"accountId"`
	Region                      string `json:"region"`
	ServerSideEncryptionEnabled bool   `json:"serverSideEncryptionEnabled"`
	VersioningEnabled           bool   `json:"versioningEnabled"`
	MFADelete                   bool   `json:"mfaDelete"`
	PolicyDenyHTTP              bool   `json:"policyDenyHttp"`

	BlockPublicAccess struct {
		BlockPublicAcls       bool `json:"blockPublicAcls"`
		BlockPublicPolicy     bool `json:"blockPublicPolicy"`
		IgnorePublicAcls      bool `json:"ignorePublicAcls"`
		RestrictPublicBuckets bool `json:"restrictPublicBuckets"`
	}
}

type BucketAuditor struct{}

func New() *BucketAuditor {
	return &BucketAuditor{}
}

func (auditor *BucketAuditor) Report(bucketName string, accountID string, region string) BucketReport {
	logBucket := log.WithFields(log.Fields{
		"bucket_name": bucketName,
	})

	bucketReport := BucketReport{Name: bucketName, AccountID: accountID, Region: region}

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	s3Client := s3.NewFromConfig(cfg)

	bucketReport.VersioningEnabled = false
	input := &s3.GetBucketVersioningInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
	versioningOutput, err := s3Client.GetBucketVersioning(context.TODO(), input)
	if err != nil {
		logBucket.Debugf("Error getting versioning status for bucket %s: %v", bucketName, err)
	} else {

		versioningStatus := versioningOutput.Status
		logBucket.Debugf("Versioning status: %#v", versioningStatus)
		if versioningStatus == "Enabled" {
			bucketReport.VersioningEnabled = true
		}

		mfaDelete := versioningOutput.MFADelete
		logBucket.Debugf("MFA Delete: %#v", mfaDelete)
		if mfaDelete == "Enabled" {
			bucketReport.MFADelete = true
		}
	}

	encryptionInput := &s3.GetBucketEncryptionInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
	encryptionOutput, err := s3Client.GetBucketEncryption(context.TODO(), encryptionInput)
	if err != nil {
		// api error ServerSideEncryptionConfigurationNotFoundError:
		// The server side encryption configuration was not found
		logBucket.Debug("Error getting bucket encryption status.")
		bucketReport.ServerSideEncryptionEnabled = false

	} else {
		bucketReport.ServerSideEncryptionEnabled = true

		for _, rule := range encryptionOutput.ServerSideEncryptionConfiguration.Rules {
			if rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == "AES256" {
				logBucket.Info("SSEAlgorithm is 'AES256'")
			}
			// c.Println(rule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID)
			// c.Println(rule.BucketKeyEnabled)
		}
	}

	publicAccessBlockInput := &s3.GetPublicAccessBlockInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
	publicAccessBlockOutput, err := s3Client.GetPublicAccessBlock(context.TODO(), publicAccessBlockInput)
	if err != nil {
		logBucket.Debug("Error getting public access block info.")
	} else {
		conf := publicAccessBlockOutput.PublicAccessBlockConfiguration
		bucketReport.BlockPublicAccess.BlockPublicAcls = conf.BlockPublicAcls
		bucketReport.BlockPublicAccess.BlockPublicPolicy = conf.BlockPublicPolicy
		bucketReport.BlockPublicAccess.IgnorePublicAcls = conf.IgnorePublicAcls
		bucketReport.BlockPublicAccess.RestrictPublicBuckets = conf.RestrictPublicBuckets

	}

	// 2.1.2 Ensure S3 Bucket Policy is set to deny HTTP requests
	// https://aws.amazon.com/premiumsupport/knowledge-center/s3-bucket-policy-for-config-rule/
	// https://docs.fugue.co/FG_R00100.html
	// {"Version":"2012-10-17","statement":[{"Sid":"AWSCloudTrailAclCheck20150319","Effect":"Allow","Principal":{"Service":"clou
	bucketReport.PolicyDenyHTTP = false

	bucketPolicyInput := &s3.GetBucketPolicyInput{Bucket: &bucketName, ExpectedBucketOwner: &accountID}
	bucketPolicyOutput, err := s3Client.GetBucketPolicy(context.TODO(), bucketPolicyInput)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			logBucket.Debugf("%s:  %s", ae.ErrorCode(), ae.ErrorMessage())
		}
	} else {
		var policyDocument policyDocument
		err = json.Unmarshal([]byte(*bucketPolicyOutput.Policy), &policyDocument)
		if err != nil {
			logBucket.Errorf("Error unmarshalling json %v", err)
		}
		logPolicy := logBucket.WithFields(log.Fields{"policy_id": policyDocument.ID})
		logPolicy.Debugf("Processing policy...")
		for _, statement := range policyDocument.Statements {
			// "Effect": "Deny" ?
			if statement.Effect == "Deny" {
				// -  "condition" (1/2): { "Bool"  ?
				for conditionOperator, rawJSON := range statement.Condition {
					if conditionOperator == "Bool" {
						denyUnsecureTransport := false
						conditionKeyValue := &map[string]string{}
						_ = json.Unmarshal(rawJSON, conditionKeyValue)
						if value, ok := (*conditionKeyValue)["aws:SecureTransport"]; ok {
							boolValue, _ := strconv.ParseBool(value)
							//- condition (2/2) { "aws:SecureTransport": true ?
							if !boolValue {
								logPolicy.Debug("aws:SecureTransport is enforced.")
								denyUnsecureTransport = true
							}
						}

						// -  "Action": "*"  or  "Action": "s3:*"  ?
						s3ActionsCovered := false
						for _, action := range statement.Action {
							if action == "*" || action == "s3:*" {
								logPolicy.Debug("s3ActionsCovered is true")
								s3ActionsCovered = true
							}
						}

						// -  "Principal": "*"  or "Principal": { "AWS": "*" } ?
						principalCovered := false
						p := string(statement.Principal)
						if p == "*" || p == "\"*\"" || p == "'*'" {
							principalCovered = true
						} else {
							var principal map[string]string
							_ = json.Unmarshal(statement.Principal, &principal)
							for key, value := range principal {
								if key == "AWS" && value == "*" {
									logPolicy.Debug("principalCovered is true")
									principalCovered = true
								}
							}
						}

						// -  "Resource":  "Resource":"<bucket arn>/*" + "Resource":"<bucket arn>" ?
						resourceBucket := false
						resourceBucketContent := false
						arn := "arn:aws:s3:::" + bucketName
						for _, r := range statement.Resource {
							if strings.HasSuffix(r, "*") && !resourceBucketContent {
								resourceBucketContent = (arn + "/*") == r
							} else if !resourceBucket {
								resourceBucket = arn == r
							}
						}
						bucketResourcesCovered := resourceBucket && resourceBucketContent
						logPolicy.Debugf("bucketResourcesCovered = %v", bucketResourcesCovered)

						bucketReport.PolicyDenyHTTP = denyUnsecureTransport && s3ActionsCovered && principalCovered && bucketResourcesCovered
						logPolicy.Debugf("policyDenyHTTP = %v", bucketReport.PolicyDenyHTTP)
					}
				}
			}
		}
	}

	// done
	return bucketReport
}
