package k8s

import "github.com/aws/aws-sdk-go-v2/service/iam"

type K8s interface {
	UpdateSecret(namespace, secretName, key, value string) error
	RestartWorkloads(namespace, kind, workload string) error
	UpdateCredentialInSecret(namespace, secretName, key, profile string, value *iam.CreateAccessKeyOutput) error
}
