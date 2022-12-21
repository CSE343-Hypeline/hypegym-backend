package models

type Trainer struct {
	UserID  int      `gorm:"primaryKey"`
	Members []Member `gorm:"foreignKey:TrainerID"`
}

type AssignMemberDto struct {
	UserID int `json:"user_id"`
}
