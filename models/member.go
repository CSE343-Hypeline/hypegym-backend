package models

type Member struct {
	UserID    int `gorm:"primaryKey"`
	TrainerID int
}
