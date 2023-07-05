package aws_cloud_impl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

func (ac *AWSConfiguration) accessAWSCloud() (*session.Session, error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ac.Region)},
	)
	if err != nil {
		log.Error("Error creating session: %v", err)
		return nil, err
	}

	return sess, nil
}
