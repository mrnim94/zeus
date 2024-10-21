package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"os"
	"zeus/handler"
	"zeus/helper/aws_cloud/aws_cloud_impl"
	"zeus/helper/k8s/k8s_impl"
	"zeus/log"
)

func init() {
	os.Setenv("APP_NAME", "zeus")
	logger := log.InitLogger(false)
	// Check if KUBERNETES_SERVICE_HOST is set
	if _, exists := os.LookupEnv("KUBERNETES_SERVICE_HOST"); !exists {
		// If not in Kubernetes, set LOG_LEVEL to DEBUG
		os.Setenv("LOG_LEVEL", "DEBUG")
	}
	logger.SetLevel(log.GetLogLevel("LOG_LEVEL"))
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}

func main() {

	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		os.Setenv("AWS_REGION", "us-east-1")
	} else {
		log.Info("Zeus is working on ", awsRegion)
	}

	awsCloud := &aws_cloud_impl.AWSConfiguration{Region: awsRegion}

	kubeconfig := flag.String("kubeconfig", ".kube/config", "location to your confighihi file")
	kube := &k8s_impl.KubeConfiguration{KubeConfig: *kubeconfig}

	rotateKeyHandler := handler.RotateKeyHandler{
		AWSCloud: aws_cloud_impl.NewAWSConnection(awsCloud),
		K8s:      k8s_impl.NewKubernetesConnection(kube),
	}

	e := echo.New()
	rotateKeyHandler.HandlerCreateDeleteKey()
	e.Logger.Fatal(e.Start(":1994"))
}
