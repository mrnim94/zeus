package main

import (
	"github.com/labstack/echo/v4"
	"os"
	"zeus/handler"
	"zeus/helper/aws_cloud/aws_cloud_impl"
	"zeus/log"
)

func init() {
	os.Setenv("APP_NAME", "zeus")
	log.InitLogger(false)
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}

func main() {

	awsCloud := &aws_cloud_impl.AWSConfiguration{Region: "us-east-1"}
	rotateKeyHandler := handler.RotateKeyHandler{
		AWSCloud: aws_cloud_impl.NewAWSConnection(awsCloud),
	}

	e := echo.New()
	rotateKeyHandler.HandlerCreateDeleteKey()
	e.Logger.Fatal(e.Start(":1323"))
}
