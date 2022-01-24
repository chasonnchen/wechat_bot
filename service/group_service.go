package service

import (
	"log"
	"time"

	"github.com/chasonnchen/wechat_bot/dao"
	"github.com/chasonnchen/wechat_bot/entity"
)

var (
	groupServiceInstance = &GroupService{}
)

func NewGroupService() *GroupService {
	return groupServiceInstance
}

type GroupService struct {
	GroupRelatinMap map[int32][]entity.GroupRelationEntity
}

func (g *GroupService) GetContactIdListByGroupId(groupId int32) []string {
	contactIdList := make([]string, 0)
	relationList, ok := g.GroupRelatinMap[groupId]
	if ok {
		for _, relation := range relationList {
			contactIdList = append(contactIdList, relation.ContactId)
		}
	}

	return contactIdList
}

func (g *GroupService) init() {
	g.load()

	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				g.load()
			}
		}
	}()
}

func (g *GroupService) load() {
	groupRelatinMap := make(map[int32][]entity.GroupRelationEntity)

	var groupRelationList []entity.GroupRelationEntity
	dao.Webot().Find(&groupRelationList)

	for _, relation := range groupRelationList {
		if len(groupRelatinMap[relation.GroupId]) < 1 {
			relationList := make([]entity.GroupRelationEntity, 0)
			groupRelatinMap[relation.GroupId] = relationList
		}
		groupRelatinMap[relation.GroupId] = append(groupRelatinMap[relation.GroupId], relation)
	}
	log.Printf("group service relation conf is %#v", groupRelatinMap)

	g.GroupRelatinMap = groupRelatinMap
}
