package util

import (
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func BuildMsgFrom(message *user.Message) string {
	var msgText string
	if message.Room() != nil {
		aliasName, err := message.Room().Alias(message.From())
		if err != nil || len(aliasName) < 1 {
			aliasName = message.From().Name()
		}
		msgText = "[" + aliasName + "@" + message.Room().Topic()
	} else {
		name := message.From().Alias()
		if len(name) < 1 {
			name = message.From().Name()
		}
		msgText = "[" + name
	}
	msgText = msgText + "]"

	return msgText
}
