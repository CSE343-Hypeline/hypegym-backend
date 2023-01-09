package models

type Program struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	UserID     uint `json:"user_id" gorm:"uniqueIndex:idx_first_second"`
	ExerciseID uint `json:"exercise_id" gorm:"uniqueIndex:idx_first_second"`
	Set        uint `json:"set"`
	Repetition uint `json:"repetition"`
}

type ProgramResponseDto struct {
	ID         int      `json:"id"`
	ExerciseID uint     `json:"exercise_id"`
	Exercise   Exercise `json:"exercise"`
	Set        uint     `json:"set"`
	Repetition uint     `json:"repetition"`
}

type ProgramRequestDto struct {
	ProgramID uint `json:"program_id"`
}
