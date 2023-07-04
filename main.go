package main

import (
	"github.com/labstack/echo/v4"
	"os"
	"zeus/handler"
	"zeus/log"
)

func init() {
	os.Setenv("APP_NAME", "zeus")
	log.InitLogger(false)
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}

func main() {
	e := echo.New()

	handler.HandlerCreateDeleleKey("us-east-1")

	e.Logger.Fatal(e.Start(":1323"))
}
