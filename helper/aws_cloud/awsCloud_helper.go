package aws_cloud

import (
	"github.com/aws/aws-sdk-go/service/iam"
)

type AWSCloud interface {
	RetentionAWSKey(userName string) (*iam.AccessKey, error)
}
