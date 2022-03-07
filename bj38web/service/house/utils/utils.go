package utils

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
)

func InitMicro() micro.Service {
	microService := micro.NewService()
	return microService
}

func GetMicroClient() client.Client{
	microService := micro.NewService(
	)
	return microService.Client()
}