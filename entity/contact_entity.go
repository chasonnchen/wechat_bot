package entity

type ContactEntity struct {
	Id   string
	Name   string
	Type   int32
	Status int32
}

func (ContactEntity) TableName() string {
	return "t_contact"
}
