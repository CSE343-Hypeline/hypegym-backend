package models

type Program struct {
	ID         int  `gorm:"primaryKey" json:"id"`
	UserID     uint `json:"user_id" `
	ExerciseID uint `gorm:"unique" json:"exercise_id"`
	Set        uint `json:"set"`
	Repetition uint `json:"repetition"`
}

type ProgramResponseDto struct {
	ID         int  `json:"id"`
	ExerciseID uint `json:"exercise_id"`
	Set        uint `json:"set"`
	Repetition uint `json:"repetition"`
}

type ProgramRequestDto struct {
	ProgramID uint `json:"program_id"`
}
