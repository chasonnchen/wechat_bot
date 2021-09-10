package service

import (
	"log"
	"time"

	"github.com/chasonnchen/wechat_bot/dao"
	"github.com/chasonnchen/wechat_bot/entity"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	roomServiceInstance = &RoomService{}
)

func NewRoomService() *RoomService {
	return roomServiceInstance
}

type RoomService struct {
	InviteConf []entity.SkillInviteEntity
}

func (r *RoomService) AutoInvite(contact _interface.IContact, message *user.Message, hello string) {
	// 1. 检查获取文本消息
	var msgText string
	if message.Type() == schemas.MessageTypeText {
		msgText = message.Text()
	}
	if len(hello) > 0 {
		msgText = hello
	}

	if len(msgText) < 1 {
		return
	}

	// 2. 开始匹配是否触发邀请
	for _, invite := range r.InviteConf {
		//if strings.Contains(msgText, invite.Keyword) {
		if msgText == invite.Keyword {
			// 命中关键字，先回复一句提示
			NewGlobleService().GetBot().Contact().Load(contact.ID()).Say("群暗号正确，已发起进群邀请。")

			// 发送邀请, 并打日志
			NewGlobleService().GetBot().Room().Load(invite.ContactId).Add(contact)

			// TODO 加入消息队列，加群成功后 定时删除好友？否则到上限之后功能就废了

			// 这里continue有可能命中多个群的邀请, 在配置上尽量避免命中多个
			break
		}
	}
}

func (r *RoomService) init() {
	r.loadInviteConf()

	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				r.loadInviteConf()
			}
		}
	}()
}

func (r *RoomService) loadInviteConf() {
	var inviteList []entity.SkillInviteEntity
	dao.Webot().Find(&inviteList)

	log.Printf("room invite conf is %#v", inviteList)
	r.InviteConf = inviteList
}
