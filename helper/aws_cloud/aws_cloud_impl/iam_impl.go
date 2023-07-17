package aws_cloud_impl

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"zeus/log"
)

func (ac *AWSConfiguration) RetentionAWSKey(userName string) (*iam.CreateAccessKeyOutput, error) {
	cfg, err := ac.accessAWSCloud()
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return nil, err
	}

	// Create IAM client
	client := iam.NewFromConfig(cfg)

	// Step 1: Create a new access key
	newAccessKey, err := client.CreateAccessKey(context.TODO(), &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return nil, err
	}

	// Step 2: List old access keys
	oldAccessKeys, err := client.ListAccessKeys(context.TODO(), &iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		log.Error("Error listing access keys: ", err)
	}

	// Step 3: Delete old access keys
	for _, keyMetadata := range oldAccessKeys.AccessKeyMetadata {
		if *keyMetadata.AccessKeyId != *newAccessKey.AccessKey.AccessKeyId {
			_, err := client.DeleteAccessKey(context.TODO(), &iam.DeleteAccessKeyInput{
				UserName:    aws.String(userName),
				AccessKeyId: aws.String(*keyMetadata.AccessKeyId),
			})
			if err != nil {
				log.Error("Error deleting access key %s: %v", *keyMetadata.AccessKeyId, err)
			} else {
				log.Info("Deleted old access key: ", *keyMetadata.AccessKeyId)
			}
		}
	}

	return newAccessKey, nil
}
