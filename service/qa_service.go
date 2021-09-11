package service

import (
	"log"
	"strings"
	"time"

	"github.com/chasonnchen/wechat_bot/dao"
	"github.com/chasonnchen/wechat_bot/entity"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	qaServiceInstance = &QaService{}
)

func NewQaService() *QaService {
	return qaServiceInstance
}

type QaService struct {
	QaConf map[string][]entity.SkillQaEntity
}

func (q *QaService) neesIgnore(message *user.Message) bool {
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
		return true
	}

	return false
}

func (q *QaService) DoQa(contact entity.ContactEntity, message *user.Message) {
	// 1. 检查是否需要忽略
	if q.neesIgnore(message) {
		return
	}

	// 2. 开始匹配问答
	for _, qaItem := range q.QaConf[contact.Id] {
		for _, keyword := range strings.Split(qaItem.QaKey, ",") {
			if strings.Contains(message.Text(), keyword) {
				if contact.Type == 2 {
					// 群里回答并at管理员
					currRoom := NewGlobleService().GetBot().Room().Load(contact.Id)
					atContact := message.From()
					if qaItem.CallOwner == 1 {
						atContact = currRoom.Owner()
					}
					currRoom.Say(strings.Trim(qaItem.QaValue, "\n"), atContact)
				} else {
					message.Say(strings.Trim(qaItem.QaValue, "\n"))
				}
				log.Printf("Message response is %s", qaItem.QaValue)
			}
		}
	}

	// 3. 单聊，通用问答匹配
	if contact.Type == 1 {
		for _, qaItem := range q.QaConf["@"] {
			for _, keyword := range strings.Split(qaItem.QaKey, ",") {
				if strings.Contains(message.Text(), keyword) {
					message.Say(strings.Trim(qaItem.QaValue, "\n"))
					log.Printf("Message response is %s", qaItem.QaValue)
				}
			}
		}
	}

	return
}

func (q *QaService) init() {
	q.load()

	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				q.load()
			}
		}
	}()
}

func (q *QaService) load() {
	qaConf := make(map[string][]entity.SkillQaEntity)
	var qaList []entity.SkillQaEntity
	dao.Webot().Find(&qaList)

	for _, qaItem := range qaList {
		if len(qaConf[qaItem.ContactId]) < 1 {
			confItem := make([]entity.SkillQaEntity, 0)
			qaConf[qaItem.ContactId] = confItem
		}
		qaConf[qaItem.ContactId] = append(qaConf[qaItem.ContactId], qaItem)
	}

	log.Printf("qa conf is %#v", qaConf)
	q.QaConf = qaConf
}
