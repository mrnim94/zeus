package handler

import (
	"github.com/go-co-op/gocron"
	"time"
	"zeus/helper"
	"zeus/helper/aws_cloud"
	"zeus/helper/k8s"
	"zeus/log"
	"zeus/model"
)

type RotateKeyHandler struct {
	AWSCloud aws_cloud.AWSCloud
	K8s      k8s.K8s
}

func (rk *RotateKeyHandler) HandlerCreateDeleteKey() {

	var cfg model.RotationKey
	helper.LoadConfigFile(&cfg)
	log.Info("Load config file schedules")

	s := gocron.NewScheduler(time.UTC)

	rotationTask := func(schedule model.Schedule) error {
		accessKey, err := rk.AWSCloud.RetentionAWSKey(schedule.UsernameOnAws)

		if err != nil {
			log.Error("Error rotating session: %v", err)
			return err
		}

		log.Info("Update New Access Key ID:", *accessKey.AccessKey.AccessKeyId)
		err = rk.K8s.UpdateSecret(schedule.NamespaceOnK8s, schedule.AccessKeyOnK8S.Name, schedule.AccessKeyOnK8S.Key, *accessKey.AccessKey.AccessKeyId)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Info("Update New Secret Access Key: ******")
		err = rk.K8s.UpdateSecret(schedule.NamespaceOnK8s, schedule.AccessKeyOnK8S.Name, schedule.SecretKeyOnK8S.Key, *accessKey.AccessKey.SecretAccessKey)
		if err != nil {
			log.Error(err)
			return err
		}

		for i, workload := range schedule.RestartWorkloads {
			workload := workload
			log.Info("Restart ", i, " ==> ", workload.Kind, " -->> ", workload.Name, "in namespace ", schedule.NamespaceOnK8s)
			err := rk.K8s.RestartWorkloads(schedule.NamespaceOnK8s, workload.Kind, workload.Name)
			if err != nil {
				log.Error(err)
				return err
			}
		}

		return nil
	}

	for i, schedule := range cfg.Schedules {
		schedule := schedule
		log.Info("Setup Schedule ", i, " ==> ", schedule.Cron)
		_, err := s.Cron(schedule.Cron).Do(rotationTask, schedule)
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
