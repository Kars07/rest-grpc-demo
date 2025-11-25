// database/database.go
package database

import (
	"log"

	"github.com/Kars07/rest-grpc-demo/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	// Using SQLite for simplicity (can switch to PostgreSQL/MySQL)
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected successfully!")
}

func GetDB() *gorm.DB {
	return DB
}
