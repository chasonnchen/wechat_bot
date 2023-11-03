package service

import (
	"log"
	"strconv"
	"time"

	"github.com/chasonnchen/wechat_bot/dao"
	"github.com/chasonnchen/wechat_bot/entity"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	forwardMediaServiceInstance = &ForwardMediaService{}
)

func NewForwardMediaService() *ForwardMediaService {
	return forwardMediaServiceInstance
}

type ForwardMediaService struct {
	ForwardList []entity.SkillForwardEntity
}

func (f *ForwardMediaService) neesIgnore(message *user.Message) bool {
	if message.Self() {
		log.Println("Message discarded because its outgoing")
		return true
	}

	if message.Age() > 2*60*time.Second {
		log.Println("Message discarded because its TOO OLD(than 2 minutes)")
		return true
	}

	if message.Type() == schemas.MessageTypeText {
		log.Println("Forward Media Message discarded because it is Text")
		return true
	}

	return false
}

func (f *ForwardMediaService) DoForward(contact entity.ContactEntity, message *user.Message) {
	// 1. 检查是否需要忽略
	if f.neesIgnore(message) {
		return
	}

	// 2. 开始匹配转发
	for _, forward := range f.ForwardList {
		if f.checkFrom(contact, message, forward) {
			// 执行转发动作
			f.forward(contact, message, forward)
		}
		// 可能触发多次转发，所以这里不退出
	}
}

func (f *ForwardMediaService) forward(contact entity.ContactEntity, message *user.Message, forward entity.SkillForwardEntity) {
	// 检查group有没有配置触发联系人
	cidList := NewGroupService().GetContactIdListByGroupId(forward.ToGroupId)
	if len(cidList) < 1 {
		log.Printf("ForwardMedia fial. This group has no relaion. forward info is %#v", forward)
		return
	}

	for _, cid := range cidList {
		// 判断下不发给当前消息来源的群，或者个人
		if contact.Id == cid || message.From().ID() == cid {
			continue
		}
		// 通过转发实现
		message.Forward(cid)
	}
}

func (f *ForwardMediaService) checkFrom(contact entity.ContactEntity, message *user.Message, forward entity.SkillForwardEntity) bool {
	time.Sleep(5 * time.Second)
	_, found := NewCacheService().Get(message.From().ID() + strconv.Itoa(int(forward.Id)))
	if found {
		log.Printf("This message match forward media.")
		return true
	}

	return false
}

func (f *ForwardMediaService) init() {
	f.load()

	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				f.load()
			}
		}
	}()
}

func (f *ForwardMediaService) load() {
	var forwardList []entity.SkillForwardEntity
	dao.Webot().Where("status = ?", "1").Find(&forwardList)

	f.ForwardList = forwardList
}
