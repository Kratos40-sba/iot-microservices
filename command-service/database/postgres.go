package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitPostgres(host, dbName, dbUser, dbPort string) (*gorm.DB, error) {

	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", host, dbUser, dbName, dbPort)}), &gorm.Config{
			SkipDefaultTransaction: true,
		})
	if err != nil {
		log.Println("Failed to connect to database!")
		return nil, err
	}
	DB = db
	return DB, nil
}
