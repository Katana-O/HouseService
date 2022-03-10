package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"userOrder/handler"
	model "userOrder/mode"
	_ "userOrder/subscriber"

	userOrder "userOrder/proto/userOrder"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.userOrder"),
		micro.Version("latest"),
	)

	// Initialise service
	// service.Init()
	model.InitDb()

	// Register Handler
	userOrder.RegisterUserOrderHandler(service.Server(), new(handler.UserOrder))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
