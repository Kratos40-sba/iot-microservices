package main

import (
	"fmt"
	"gihub.com/Kratos40-sba/visualisation-service/config"
	"gihub.com/Kratos40-sba/visualisation-service/controllers"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var msg mqtt.Message

const (
	HttpServer     = "HTTP_SERVER" // :9292
	MqttClientName = "MQTT_CLINT_NAME"
	MqttHost       = "MQTT_HOST" // 192.168.1.6
	MqttPort       = "MQTT_PORT" // 1883
	DhtTopic       = "esp/sensor"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	mqttClient := controllers.NewMqttConnection(os.Getenv(MqttHost), os.Getenv(MqttPort), os.Getenv(MqttClientName))

	mqttClient.Subscribe(DhtTopic)
	router.GET("/ws", controllers.WsHandler)

	visualisationAPI := router.Group("/api/v1/visual")
	{
		visualisationAPI.GET("/", controllers.Home)
	}
	log.Fatalln(router.Run(fmt.Sprintf(":%s", os.Getenv(HttpServer))))

}
