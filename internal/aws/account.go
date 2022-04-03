package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// GetAccountID get the Account ID for the currently logged User.
func GetAccountID(config *aws.Config) (string, error) {
	stsClient := sts.NewFromConfig(*config)
	id, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return *id.Account, err
}
