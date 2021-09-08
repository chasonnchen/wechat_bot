package logic

import (
	"log"

	"github.com/chasonnchen/wechat_bot/entity"
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
		contact.Name = message.From().Alias()
		contact.Type = 1
		contact.Status = 1
	}

	return contact
}

func (m *MessageLogic) buildMsgText(message *user.Message) string {
	var msgText string
	if message.Room() != nil {
		aliasName, err := message.Room().Alias(message.From())
		if err != nil || len(aliasName) < 1 {
			aliasName = message.From().Name()
		}
		msgText = "[" + aliasName + "@" + message.Room().Topic()
	} else {
		msgText = "[" + message.From().Alias()
	}
	msgText = msgText + "]: "

	if message.Type() != schemas.MessageTypeText {
		msgText = msgText + "[say something not Text.]"
	}
	msgText = msgText + message.Text()

	return msgText
}

func (m *MessageLogic) Do(message *user.Message) {
	// 0. log
	log.Printf("MessageLogic recive message: %s", m.buildMsgText(message))
	contact := m.buildContact(message)

	// 1. 更新联系人
	service.NewContactService().Upsert(contact)

	// 2. 问答
	service.NewQaService().DoQa(contact, message)

	// 3. 转发
	service.NewForwardService().DoForward(contact, message)
}
