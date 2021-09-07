package service

import (
    "strings"

	"github.com/chasonnchen/wechat_bot/dao"
	"github.com/chasonnchen/wechat_bot/entity"
)

var (
	contactServiceInstance = &ContactService{}
)

func NewContactService() *ContactService {
	return contactServiceInstance
}

type ContactService struct {
    ContactList map[string]entity.ContactEntity
}

func (c *ContactService) getContactTypeById(contactId string) int32 {
    if strings.HasSuffix(contactId, "@chatroom") {
        return 2
    }
    return 1
}

func (c *ContactService) Upsert(contact entity.ContactEntity) {
    //  先检查在不在List里
    contactOri, ok := c.ContactList[contact.Id]
    if ok {
        // 判等
        if contactOri.Name != contact.Name {
            contactOri.Name = contact.Name
            c.ContactList[contact.Id] = contactOri

            // 更新DB
            dao.Webot().Model(&contactOri).Update("name", contactOri.Name)
        }
    } else {
        // 插入
        // 注意这里填充type
        contact.Type = c.getContactTypeById(contact.Id)
        c.ContactList[contact.Id] = contact
        dao.Webot().Create(&contact)
    }
}

func (c *ContactService) init() {
    // 加载数据库中已保存的联系人
    contactMap := make(map[string]entity.ContactEntity)
    var contactList []entity.ContactEntity
    dao.Webot().Find(&contactList)

    for _, contact := range contactList{
        contactMap[contact.Id] = contact
    }

    c.ContactList = contactMap
    return
}
