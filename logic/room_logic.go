package logic

import (
	"log"

	"github.com/chasonnchen/wechat_bot/service"

	"github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	roomLogicInstance = &RoomLogic{}
)

func NewRoomLogic() *RoomLogic {
	return roomLogicInstance
}

type RoomLogic struct {
}

func (r *RoomLogic) DoInvite(roomInvitation *user.RoomInvitation) {
	// 0. log
	contact, _ := roomInvitation.Inviter()
	roomTopic, _ := roomInvitation.Topic()
	log.Printf("RoomLogic recive invite by [%s][%s] to room[%s]: %s", contact.ID(), contact.Name(), roomTopic)

	// 直接自动通过
	// roomInvitation.Accept()
	// 回复邀请人
	service.NewContactService().SayTextToContact("fenglinyexing", "主人，我收到一个入群邀请")
}

func (r *RoomLogic) buildNameString(inviteeList []_interface.IContact) string {
	nameText := "{"
	for _, invitee := range inviteeList {
		nameText = nameText + "[" + invitee.ID() + "][" + invitee.Name() + "], "
	}
	nameText = nameText + "}"

	return nameText
}

func (r *RoomLogic) DoJoin(room *user.Room, inviteeList []_interface.IContact, inviter _interface.IContact) {
	log.Printf("RoomLogic recive join to [%s][%s].  %s intived by [%s][%s]", room.ID(), room.Topic(), r.buildNameString(inviteeList), inviter.ID(), inviter.Name())

	contact := service.NewContactService().GetById(room.ID())
	if len(contact.Hello) > 0 {
		service.NewGlobleService().GetBot().Room().Load(room.ID()).Say(contact.Hello, inviteeList...)
	}
}
