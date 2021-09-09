package task

import (
	"log"
	"strconv"
	"time"

	"github.com/chasonnchen/wechat_bot/lib/sina7x24"

	"github.com/wechaty/go-wechaty/wechaty"
)

var (
	sina7x24Task = &Sina7x24Task{LastId: 0}
)

type Sina7x24Task struct {
	LastId int32
	Bot    *wechaty.Wechaty
}

func NewSina7x24Task(bot *wechaty.Wechaty) *Sina7x24Task {
	sina7x24Task.Bot = bot
	return sina7x24Task
}

func (s *Sina7x24Task) Start() {
	s.work()
	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				s.work()
			}
		}
	}()
}

func (s *Sina7x24Task) work() {
	msg, id := sina7x24.NewClient().GetMsgs(10, s.LastId)
	if id > 0 {
		s.LastId = id
	}

	// 晚上10点半到早上8点半 不推送
	layout := "1504"
	timeStr, _ := strconv.Atoi(time.Now().Format(layout))
	if timeStr > 2230 || timeStr < 830 {
		log.Println("It is not good time")
		return
	}

	if len(msg) > 0 {
		//s.Bot.Room().Load("18543635738@chatroom").Say(msg)
		//s.Bot.Room().Load("21083279973@chatroom").Say(msg)
	}
}