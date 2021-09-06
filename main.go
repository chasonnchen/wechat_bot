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
	"github.com/chasonnchen/wechat_bot/service"
	"github.com/chasonnchen/wechat_bot/task"

	"github.com/wechaty/go-wechaty/wechaty"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func main() {

	// 1. 初始化
	/*configs.InitConfig()
	  dao.InitDao()
	  service.InitService()
	  task.InitTask()*/

	// 2. 启动bot服务
	var bot = wechaty.NewWechaty(wechaty.WithPuppetOption(wp.Option{
		//Endpoint: "127.0.0.1:8788",
		//Token: "puppet_padlocal_96b1f6b62b074fa3bd3a8ae4c55fa56a",
		Timeout: time.Duration(2 * time.Minute),
	}))

	bot.OnScan(func(ctx *wechaty.Context, qrCode string, status schemas.ScanStatus, data string) {
		fmt.Printf("Scan QR Code to login: %v\nhttps://wechaty.js.org/qrcode/%s\n", status, qrCode)
	}).OnLogin(func(ctx *wechaty.Context, user *user.ContactSelf) {
		fmt.Printf("User %s logined, user info is %v\n", user.Name(), user)
	}).OnMessage(onMessage).OnLogout(func(ctx *wechaty.Context, user *user.ContactSelf, reason string) {
		fmt.Printf("User %s logouted: %s\n", user, reason)
	})

	var err = bot.Start()
	if err != nil {
		panic(err)
	}

	// 1. 初始化
	configs.InitConfig()
	dao.InitDao()
	service.InitService()
	task.InitTask(bot)

	var quitSig = make(chan os.Signal)
	signal.Notify(quitSig, os.Interrupt, os.Kill)

	select {
	case <-quitSig:
		log.Fatal("exit.by.signal")
	}
}
