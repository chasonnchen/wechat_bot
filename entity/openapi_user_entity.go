package entity

type OpenapiUserEntity struct {
	Id        int32
	Name      string
	ContactId string
	Tel string
	AppId   int32
	AppKey string
	Status    int32
}

func (OpenapiUserEntity) TableName() string {
	return "t_openapi_user"
}
