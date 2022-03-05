package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main () {

	router := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "127.0.0.1:8083", "", []byte("bj38"))
	router.Use(sessions.Sessions("mysession", store))
	router.GET("/test", func(ctx * gin.Context) {
		// ctx.SetCookie("mytest", "chuanzhi", 60 * 60, "", "", false, false)
		// cookieVal, _ := ctx.Cookie("mytest")
		// fmt.Println("get cookie:", cookieVal)
		s := sessions.Default(ctx)
		//s.Set("itcast", "itheima")
		//s.Save()
		v := s.Get("itcast")
		fmt.Println("v:", v)
		ctx.Writer.WriteString("ass hole !")
	})
	router.Run(":9999")
}
