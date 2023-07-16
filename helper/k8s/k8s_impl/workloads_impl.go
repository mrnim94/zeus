package k8s_impl

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
	"zeus/log"
)

func (kc *KubeConfiguration) RestartWorkloads(namespace, kind, workload string) error {
	config, err := kc.accessKubernetes()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), workload, metav1.GetOptions{})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Changing the deployment spec causes a rollout
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}

	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Successfully restarted deployment ", workload, "  in namespace ", namespace)

	return nil
}
