package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitPostgres(hostName, dbName, dbUser, dbPort string) (*gorm.DB, error) {

	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", hostName, dbUser, dbName, dbPort)}), &gorm.Config{
			SkipDefaultTransaction: true,
		})
	if err != nil {
		log.Println("Failed to connect to database!")
		return nil, err
	}
	DB = db
	/*
		err = DB.AutoMigrate(&model.Event{})
			if err != nil {
				log.Fatalln("Error while migrating model : ",err)
				return nil , err
			}
	*/

	return DB, nil
}
