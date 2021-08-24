package service

import (
	"fmt"
	"gihub.com/Kratos40-sba/command-service/database"
	"gorm.io/gorm"
	"log"
	"net/smtp"
)
const (
	HOSTFORMAT = "%s:%s"
)

type Notification struct {
	gorm.Model
	From string `json:"from"`
	To []string `json:"to"`
	Host string `json:"-"`
	Port string `json:"-"`
	Message []byte  `json:"message"`
}

func NewNotification(from ,port  , host string , message []byte,to []string ) *Notification  {



	return &Notification{
		From: from,
		To: to,
		Message: message,
		Host: host,
		Port: port,
	}
}
func(n *Notification) CreateNotification() error {
	record := database.DB.Create(&n)
	if record.Error != nil {
		return record.Error
	}
	return nil
}

func (n *Notification)SendEmail(password string) error{
	auth := smtp.PlainAuth("",n.From,password,n.Host)
	err := smtp.SendMail(fmt.Sprintf(HOSTFORMAT,n.Host,n.Port),auth,n.From,n.To,n.Message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}