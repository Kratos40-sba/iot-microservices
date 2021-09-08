package main

import (
	"fmt"
	"gihub.com/Kratos40-sba/command-service/config"
	"gihub.com/Kratos40-sba/command-service/controllers"
	"gihub.com/Kratos40-sba/command-service/database"
	"gihub.com/Kratos40-sba/command-service/middleware"
	"gihub.com/Kratos40-sba/command-service/model"
	"gihub.com/Kratos40-sba/command-service/service"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

const (
	HttpServer     = "HTTP_SERVER" // :8080
	DbName         = "POSTGRES_DB_NAME"
	HostName       = "POSTGRES_HOST"
	DbUser         = "POSTGRES_USER_NAME"
	DbPort         = "POSTGRES_PORT" // 5432
	MqttClientName = "MQTT_CLINT_NAME"
	MqttHost       = "MQTT_HOST" // 192.168.1.5
	MqttPort       = "MQTT_PORT" // 1883
	ActionTopic    = "esp/action"
	DhtTopic       = "esp/sensor"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error while loading config : ", err)
	}
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// setup DATABASE
	db, err := database.InitPostgres(os.Getenv(HostName), os.Getenv(DbName), os.Getenv(DbUser), os.Getenv(DbPort))
	if err != nil {
		log.Fatalln("err ", err)
	}
	err = db.AutoMigrate(&model.Command{})
	err = db.AutoMigrate(&model.Trigger{})
	err = db.AutoMigrate(&model.TimedCommand{})
	if err != nil {
		log.Fatalln("Error while migrating model : ", err)
	}
	mqttClient := service.NewMqttConnection(os.Getenv(MqttHost), os.Getenv(MqttPort), os.Getenv(MqttClientName))
	go mqttClient.Subscribe(DhtTopic)
	commandAPI := r.Group("/api/v1/command").Use(middleware.Authentication())
	{
		commandAPI.POST("/add", controllers.AddCommand(mqttClient, ActionTopic))
		commandAPI.POST("/timed", controllers.AddTimedCommand(mqttClient, ActionTopic))
		commandAPI.GET("/all", controllers.GetAllCommands)
		commandAPI.POST("/trigger", controllers.AddTrigger)
		commandAPI.DELETE("/trigger/:id", controllers.DeleteTrigger)

	}
	log.Fatalln(r.Run(fmt.Sprintf(":%s", os.Getenv(HttpServer))))
}
