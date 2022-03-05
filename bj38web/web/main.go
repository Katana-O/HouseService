package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"main/bj38web/web/controller"
	"main/bj38web/web/model"
)

func LoginFilter() gin.HandlerFunc {
	return func(ctx * gin.Context) {
		s := sessions.Default(ctx)
		userName := s.Get("userName")
		if userName == nil {
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}

func main() {
	model.InitRedis()
	model.InitDb()
	router := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "127.0.0.1:8083", "", []byte("bj38"))
	router.Use(sessions.Sessions("mysession", store))
	router.Static("/home", "view")
	r1 := router.Group("/api/v1.0")
	{
		r1.GET("/session", controller.GetSession)
		r1.GET("/imagecode/:uuid", controller.GetImageCd)
		r1.POST("/users", controller.PostRet)
		r1.GET("/areas", controller.GetArea)
		r1.POST("/sessions", controller.PostLogin)
		r1.Use(LoginFilter())
		r1.DELETE("/session", controller.DeleteSession)
		r1.GET("/user", controller.GetUserInfo)
		r1.PUT("/user/name", controller.PutUserInfo)
		r1.POST("/user/avatar", controller.PostAvatar)
	}
	router.Run(":8080")
}
