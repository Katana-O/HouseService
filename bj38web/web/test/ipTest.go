package main

import (
	"github.com/gin-gonic/gin"
)

func main () {
	router := gin.Default()
	router.GET("/", func(ctx * gin.Context) {
		ctx.Writer.WriteString("ass hole !")
	})
	router.Run("10.0.2.15:8088")
}
