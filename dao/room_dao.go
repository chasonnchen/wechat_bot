package dao

import (
	"github.com/chasonnchen/wechat_bot/entity"
)

var (
	roomDao = &RoomDao{}
)

type RoomDao struct {
}

func NewRoomDao() *RoomDao {
	return roomDao
}

func (r *RoomDao) GetAll() (roomList []entity.RoomEntity) {
	conn := getDb("webot")
	conn.Find(&roomList)
	return
}

func (r *RoomDao) Insert(room entity.RoomEntity) {
	conn := getDb("webot")
	conn.Create(&room)
	return
}
