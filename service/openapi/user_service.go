package openapi

import (
	"log"
	"time"

	"github.com/chasonnchen/wechat_bot/dao"
	"github.com/chasonnchen/wechat_bot/entity"
)

var (
	userServiceInstance = &UserService{}
)

func NewUserService() *UserService {
	return userServiceInstance
}

type UserService struct {
	UserMap map[int32]entity.OpenapiUserEntity
}

func (o *UserService) init() {
	o.load()

	go func() {
		for {
			select {
			case <-time.After(time.Second * 3600):
				o.load()
			}
		}
	}()
}

func (o *UserService) GetUserByAppId(appId int32) entity.OpenapiUserEntity {
	return o.UserMap[appId]
}

func (o *UserService) GetAppKeyByAppId(appId int32) string {
	user := o.UserMap[appId]
	return user.AppKey
}

func (o *UserService) load() {
	userMap := make(map[int32]entity.OpenapiUserEntity)
	var userList []entity.OpenapiUserEntity
	dao.Webot().Where("status = ?", "1").Find(&userList)

	for _, user := range userList {
		userMap[user.AppId] = user
	}

	log.Printf("openapi user conf is %#v", userMap)
	o.UserMap = userMap
}
