/**
 *   Wechaty - https://github.com/wechaty/wechaty
 *
 *   @copyright 2020-now Wechaty
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 *
 */
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/chasonnchen/wechat_bot/configs"
	"github.com/chasonnchen/wechat_bot/dao"
	"github.com/chasonnchen/wechat_bot/logic"
	"github.com/chasonnchen/wechat_bot/service"
	"github.com/chasonnchen/wechat_bot/service/openapi"
	"github.com/chasonnchen/wechat_bot/task"

	"github.com/wechaty/go-wechaty/wechaty"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func main() {

	// 1. 启动bot服务
	var bot = wechaty.NewWechaty(wechaty.WithPuppetOption(wp.Option{
		Endpoint: "127.0.0.1:30009",
		Token:    "2fdb00a5-5c31-4018-84ac-c64e5f995057",
		Timeout:  time.Duration(2 * time.Minute),
	}))

	bot.OnScan(func(ctx *wechaty.Context, qrCode string, status schemas.ScanStatus, data string) {
		log.Printf("Scan QR Code to login: %v\nhttps://wechaty.js.org/qrcode/%s\n", status, qrCode)
	}).OnLogin(func(ctx *wechaty.Context, user *user.ContactSelf) {
		log.Printf("User %s login success! \n", user.Name())
	}).OnMessage(onMessage).OnLogout(func(ctx *wechaty.Context, user *user.ContactSelf, reason string) {
		fmt.Printf("User %s logouted: %s\n", user, reason)
	}).OnFriendship(onFriendship).OnRoomInvite(onRoomInvite).OnRoomJoin(onRoomJoin)
	//.OnRoomLeave().OnRoomTopic()

	var err = bot.Start()
	if err != nil {
		panic(err)
	}
	// 2. 放一个全局bot
	service.NewGlobleService().SetBot(bot)

	// 3. 初始化业务模块
	configs.InitConfig()
	dao.InitDao()
	service.InitService()

	openapi.InitOpenapi()
	task.InitTask()

	// 4. 启动一个http服务
	initServer(":8877")

	var quitSig = make(chan os.Signal)
	signal.Notify(quitSig, os.Interrupt, os.Kill)

	select {
	case <-quitSig:
		log.Fatal("exit.by.signal")
	}
}

func onMessage(ctx *wechaty.Context, message *user.Message) {
	logic.NewMessageLogic().Do(message)
}

func onFriendship(ctx *wechaty.Context, friendship *user.Friendship) {
	logic.NewFriendshipLogic().Do(friendship)
}

func onRoomInvite(ctx *wechaty.Context, roomInvitation *user.RoomInvitation) {
	logic.NewRoomLogic().DoInvite(roomInvitation)
}

func onRoomJoin(ctx *wechaty.Context, room *user.Room, inviteeList []_interface.IContact, inviter _interface.IContact, date time.Time) {
	logic.NewRoomLogic().DoJoin(room, inviteeList, inviter)
}
