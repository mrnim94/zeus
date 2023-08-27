package aws_cloud_impl

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"zeus/log"
	"zeus/model"
)

func (ac *AWSConfiguration) DeleteAWSKey(userName string, accessKeyId string) error {

	cfg, err := ac.accessAWSCloud()
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return err
	}

	client := iam.NewFromConfig(cfg)

	_, err = client.DeleteAccessKey(context.TODO(), &iam.DeleteAccessKeyInput{
		UserName:    aws.String(userName),
		AccessKeyId: aws.String(accessKeyId),
	})
	if err != nil {
		log.Error("Error deleting access key %s: %v", accessKeyId, err)
		return err
	} else {
		log.Info("Deleted old access key: ", accessKeyId)
	}
	return nil
}

func (ac *AWSConfiguration) RetentionAWSKey(userName string) (*iam.CreateAccessKeyOutput, model.OldAWSCredential, error) {
	var oldCreds model.OldAWSCredential
	cfg, err := ac.accessAWSCloud()
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return nil, oldCreds, err
	}

	// Create IAM client
	client := iam.NewFromConfig(cfg)

	// Step 1: Create a new access key
	newAccessKey, err := client.CreateAccessKey(context.TODO(), &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return nil, oldCreds, err
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
			oldCreds.ListOLDKeys = append(oldCreds.ListOLDKeys, *keyMetadata.AccessKeyId)
		}
	}

	return newAccessKey, oldCreds, nil
}
