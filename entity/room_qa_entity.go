package entity

type RoomQaEntity struct {
	Id      int32
	RoomId  int32
	QaName  string
	QaKey   string
	QaValue string
	Status  int32
}

func (RoomQaEntity) TableName() string {
	return "t_room_qa"
}
