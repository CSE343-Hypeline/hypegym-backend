package models

import (
	"hypegym-backend/models/enums"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string     `json:"name"`
	Email       string     `json:"email" gorm:"unique"`
	Password    string     `json:"password"`
	PhoneNumber string     `json:"phone_number" gorm:"unique"`
	Address     string     `json:"address"`
	GymID       uint       `json:"gym_id"`
	Role        enums.Role `json:"role" gorm:"type:enum('SUPERADMIN', 'ADMIN', 'PT', 'MEMBER')"`
}

type UserLoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponseDto struct {
	ID    uint   `json:"ID"`
	Email string `json:"email"`
	Role  string `json:"role"`
	GymID uint   `json:"gym_id"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
