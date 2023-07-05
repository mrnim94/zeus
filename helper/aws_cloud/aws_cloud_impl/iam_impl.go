package aws_cloud_impl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"zeus/log"
)

func (ac *AWSConfiguration) RetentionAWSKey(userName string) (*iam.AccessKey, error) {

	sess, err := ac.accessAWSCloud()
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return nil, err
	}

	// Create IAM service client
	svc := iam.New(sess)

	// Provide the username

	// Step 1: Create a new access key
	newAccessKey, err := svc.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return nil, err
	}

	// Step 2: List old access keys
	oldAccessKeys, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		log.Error("Error listing access keys: %v", err)
	}

	// Step 3: Delete old access keys
	for _, keyMetadata := range oldAccessKeys.AccessKeyMetadata {
		if *keyMetadata.AccessKeyId != *newAccessKey.AccessKey.AccessKeyId {
			_, err := svc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
				UserName:    aws.String(userName),
				AccessKeyId: aws.String(*keyMetadata.AccessKeyId),
			})
			if err != nil {
				log.Error("Error deleting access key %s: %v", *keyMetadata.AccessKeyId, err)
			} else {
				log.Info("Deleted old access key: %s\n", *keyMetadata.AccessKeyId)
			}
		}
	}

	return newAccessKey.AccessKey, nil
}
