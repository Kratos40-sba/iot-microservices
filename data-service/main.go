package main

import (
	"fmt"
	"gihub.com/Kratos40-sba/data-service/config"
	"gihub.com/Kratos40-sba/data-service/controllers"
	"gihub.com/Kratos40-sba/data-service/database"
	"gihub.com/Kratos40-sba/data-service/middleware"
	"gihub.com/Kratos40-sba/data-service/model"
	"gihub.com/Kratos40-sba/data-service/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

const (
	HttpServer             = "HTTP_SERVER" // :8080
	DbName                 = "POSTGRES_DB_NAME"
	DbUser                 = "POSTGRES_USER_NAME"
	DbPort                 = "POSTGRES_PORT" // 5432
	MqttClientName         = "MQTT_CLINT_NAME"
	MqttHost               = "MQTT_HOST" // 192.168.1.6
	MqttPort               = "MQTT_PORT" // 1883
	DhtTopic               = "esp/sensor"
	_NotificationTopic     = "esp/notification"
)
func main()  {
  // Load CONFIG from .env
  err := config.LoadConfig()
  if err != nil {
    log.Fatalln("Error while loading config : ",err)
  }
	router := gin.Default()
	// setup DATABASE
	db , err := database.InitPostgres(os.Getenv(DbName),os.Getenv(DbUser),os.Getenv(DbPort))
	if err != nil {
		log.Fatalln("err ", err)

	}
	err = db.AutoMigrate(&model.Event{})
	if err != nil {
		log.Fatalln("Error while migrating model : ",err)

	}
	// setup MQTT CLIENT
	mqttClient := service.NewMqttConnection(os.Getenv(MqttHost),os.Getenv(MqttPort),os.Getenv(MqttClientName))
	mqttClient.Subscribe(DhtTopic)
	// setup ROUTES
	sensorAPI  := router.Group("/api/v1/data").Use(middleware.Authentication())
	{
		sensorAPI.GET("/mqttHealth", controllers.Hello(mqttClient))
		// todo get only 25% for performance
		sensorAPI.GET("/all",controllers.GetAllEvents)
		sensorAPI.GET("/last",controllers.GetLastEvents)
		sensorAPI.GET("/byDate",controllers.GetEventByDate)
		// todo add threshold controllers
	}
	router.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound,gin.H{"Message":"No route for this url "})
	})
	log.Fatalln(router.Run(fmt.Sprintf(":%s",os.Getenv(HttpServer))))

}