package main

import (
	"bj38web/service/getCaptcha/handler"
	getCaptcha "bj38web/service/getCaptcha/proto/getCaptcha"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
)

func main() {
	//consulReg := consul.NewRegistry()
	// New Service
	service := micro.NewService(
		micro.Address("127.0.0.1:8081"),
		micro.Name("getCaptcha"), // go.micro.srv.getCaptcha
		micro.Version("latest"),
		//micro.Registry(consulReg),
	)
	//fmt.Println("DefaultName 1 :", service.Name())

	// Initialise service
	//service.Init()

	fmt.Println("DefaultName :", service.Name())

	// Register Handler
	getCaptcha.RegisterGetCaptchaHandler(service.Server(), new(handler.GetCaptcha))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
