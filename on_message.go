package main

import (
	"log"

    "github.com/chasonnchen/wechat_bot/service"

	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func onMessage(ctx *wechaty.Context, message *user.Message) {
	log.Printf("Webot Recive Message: %#v", message)
    if message.Room() != nil {
        service.NewRoomQaService().OnMessage(ctx, message)
        return
    }
}
