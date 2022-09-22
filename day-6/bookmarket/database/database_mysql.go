package database

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	dbSource := os.Getenv("DB_SOURCE")

	// Open connection to db
	db, err := gorm.Open(mysql.Open(dbSource), &gorm.Config{})
	if err != nil {
		log.Fatal("Connection to database failed!, err: ", err.Error())
	}

	return db
}

func SetupTestDB() *gorm.DB {
	dbSource := os.Getenv("DB_SOURCE_TEST")
	// Open connection to db
	db, err := gorm.Open(mysql.Open(dbSource), &gorm.Config{})
	if err != nil {
		log.Fatal("Connection to database failed!, err: ", err.Error())
	}

	return db
}
