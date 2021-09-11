package logic

import (
	"log"

	"github.com/chasonnchen/wechat_bot/entity"
	"github.com/chasonnchen/wechat_bot/service"

	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	friendshipLogicInstance = &FriendshipLogic{}
)

func NewFriendshipLogic() *FriendshipLogic {
	return friendshipLogicInstance
}

type FriendshipLogic struct {
}

func (f *FriendshipLogic) Do(friendship *user.Friendship) {
	// 0. log
	contact := friendship.Contact()
	log.Printf("FriendshipLogic recive [%s] from [%s][%s]: %s", friendship.Type(), contact.ID(), contact.Name(), friendship.Hello())

	if friendship.Type().String() == "FriendshipTypeReceive" {
		// 收到添加好友请求

		// 直接自动通过
		friendship.Accept()
		// 更新到好友列表
		service.NewContactService().Upsert(entity.ContactEntity{
			Id:     contact.ID(),
			Name:   contact.Name(),
			Type:   1,
			Status: 1,
		})

		// 发送通用欢迎语
		service.NewContactService().SayTextToContact(contact.ID(), "我们已经是好友了，开始聊天吧~\n如果您有群暗号，直接发给我可以在自动邀请您进群哦~\n（本微信是机器人，功能测试中）")

		// 通用QA

		// 如果上面hello命中口令， 发送邀请定制欢迎语，并发送邀请
		service.NewRoomService().AutoInvite(contact, nil, friendship.Hello())
	}
}
