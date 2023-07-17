package aws_cloud_impl

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"zeus/helper/aws_cloud"
	"zeus/log"
)

type AWSConfiguration struct {
	Region string
}

func NewAWSConnection(ac *AWSConfiguration) aws_cloud.AWSCloud {
	return &AWSConfiguration{
		Region: ac.Region,
	}
}

func (ac *AWSConfiguration) accessAWSCloud() (aws.Config, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(ac.Region),
	)
	if err != nil {
		log.Error("Unable to load AWS SDK config, ", err)
	}

	return cfg, nil
}
