package controllers

import (
	"gihub.com/Kratos40-sba/command-service/database"
	"gihub.com/Kratos40-sba/command-service/model"
	"gihub.com/Kratos40-sba/command-service/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)


func AddCommand (mqtt *service.MqttConnection,topic string) func(x *gin.Context) {
	return func(x *gin.Context) {
		var command = model.Command{
			Time: time.Now().Format("15:04:05"),
			Date: time.Now().Format("2006-01-02"),
			Email: x.GetString("email"),
		}
		err := x.ShouldBindJSON(&command)
		if err != nil {
			log.Println(err.Error())
			x.JSON(http.StatusBadRequest , gin.H{
				"Error":"JSON bad format",
			})
			x.Abort()
			return
		}
	/*
		if command.CommandT != model.ON && command.CommandT != model.OFF {
			x.JSON(http.StatusBadRequest , gin.H{
				"Error":"JSON bad format use ON or OFF",
			})
			x.Abort()
			return
		}
	 */
		err = command.CreateCommand()
		if err != nil {
			log.Println(err.Error())
			x.JSON(http.StatusInternalServerError , gin.H{
				"Error":"Couldn't save the record in the database",
			})
			x.Abort()
			return
		}
		// todo make this in goroutine
		switch command.CommandT {
			case model.ON:
				mqtt.Publish(topic,"1")
			case model.OFF:
				mqtt.Publish(topic,"0")
			default:
				x.JSON(http.StatusBadRequest , gin.H{
					"Error":"JSON bad format use ON or OFF",
				})
				x.Abort()
				return
		}

		x.JSON(http.StatusCreated, command)
	}

}
func GetAllCommands(x *gin.Context)  {
	var commands []model.Command
	var email = x.GetString("email")
	rs := database.DB.Where("email= ?",email).Find(&commands)
	if rs.Error == gorm.ErrEmptySlice {
		log.Println("Empty database")
		x.JSON(http.StatusNotFound,gin.H{
			"message":"no records found",
		})
		x.Abort()
		return
	}
	x.JSON(http.StatusOK,commands)
}
func AddTrigger(x *gin.Context) {
	    var trigger = model.Trigger{
			Time: time.Now().Format("15:04:05"),
			Date: time.Now().Format("2006-01-02"),
			Email : x.GetString("email"),
		}
	err := x.ShouldBindJSON(&trigger)
	if err != nil {
		log.Println(err.Error())
		x.JSON(http.StatusBadRequest , gin.H{
			"Error":"JSON bad format",
		})
		x.Abort()
		return
	}
	// sensor type and action type in a switch statement
	switch  {
	case trigger.SensorT == model.HUM && trigger.CommandT == model.ON :
      trigger.CreateTrigger()
      x.JSON(200,gin.H{"message":trigger})
	case trigger.SensorT == model.HUM && trigger.CommandT == model.OFF:
		trigger.CreateTrigger()
		x.JSON(200,gin.H{"message":trigger})
	case trigger.SensorT == model.TEMP&& trigger.CommandT == model.ON :
		trigger.CreateTrigger()
		x.JSON(200,gin.H{"message":trigger})
	case trigger.SensorT == model.TEMP&& trigger.CommandT == model.OFF:
		trigger.CreateTrigger()
		x.JSON(200,gin.H{"message":trigger})
	case trigger.SensorT == model.SOIL&& trigger.CommandT == model.ON :
		trigger.CreateTrigger()
		x.JSON(200,gin.H{"message":trigger})
		case trigger.SensorT == model.SOIL&& trigger.CommandT == model.OFF:
			trigger.CreateTrigger()
			x.JSON(200,gin.H{"message":trigger})
	default:
			x.JSON(http.StatusBadRequest , gin.H{
				"Error":"JSON bad format use ON or OFF",
			})
			x.Abort()
			return
	}


		// get values to compare with "inference"
		// get email from context
		// add unsubscribe method to

}