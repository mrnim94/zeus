package k8s_impl

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"zeus/log"
)

func (kc *KubeConfiguration) UpdateCredentialInSecret(namespace, secretName, key, profile string, value *iam.CreateAccessKeyOutput) error {
	config, err := kc.accessKubernetes()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	secret.Data[key] = updateSecretData(secret.Data[key], profile, *value.AccessKey.AccessKeyId, *value.AccessKey.SecretAccessKey)

	// Update the secret in Kubernetes
	_, err = clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Credentials updated successfully in Kubernetes!")
	return nil
}

func (kc *KubeConfiguration) UpdateSecret(namespace, secretName, key, value string) error {

	config, err := kc.accessKubernetes()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	secret.StringData = map[string]string{
		key: value,
	}

	_, err = clientset.CoreV1().Secrets(namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// sub functions
func updateSecretData(data []byte, profile, accessKey, secretKey string) []byte {
	lines := strings.Split(string(data), "\n")
	inProfile := false
	updatedData := ""

	for _, line := range lines {
		if strings.HasPrefix(line, fmt.Sprintf("[%s]", profile)) {
			inProfile = true
		} else if strings.HasPrefix(line, "[") {
			inProfile = false
		}

		if inProfile {
			if strings.HasPrefix(line, "aws_access_key_id") {
				line = fmt.Sprintf("aws_access_key_id = %s", accessKey)
			} else if strings.HasPrefix(line, "aws_secret_access_key") {
				line = fmt.Sprintf("aws_secret_access_key = %s", secretKey)
			}
		}

		updatedData += line + "\n"
	}

	return []byte(updatedData)
}
