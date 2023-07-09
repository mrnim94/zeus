package k8s_impl

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"zeus/log"
)

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
