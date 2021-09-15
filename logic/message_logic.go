package logic

import (
	"log"
	"strings"

	"github.com/chasonnchen/wechat_bot/entity"
	"github.com/chasonnchen/wechat_bot/lib/baidu/unit"
	"github.com/chasonnchen/wechat_bot/lib/util"
	"github.com/chasonnchen/wechat_bot/service"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

const (
	contactTypeRoom = 2
	contactTypeUser = 1
)

var (
	messageLogicInstance = &MessageLogic{}
)

func NewMessageLogic() *MessageLogic {
	return messageLogicInstance
}

type MessageLogic struct {
}

func (m *MessageLogic) buildContact(message *user.Message) entity.ContactEntity {
	messageRoom := message.Room()
	contact := entity.ContactEntity{}

	if messageRoom != nil {
		contact.Id = messageRoom.ID()
		contact.Name = messageRoom.Topic()
		contact.Type = 2
		contact.Status = 1
	} else {
		contact.Id = message.From().ID()
		contact.Name = message.From().Name()
		if len(message.From().Alias()) > 0 {
			contact.Name = message.From().Alias()
		}
		contact.Type = 1
		contact.Status = 1
	}

	return contact
}

func (m *MessageLogic) buildMsgText(message *user.Message) string {
	msgText := "[" + message.From().ID() + "]" + util.BuildMsgFrom(message)

	if message.Type() != schemas.MessageTypeText {
		msgText = msgText + "[say something not Text.]"
	} else {
		msgText = msgText + ": " + message.Text()
	}

	return msgText
}

func (m *MessageLogic) Do(message *user.Message) {
	// 0. log
	log.Printf("MessageLogic recive message: %s", m.buildMsgText(message))
	contact := m.buildContact(message)

	// 1. 更新联系人
	contact = service.NewContactService().Upsert(contact)

	// 2. 问答
	service.NewQaService().DoQa(contact, message)

	// 3. 转发
	service.NewForwardService().DoForward(contact, message)

	// 4. 暗号加群
	service.NewRoomService().AutoInvite(message.From(), message, "")

	// 5. 智能聊天
	// 5.1 好友聊天 && 打开智能聊天配置 && 文本消息
	if contact.Type == 1 && contact.OpenAi == 1 && message.Type() == schemas.MessageTypeText {
		log.Print("start ai\n")
		aiRes, _ := unit.NewClient().Chat(contact.Id, message.Text())
		if len(aiRes) > 1 {
			message.Say(aiRes)
		}
	}
	// 5.2 群里聊天 && @机器人 && 文本消息 && 非群公告
	// @发言人 并回复智能聊天的结果
	selfAliasName, _ := message.Room().Alias(service.NewGlobleService().GetBot().UserSelf())
	if len(selfAliasName) < 1 {
		selfAliasName = service.NewGlobleService().GetBot().UserSelf().Name()
	}
	log.Printf("self room alias name is %s", selfAliasName)
	if contact.Type == 2 && message.MentionSelf() && message.Type() == schemas.MessageTypeText && strings.Contains(message.Text(), selfAliasName) {
		aiRes, _ := unit.NewClient().Chat(message.From().ID(), message.Text()[strings.Index(message.Text(), string(rune(8197)))+3:])
		if len(aiRes) > 1 {
			message.Room().Say(aiRes, message.From())
		}
	}
}
