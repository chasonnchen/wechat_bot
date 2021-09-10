package entity

type SkillInviteEntity struct {
	Id        int32
	Name      string
	Keyword   string
	ContactId string
	Hello     string
	Status    int32
}

func (SkillInviteEntity) TableName() string {
	return "t_skill_invite"
}
