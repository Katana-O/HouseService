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
		r1.POST("/user/auth", controller.PostUserAuth)
		r1.GET("/user/auth", controller.GetUserInfo)
		r1.GET("/user/houses",controller.GetUserHouses)
		r1.POST("/houses", controller.PostHouses)
		r1.POST("/houses/:id/images",controller.PostHousesImage)
		//展示房屋详情
		r1.GET("/houses/:id",controller.GetHouseInfo)
		//展示首页轮播图
		r1.GET("/house/index",controller.GetIndex)
		//搜索房屋
		r1.GET("/houses",controller.GetHouses)
	}
	router.Run(":8080")
}
