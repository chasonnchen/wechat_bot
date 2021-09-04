package dao

import (
    "github.com/chasonnchen/wechat_bot/entity"
)

var (
    roomQaDao = &RoomQaDao{}
)

type RoomQaDao struct {
}

func NewRoomQaDao() *RoomQaDao{
    return roomQaDao
}

func (r *RoomQaDao) GetAll() (roomQaList []entity.RoomQaEntity) {
    conn := getDb("webot")
    conn.Find(&roomQaList)
    return
}
