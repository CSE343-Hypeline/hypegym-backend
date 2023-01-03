package initializers

import (
	"context"
	"hypegym-backend/models"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var DBError error
var RDB *redis.Client
var CTX context.Context = context.Background()

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
	DB, DBError = gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: newLogger})
	if DBError != nil {
		log.Fatal(DBError)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database!")
}

func ConnectToRedis() {
	connectionString := os.Getenv("REDIS_ADD")
	RDB = redis.NewClient(&redis.Options{
		Addr:     connectionString,
		Password: "",
		DB:       0,
	})
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
	superadmin.Gender = "OTHER"

	var gym models.Gym
	gym.Name = "HYPEGYM"
	gym.Address = "CSE343"
	gym.Email = "hypegym@gmail.com"
	gym.PhoneNumber = "11111111111"

	DB.Create(&gym)
	DB.Create(&superadmin)
	log.Println("Database Migration Completed!")
}
