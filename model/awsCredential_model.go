package model

type AWSCredential struct {
	AccessKey string
	SecretKey string
}

type OldAWSCredential struct {
	ListOLDKeys []string
}
