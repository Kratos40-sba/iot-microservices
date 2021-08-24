package model

import (
	"gihub.com/Kratos40-sba/auth-service/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" ,gorm:"unique"` //todo make the email unique
	Password string `json:"password"`
}

func(u *User) CreateUserRecord () error {

	record := database.DB.Create(&u)
	if record.Error != nil {
		return record.Error
	}
	return nil
}
func(u *User) HashRfidNumber(rfid string) error  {
	bytes , err := bcrypt.GenerateFromPassword([]byte(rfid),14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}
func(u *User) CheckPassword(providedRfid string) error  {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password),[]byte(providedRfid))
	if err != nil {
		return err
	}
	return nil
}
