package initializers

import (
	"hypegym-backend/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var dbError error

func ConnectToDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	connectionString := os.Getenv("DB_URL")
	DB, dbError = gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: newLogger})
	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database!")
}

func MigrateDB() {
	DB.AutoMigrate(&models.User{}, &models.Gym{}, &models.Member{}, &models.Trainer{})

	var superadmin models.User
	superadmin.Name = "super"
	superadmin.Email = "superadmin@superadmin"
	superadmin.HashPassword("superadmin")
	superadmin.PhoneNumber = "11111111111"
	superadmin.Address = "Gebze No: 1"
	superadmin.GymID = 1
	superadmin.Role = "SUPERADMIN"

	/* 	var pt models.User
	   	pt.Name = "pt"
	   	pt.Email = "pt@pt"
	   	pt.HashPassword("superadmin")
	   	pt.PhoneNumber = "1111111asd1111"
	   	pt.Address = "Gebze No: 1"
	   	pt.GymID = 1
	   	pt.Role = "SUPERADMIN"

	   	var member models.User
	   	member.Name = "member"
	   	member.Email = "member@member"
	   	member.HashPassword("superadmin")
	   	member.PhoneNumber = "11111111111"
	   	member.Address = "Gebze No: 1"
	   	member.GymID = 1
	   	member.Role = "SUPERADMIN" */

	var gym models.Gym
	gym.Name = "HYPEGYM"
	gym.Address = "CSE343"
	gym.Email = "hypegym@gmail.com"
	gym.PhoneNumber = "11111111111"

	DB.Create(&gym)
	DB.Create(&superadmin)
	/* 	DB.Create(&pt)
	   	DB.Create(&member) */
	log.Println("Database Migration Completed!")
}
