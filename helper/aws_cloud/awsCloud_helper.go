package aws_cloud

import (
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"zeus/model"
)

type AWSCloud interface {
	RetentionAWSKey(userName string) (*iam.CreateAccessKeyOutput, model.OldAWSCredential, error)
	DeleteAWSKey(userName string, accessKeyId string) error
	DeleteAllAWSKey(userName string) error
}
