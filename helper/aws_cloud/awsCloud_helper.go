package aws_cloud

import (
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type AWSCloud interface {
	RetentionAWSKey(userName string) (*iam.CreateAccessKeyOutput, error)
}
