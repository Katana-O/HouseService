package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	"house/handler"
	"house/model"
	house "house/proto/house"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.house"),
		micro.Version("latest"),
	)

	// Initialise service
	// service.Init()
	model.InitDb()

	// Register Handler
	house.RegisterHouseHandler(service.Server(), new(handler.House))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
