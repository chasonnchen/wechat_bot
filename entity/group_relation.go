package entity

type GroupRelationEntity struct {
	Id        int32
	GroupId   int32
	ContactId string
	Status    int32
}

func (GroupRelationEntity) TableName() string {
	return "t_group_relation"
}
