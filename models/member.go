package models

type Member struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id"`
	TrainerID int       `json:"trainer_id"`
	Programs  []Program `json:"programs" gorm:"foreignKey:UserID"`
}
