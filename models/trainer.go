package models

type Trainer struct {
	UserID  int      `gorm:"primaryKey"`
	Members []Member `gorm:"foreignKey:TrainerID"`
}

type TrainerDto struct {
	UserID int `json:"user_id"`
}
