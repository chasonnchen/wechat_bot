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
	"strings"
	"time"

	"github.com/wechaty/go-wechaty/wechaty"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func main() {
	var bot = wechaty.NewWechaty(wechaty.WithPuppetOption(wp.Option{
		Endpoint: "127.0.0.1:10002",
		//Token: "3d415ebb-7a6f-4cba-b602-1f4ae400f011",
		Timeout: time.Duration(2 * time.Minute),
	}))

	bot.OnScan(func(ctx *wechaty.Context, qrCode string, status schemas.ScanStatus, data string) {
		fmt.Printf("Scan QR Code to login: %v\nhttps://wechaty.js.org/qrcode/%s\n", status, qrCode)
	}).OnLogin(func(ctx *wechaty.Context, user *user.ContactSelf) {
		fmt.Printf("User %s logined\n", user.Name())
	}).OnMessage(onMessage).OnLogout(func(ctx *wechaty.Context, user *user.ContactSelf, reason string) {
		fmt.Printf("User %s logouted: %s\n", user, reason)
	})

	var err = bot.Start()
	if err != nil {
		panic(err)
	}

	var quitSig = make(chan os.Signal)
	signal.Notify(quitSig, os.Interrupt, os.Kill)

	select {
	case <-quitSig:
		log.Fatal("exit.by.signal")
	}
}

func onMessage(ctx *wechaty.Context, message *user.Message) {
	log.Println(message)

	if message.Self() {
		log.Println("Message discarded because its outgoing")
	}

	if message.Age() > 2*60*time.Second {
		log.Println("Message discarded because its TOO OLD(than 2 minutes)")
	}

	if message.Type() != schemas.MessageTypeText {
		log.Println("Message discarded because it dose not Text")
		return
	}

	if message.Room() != nil {
		fmt.Printf("message from room,info is %v\n", message.Room())
	}

	// 1. 配置关键字和答复的map
	msgMap := make(map[string]string)
	msgMap["在哪儿"] = "客官好！~猫眼宠物店地址是：浙江省嘉兴市秀洲区花园街311号，电话：15157409090,15805839253。欢迎随时来店里体验~"

	// 2. 命中关键字自动答复并返回
	for k, value := range msgMap {
		if strings.Contains(message.Text(), k) {
			_, err := message.Say(value)
			if err != nil {
				log.Println(err)
				return
			}
			//log.Println("REPLY: " + v)
			return
		}
	}

	// TODO 未命中关键字时，ai聊天
	log.Println("Message discarded because not match any keyword.")
	return
	//log.Println("REPLY: dong")

	// 2. reply image(qrcode image)
	//fileBox, _ := file_box.FromUrl("https://wechaty.github.io/wechaty/images/bot-qr-code.png", "", nil)
	//_, err = message.Say(fileBox)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Printf("REPLY: %s\n", fileBox)
}
