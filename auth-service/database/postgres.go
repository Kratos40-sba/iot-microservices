package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(hostName, userName, dbName, portDB string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s user=%v dbname=%v port=%v sslmode=disable", hostName, userName, dbName, portDB)
	DB, err = gorm.Open(postgres.New(postgres.Config{DSN: dsn}),
		&gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		return nil, err
	}

	return DB, nil
}
