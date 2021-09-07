package entity

type GroupEntity struct {
	Id   int32
	Name   string
	Status int32
}

func (GroupEntity) TableName() string {
	return "t_group"
}
