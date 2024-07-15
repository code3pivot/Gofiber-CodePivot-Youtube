package configs

import (
	"log"

	"blogappgolang/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Declare a structure it will be a pointer which will hold the database connection instance
type Dbinstance struct {
	Db *gorm.DB
}

// Global Variable It will hold database instance throughout the application
var DB Dbinstance

func ConnectDb() {
	dsn := "host=localhost user=postgres password='123456789' dbname=blogappgo port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("DATABASE CONNECTED")
	db.Logger = logger.Default.LogMode(logger.Info)
	// log.Println("running migrations")
	// Auto Migration Of Models
	db.AutoMigrate(&models.User{}, &models.UserToken{}, &models.Blog{}, &models.Category{})
	DB = Dbinstance{
		Db: db,
	}
}
