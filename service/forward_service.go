package service

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chasonnchen/wechat_bot/dao"
	"github.com/chasonnchen/wechat_bot/entity"
	"github.com/chasonnchen/wechat_bot/lib/util"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	forwardServiceInstance = &ForwardService{}
)

func NewForwardService() *ForwardService {
	return forwardServiceInstance
}

type ForwardService struct {
	ForwardList     []entity.SkillForwardEntity
	GroupRelatinMap map[int32][]entity.GroupRelationEntity
}

func (f *ForwardService) neesIgnore(message *user.Message) bool {
	if message.Self() {
		log.Println("Message discarded because its outgoing")
		return true
	}

	if message.Age() > 2*60*time.Second {
		log.Println("Message discarded because its TOO OLD(than 2 minutes)")
		return true
	}

	if message.Type() != schemas.MessageTypeText {
		log.Println("Message discarded because it dose not Text")
		// TODO return true
	}

	return false
}

func (f *ForwardService) DoForward(contact entity.ContactEntity, message *user.Message) {
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

func (f *ForwardService) buildMsgHead(message *user.Message) string {
	return util.BuildMsgFrom(message) + ":\n" + message.Text()
}

func (f *ForwardService) forward(contact entity.ContactEntity, message *user.Message, forward entity.SkillForwardEntity) {
	// 检查group有没有配置触发联系人
	groupRelationList, ok := f.GroupRelatinMap[forward.ToGroupId]
	if !ok || len(groupRelationList) < 1 {
		log.Printf("Forward fial. This group has no relaion. forward info is %#v", forward)
		return
	}

	// 这里写一个cache，记录 某个发言人，触发了某条转发规则，记1分钟，用于后面媒体消息的自动识别和转发

	NewCacheService().Set(message.From().ID()+strconv.Itoa(int(forward.Id)), 1, 1*time.Minute)

	for _, relation := range groupRelationList {
		// 判断下不发给当前消息来源的群，或者个人
		if contact.Id == relation.ContactId || message.From().ID() == relation.ContactId {
			continue
		}
		// 通过转发实现 (只配置关注人，未配置关键字，且接收方是个人，且消息类型是非文本)
		if len(forward.Spekers) > 0 && len(forward.Keywords) < 1 && NewContactService().GetById(relation.ContactId).Type == 1 && message.Type() != schemas.MessageTypeText {
			// 因为直接转发，不知道消息的来源，所以先发一条文本提示
			NewContactService().SayTextToContact(relation.ContactId, "主人好~\n您订阅的"+util.BuildMsgFrom(message)+"发了一条媒体消息，转发请查收：")
			message.Forward(relation.ContactId)
		}
		// 非文本消息不使用say发送
		if message.Type() != schemas.MessageTypeText {
			continue
		}
		// 通过say实现
		NewContactService().SayTextToContact(relation.ContactId, f.buildMsgHead(message))
	}
}

func (f *ForwardService) checkFrom(contact entity.ContactEntity, message *user.Message, forward entity.SkillForwardEntity) bool {
	// 检查group有没有配置触发联系人
	groupRelationList, ok := f.GroupRelatinMap[forward.FromGroupId]
	if !ok || len(groupRelationList) < 1 {
		log.Printf("This group has no relaion. forward info is %#v", forward)
		return false
	}

	// 检查是否符合联系人来源（包括个人和群组）
	notMatchContactId := true
	for _, relation := range groupRelationList {
		if relation.ContactId == contact.Id {
			notMatchContactId = false
			break
		}
	}
	if notMatchContactId {
		return false
	}

	// 检查是否符合发言人要求, 只有群消息 且配置了发言人字段才检查
	if contact.Type == 2 && len(forward.Spekers) > 0 {
		if !strings.Contains(forward.Spekers, message.From().ID()) {
			return false
		}
	}

	// 检查是否符合关键字
	if len(forward.Keywords) > 0 {
		notMatchKeyword := true
		for _, keyword := range strings.Split(forward.Keywords, ",") {
			if strings.Contains(message.Text(), keyword) {
				notMatchKeyword = false
			}
		}
		if notMatchKeyword {
			return false
		}
	}

	// 检查是否有屏蔽字，有的话不转发
	if len(forward.BadKeywords) > 0 {
		for _, keyword := range strings.Split(forward.BadKeywords, ",") {
			if strings.Contains(message.Text(), keyword) {
				return false
			}
		}
	}

	log.Printf("This message match forward. forward info is %#v", forward)
	return true
}

func (f *ForwardService) init() {
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

func (f *ForwardService) load() {
	groupRelatinMap := make(map[int32][]entity.GroupRelationEntity)

	var forwardList []entity.SkillForwardEntity
	dao.Webot().Where("status = ?", "1").Find(&forwardList)
	var groupRelationList []entity.GroupRelationEntity
	dao.Webot().Where("status = ?", "1").Find(&groupRelationList)

	for _, relation := range groupRelationList {
		if len(groupRelatinMap[relation.GroupId]) < 1 {
			relationList := make([]entity.GroupRelationEntity, 0)
			groupRelatinMap[relation.GroupId] = relationList
		}
		groupRelatinMap[relation.GroupId] = append(groupRelatinMap[relation.GroupId], relation)
	}
	log.Printf("forward conf is %#v", forwardList)
	log.Printf("group relation conf is %#v", groupRelatinMap)

	f.ForwardList = forwardList
	f.GroupRelatinMap = groupRelatinMap
}
