package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Test1(ctx * gin.Context) {
	t := time.Now()
	fmt.Println("Test1 111")
	ctx.Next()
	fmt.Println(time.Now().Sub(t))
	fmt.Println("Test1 222")
}

func Test2() gin.HandlerFunc {
	return func(ctx * gin.Context) {
		fmt.Println("Test2 333333333")
		ctx.Abort()
		fmt.Println("TEst 2 end")
	}
}

func main () {
	router := gin.Default()
	router.Use(Test1)
	router.Use(Test2())
	router.GET("/test", func(context *gin.Context) {
		fmt.Println("Default 222")
		context.Writer.WriteString("hello bitch")
	})
	router.Run(":9999")
}
