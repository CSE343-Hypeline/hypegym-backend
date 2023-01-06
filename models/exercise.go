package models

import (
	"gorm.io/datatypes"
)

type Exercise struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name" gorm:"unique"`
	Level        string         `json:"level"`
	Equipment    string         `json:"equipment"`
	Instructions datatypes.JSON `json:"instructions"`
}

type ExerciseResponseDto struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
}
