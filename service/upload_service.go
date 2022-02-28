package service

import (
	"log"
	"strconv"
	"time"

	"github.com/chasonnchen/wechat_bot/configs"
	"github.com/chasonnchen/wechat_bot/entity"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	uploadServiceInstance = &UploadService{}
)

func NewUploadService() *UploadService {
	return uploadServiceInstance
}

type UploadService struct {
}

func (q *UploadService) neesIgnore(message *user.Message) bool {
	if message.Self() {
		log.Println("Message discarded because its outgoing")
		return true
	}

	if message.Age() > 2*60*time.Second {
		log.Println("Message discarded because its TOO OLD(than 2 minutes)")
		return true
	}

	// 这里先只放开图片上传
	if message.Type() != schemas.MessageTypeImage {
		log.Println("Message discarded because it dose not Image")
		return true
	}

	return false
}

func (q *UploadService) DoUpload(contact entity.ContactEntity, message *user.Message) {
	// 1. 检查是否需要忽略
	if q.neesIgnore(message) {
		return
	}
	// 2. 检查发言人
	if contact.Id != "fenglinyexing" && contact.Id != "wxid_uvf9lcl1otse21" {
		return
	}

	// 3. 保存图片
	file, err := message.ToFileBox()
	if err != nil {
		log.Println("Message save file fail. err is %v", err)
		message.Say("图片上传失败，请几分钟后重试，或者联系技术负责人。")
		return
	}
	log.Println("Message save file success. file name is %s", file.Name)
	// 获取毫秒时间戳做为图片文件名，放到nginx的static img下面
	fileName := strconv.FormatInt(time.Now().UnixNano()/1e6, 10) + ".jpg"
	confUpload := configs.GetConf().Upload
	err = file.ToFile(confUpload.Path+fileName, true)
	if err != nil {
		log.Println("Message save file to user path fail. err is %v", err)
		message.Say("图片上传失败，请几分钟后重试，或者联系技术负责人。")
		return
	}

	// 4. 返回图片地址
	message.Say("http://static.webot.cc/img/" + fileName)

	return
}
