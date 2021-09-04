package entity

import(
)

type RoomEntity struct {
    Id int32
    WeRoomId string
    RoomName string
    AiStatus int32
    RoomStatus int32
}

func (RoomEntity) TableName() string {
    return "t_room"
}
