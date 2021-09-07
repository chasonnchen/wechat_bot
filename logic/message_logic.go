package logic

import (
    "log"

	"github.com/chasonnchen/wechat_bot/entity"
	"github.com/chasonnchen/wechat_bot/service"

    "github.com/wechaty/go-wechaty/wechaty/user"
    "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
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
        contact.Type = 1
        contact.Status = 1
    }

    return contact
}

func (m *MessageLogic) buildMsgText(message *user.Message) string {
    msgText := "[" + message.From().Name()
    if message.Room() != nil {
        msgText = msgText + "@" + message.Room().Topic()
    }
    msgText = msgText + "]: "

    if message.Type() != schemas.MessageTypeText {
        msgText = msgText + "[send something but not Text.]"
    }
    msgText = msgText + message.Text()

    return msgText
}

func (m *MessageLogic) Do(message *user.Message) {
    // 0. log
    log.Printf("MessageLogic recive message: %s", m.buildMsgText(message))

    // 1. 更新联系人
    service.NewContactService().Upsert(m.buildContact(message))

    // 2. 问答
    service.NewQaService().DoQa(message)

    // 3. 转发
}
