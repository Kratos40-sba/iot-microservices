package model

import (
	"gihub.com/Kratos40-sba/command-service/database"
	"gorm.io/gorm"
)

type CommandType string
type SensorType string

const (
	OFF  CommandType = "OFF"
	ON   CommandType = "ON"
	TEMP SensorType  = "DHT-TEMP"
	HUM  SensorType  = "DHT-HUM"
	SOIL SensorType  = "SOIL"
)

type Command struct {
	gorm.Model
	Time     string      `json:"time"`
	Date     string      `json:"date"`
	Email    string      `json:"email"`
	CommandT CommandType `json:"command_type"`
}
type TimedCommand struct {
	gorm.Model
	Time     string `json:"time"`
	Date     string `json:"date"`
	Email    string `json:"email"`
	Duration int    `json:"duration"`
}
type Trigger struct {
	gorm.Model
	Time        string      `json:"time"`
	Date        string      `json:"date"`
	Email       string      `json:"email"`
	Value       float64     `json:"value"`
	SensorT     SensorType  `json:"sensor_type"`
	CommandT    CommandType `json:"command_type"`
	GreaterThen bool        `json:"greater_then"`
}

func GetTriggers() []Trigger {
	var ts []Trigger
	database.DB.Find(&ts)
	return ts
}

func (c *Command) CreateCommand() error {
	record := database.DB.Create(&c)
	if record.Error != nil {
		return record.Error
	}
	return nil
}
func (t *Trigger) CreateTrigger() error {
	record := database.DB.Create(&t)
	if record.Error != nil {
		return record.Error
	}
	return nil
}
func (c *TimedCommand) CreateTimedCommand() error {
	record := database.DB.Create(&c)
	if record.Error != nil {
		return record.Error
	}
	return nil
}
