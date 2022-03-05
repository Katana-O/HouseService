package main

import (
	"github.com/micro/go-micro"
	_ "github.com/micro/go-micro/service"
	"github.com/micro/go-micro/util/log"
	"user/handler"
	"user/model"
	user "user/proto/user"
)

func main() {
	model.InitDb()

	model.InitRedis()

	// New Service
	service := micro.NewService(
		micro.Address("127.0.0.1:8082"),
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Initialise service
	// service.Init()

	// Register Handler
	user.RegisterUserHandler(service.Server(), new(handler.User))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
