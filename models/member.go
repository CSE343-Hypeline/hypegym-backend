package models

type Member struct {
	UserID    int       `gorm:"primaryKey"`
	TrainerID int       `json:"trainer_id"`
	Programs  []Program `json:"programs" gorm:"foreignKey:UserID"`
}
