package database

import (
	"log"

	"github.com/sung2708/shorten_url/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&model.User{}, &model.URL{}, &model.OTP{})
	if err != nil {
		log.Fatal("Failed to auto migrate users, urls, or email verification tokens:", err)
	}

	log.Println("Successfully migrated users, urls, and email verification tokens")
	return db
}
