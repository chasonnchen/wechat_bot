package entity

type SkillQaEntity struct {
	Id      int32
	ContactId  string
	Name  string
    // 不能为空，支持使用英文逗号分割配置多个关键字，命中任何一个返回
	QaKey   string
	QaValue string
	Status  int32
}

func (SkillQaEntity) TableName() string {
	return "t_skill_qa"
}
