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

		for _, location := range schedule.Locations {
			switch location.Style {
			case "CredentialOnK8s":
				log.Info("Action rotates credential: ", location.SecretName)
				rk.K8s.UpdateCredentialInSecret(schedule.NamespaceOnK8s, location.SecretName, location.CredentialOnK8S, location.Profile, accessKey)

			case "AccessKeyOnK8s":
				log.Info("Update New Access Key ID:", *accessKey.AccessKey.AccessKeyId)
				err = rk.K8s.UpdateSecret(schedule.NamespaceOnK8s, location.SecretName, location.AccessKeyOnK8S, *accessKey.AccessKey.AccessKeyId)
				if err != nil {
					log.Error(err)
					return err
				}
				log.Info("Update New Secret Access Key: ******")
				err = rk.K8s.UpdateSecret(schedule.NamespaceOnK8s, location.SecretName, location.SecretKeyOnK8S, *accessKey.AccessKey.SecretAccessKey)
				if err != nil {
					log.Error(err)
					return err
				}

			default:
				log.Info("You don't define style in schedules[i].locations[i].style")
			}
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
}
