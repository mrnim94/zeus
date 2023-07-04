package handler

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/go-co-op/gocron"
	"time"
	"zeus/helper"
	"zeus/log"
	"zeus/model"
)

func HandlerCreateDeleleKey(region string) {

	fmt.Println(region)
	var cfg model.RotationKey
	helper.LoadConfigFile(&cfg)

	s := gocron.NewScheduler(time.UTC)

	rotationTask := func(userName string) {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		if err != nil {
			log.Error("Error creating session: %v", err)
		}
		// Create IAM service client
		svc := iam.New(sess)

		// Provide the username
		//userName := "nimtechnology"

		// Step 1: Create a new access key
		newAccessKey, err := svc.CreateAccessKey(&iam.CreateAccessKeyInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			log.Error("Error creating new access key: %v", err)
		}
		fmt.Println("New Access Key ID:", *newAccessKey.AccessKey.AccessKeyId)
		fmt.Println("New Secret Access Key:", *newAccessKey.AccessKey.SecretAccessKey)

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
	}

	for i, task := range cfg.Tasks {
		task := task
		log.Info("Setup Schedule ", i, " ==> ", task.Cron)
		_, err := s.Cron(task.Cron).Do(rotationTask, task.UsernameOnAws)
		if err != nil {
			log.Error(err)
		}

	}
	s.StartAsync()

	//sess, err := session.NewSession(&aws.Config{
	//	Region: aws.String(region)},
	//)
	//if err != nil {
	//	log.Error("Error creating session: %v", err)
	//}
	//// Create IAM service client
	//svc := iam.New(sess)
	//
	//// Provide the username
	////userName := "nimtechnology"
	//
	//// Step 1: Create a new access key
	//newAccessKey, err := svc.CreateAccessKey(&iam.CreateAccessKeyInput{
	//	UserName: aws.String(userName),
	//})
	//if err != nil {
	//	log.Error("Error creating new access key: %v", err)
	//}
	//fmt.Println("New Access Key ID:", *newAccessKey.AccessKey.AccessKeyId)
	//fmt.Println("New Secret Access Key:", *newAccessKey.AccessKey.SecretAccessKey)
	//
	//// Step 2: List old access keys
	//oldAccessKeys, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
	//	UserName: aws.String(userName),
	//})
	//if err != nil {
	//	log.Error("Error listing access keys: %v", err)
	//}
	//
	//// Step 3: Delete old access keys
	//for _, keyMetadata := range oldAccessKeys.AccessKeyMetadata {
	//	if *keyMetadata.AccessKeyId != *newAccessKey.AccessKey.AccessKeyId {
	//		_, err := svc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
	//			UserName:    aws.String(userName),
	//			AccessKeyId: aws.String(*keyMetadata.AccessKeyId),
	//		})
	//		if err != nil {
	//			log.Error("Error deleting access key %s: %v", *keyMetadata.AccessKeyId, err)
	//		} else {
	//			log.Info("Deleted old access key: %s\n", *keyMetadata.AccessKeyId)
	//		}
	//	}
	//}
}
