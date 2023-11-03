package main

import (
	"fmt"

	"github.com/chasonnchen/wechat_bot/api"

	"github.com/gin-gonic/gin"
)

func initServer(port string) {
	ginServer := initRouter()
	if err := ginServer.Run(port); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
		panic(err)
	}
}

func initRouter() *gin.Engine {
	server := gin.Default()

	openapiGroup := server.Group("/openapi")
	openapiGroup.POST("/room/getAll", api.RoomGetAll)
	openapiGroup.POST("/msg/send", api.MsgSend)
	//openapiGroup.GET("/vip/msg/send", api.VipMsgSend)

	return server
}
