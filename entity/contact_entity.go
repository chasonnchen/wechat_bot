package entity

type ContactEntity struct {
	Id     string
	Name   string
	Type   int32
	Appid   int32
	Hello  string
	OpenAi int32
	Status int32
}

func (ContactEntity) TableName() string {
	return "t_contact"
}
