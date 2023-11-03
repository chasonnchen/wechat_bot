package logic

import (
	"log"
	"time"
    "math/rand"
    "strconv"

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
        return
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
		service.NewContactService().SayTextToContact(contact.ID(), "Hi,我们是Webot团队~\n咱们已经是好友啦，可以看看下面功能详细介绍哦~")
		time.Sleep(3 * time.Second)
        service.NewContactService().SayTextToContact(contact.ID(), "官方网站：https://webot.cc")
		service.NewContactService().SayTextToContact(contact.ID(), "github地址：https://github.com/chasonnchen/wechat_bot")
		time.Sleep(5 * time.Second)
        service.NewContactService().SayTextToContact(contact.ID(), "1. 如果您有群暗号，发给我可以自动邀请您进群~\n\n2. 如果您想使用此机器人，请先看下官网介绍和使用接入方式哦~\n\n3. 如咨询其他问题请留言，正在分配客服~")
		time.Sleep(10 * time.Second)
        rand.Seed(time.Now().UnixNano())
        name := rand.Intn(30) + 10
		service.NewContactService().SayTextToContact(contact.ID(), "您好~ 我是"+ strconv.Itoa(int(name)) +"号客服，您请说")

		// 通用QA

		// 如果上面hello命中口令， 发送邀请定制欢迎语，并发送邀请
		service.NewRoomService().AutoInvite(contact, nil, friendship.Hello())
	}
}
