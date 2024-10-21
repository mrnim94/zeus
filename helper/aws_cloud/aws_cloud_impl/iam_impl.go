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

func (ac *AWSConfiguration) DeleteAllAWSKey(userName string) error {
	cfg, err := ac.accessAWSCloud()
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return err
	}

	client := iam.NewFromConfig(cfg)

	// List all access keys for the user
	listKeysInput := &iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	}

	resp, err := client.ListAccessKeys(context.TODO(), listKeysInput)
	if err != nil {
		log.Error("unable to list access keys for user %s: %w", userName, err)
		return err
	}
	// Iterate through each access key and delete it
	for _, key := range resp.AccessKeyMetadata {
		err := deleteAccessKey(client, userName, *key.AccessKeyId)
		if err != nil {
			log.Errorf("Failed to delete access key %s for user %s: %v", *key.AccessKeyId, userName, err)
		} else {
			log.Infof("Deleted access key %s for user %s", *key.AccessKeyId, userName)
		}
	}

	return nil
}

func deleteAccessKey(client *iam.Client, userName string, accessKeyID string) error {
	_, err := client.DeleteAccessKey(context.TODO(), &iam.DeleteAccessKeyInput{
		UserName:    aws.String(userName),
		AccessKeyId: aws.String(accessKeyID),
	})
	if err != nil {
		log.Errorf("unable to delete access key %s for user %s: %w", accessKeyID, userName, err)
		return err
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
