package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadConfig() (err error)  {
	err = godotenv.Load()
	if err != nil {
		log.Println("Error While loading app.env file")
		return err
	}
	return nil
}

