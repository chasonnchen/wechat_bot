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

func (q *QaService) getContactIdFromMessage(message *user.Message) string {
	if message.Room() != nil {
		return message.Room().ID()
	}

	return message.From().ID()
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

func (q *QaService) DoQa(message *user.Message) {
	// 1. 检查是否需要忽略
	if q.neesIgnore(message) {
		return
	}

	// 2. 开始匹配问答
	contactId := q.getContactIdFromMessage(message)

	for _, qaItem := range q.QaConf[contactId] {
		for _, keyword := range strings.Split(qaItem.QaKey, ",") {
			if strings.Contains(message.Text(), keyword) {
				_, err := message.Say(qaItem.QaValue)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("Message response is %s", qaItem.QaValue)
				return
			}
		}
	}

	// TODO 未命中关键字时，ai聊天
	log.Println("Message discarded because not match any keyword.")
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
