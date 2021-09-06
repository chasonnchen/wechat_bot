package service

import (
	"log"
	"strings"
	"time"

	"github.com/chasonnchen/wechat_bot/dao"
	"github.com/chasonnchen/wechat_bot/entity"

	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	roomQaService = &RoomQaService{}
)

func NewRoomQaService() *RoomQaService {
	return roomQaService
}

type RoomQaService struct {
	RoomConf   map[string]entity.RoomEntity
	RoomQaConf map[string][]entity.RoomQaEntity
}

func (r *RoomQaService) OnMessage(ctx *wechaty.Context, message *user.Message) {
	// 1. 参数检查
	if message.Self() {
		log.Println("Message discarded because its outgoing")
		return
	}

	if message.Age() > 2*60*time.Second {
		log.Println("Message discarded because its TOO OLD(than 2 minutes)")
		return
	}

	if message.Type() != schemas.MessageTypeText {
		log.Println("Message discarded because it dose not Text")
		return
	}

	// 2. 开始匹配问答
	weRoomId := message.Room().ID()
	roomConf, ok := r.RoomConf[weRoomId]
	if ok {
		if roomConf.RoomStatus == 1 || roomConf.AiStatus == 1 {
			for _, qaItem := range r.RoomQaConf[weRoomId] {
				if strings.Contains(message.Text(), qaItem.QaKey) {
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
	} else {
		newRoom := entity.RoomEntity{
			WeRoomId:   weRoomId,
			RoomName:   message.Room().Topic(),
			AiStatus:   1,
			RoomStatus: 1,
		}
		dao.NewRoomDao().Insert(newRoom)
	}

	// TODO 未命中关键字时，ai聊天
	log.Println("Message discarded because not match any keyword.")
	return
}

func (r *RoomQaService) init() {
	r.loadConf()

	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				r.loadConf()
			}
		}
	}()
}

func (r *RoomQaService) loadConf() {
	roomConf := make(map[string]entity.RoomEntity)
	roomQaConf := make(map[string][]entity.RoomQaEntity)

	roomList := dao.NewRoomDao().GetAll()
	roomQaList := dao.NewRoomQaDao().GetAll()

	roomMap := make(map[int32]entity.RoomEntity)
	for _, roomItem := range roomList {
		roomConf[roomItem.WeRoomId] = roomItem
		roomMap[roomItem.Id] = roomItem
	}

	for _, qaItem := range roomQaList {
		WeRoomId := roomMap[qaItem.RoomId].WeRoomId
		if len(roomQaConf[WeRoomId]) < 1 {
			confItem := make([]entity.RoomQaEntity, 0)
			roomQaConf[WeRoomId] = confItem
		}
		roomQaConf[WeRoomId] = append(roomQaConf[WeRoomId], qaItem)
	}

	log.Printf("r conf is %#v", roomConf)
	log.Printf("q conf is %#v", roomQaConf)
	r.RoomConf = roomConf
	r.RoomQaConf = roomQaConf
}
