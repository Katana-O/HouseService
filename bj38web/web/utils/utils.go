package utils

import "github.com/micro/go-micro"

func InitMicro() micro.Service {
	microService := micro.NewService()
	return microService
}
