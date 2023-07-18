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
	log.InitLogger(false)
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}

func main() {

	awsCloud := &aws_cloud_impl.AWSConfiguration{Region: "us-east-1"}

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
