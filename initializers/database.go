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
	DB.AutoMigrate(&models.User{}, &models.Gym{})

	var user models.User
	user.Email = "superadmin@superadmin"
	user.HashPassword("superadmin")
	user.Role = "SUPERADMIN"
	user.GymID = 1

	var gym models.Gym
	gym.Name = "HYPEGYM"
	gym.Address = "GTU"

	DB.Create(&gym)
	DB.Create(&user)

	log.Println("Database Migration Completed!")
}
