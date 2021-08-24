package service

import (
	"fmt"
	"gihub.com/Kratos40-sba/data-service/database"
	"gihub.com/Kratos40-sba/data-service/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	HostFormat     = "tcp://%s:%s"
	QOS            = 1
)
var (

	connectionHandler mqtt.OnConnectHandler =func(client mqtt.Client) {
		log.Println("Client is connected")}
	connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Println("Connection of Client lost caused by Error  : ", err)
	}

)
type MqttConnection struct {
	mqttClient mqtt.Client
}

func NewMqttConnection (host , port , clientName string) (mqttConnection *MqttConnection)  {
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf(HostFormat,host,port))
	options.SetClientID(clientName)
	options.AutoReconnect = true
	//options.WillEnabled =true
	options.OnConnect = connectionHandler
	options.OnConnectionLost = connectionLostHandler
	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln("Connection Problem :", token.Error())
	}
	mqttConnection = &MqttConnection{client}
	return  mqttConnection
}

// IsClientConnected for the hello API
func (conn *MqttConnection) IsClientConnected() bool {
	connected := conn.mqttClient.IsConnected()
	if !connected {
		log.Println("MQTT client is not connected")
	}
	return connected
}
// onMessageReceived  handle message when comes to the client from broker
func onMessageReceived() func(client mqtt.Client , msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		switch msg.Topic() {
		case "esp/sensor":
			// get temp and humidity and insert in db table dht
			// add mutex
			handleDhtMessage(client ,msg)
			// service.Process(msg)
			// do some commands
		case "esp/rfid":
		// forward/publish  the message to the authentication service
		default :
			log.Println("Client is not subscribed to this Topic : ", msg.Topic())

		}
	}
}
func handleDhtMessage(client mqtt.Client , msg mqtt.Message)  {

	event := model.Event{Time: time.Now().Format("15:04:05") , Date: time.Now().Format("2006-01-02" ) }
	p := strings.Split(string(msg.Payload()), "||")
	event.Temperature, _ = strconv.ParseFloat(p[0], 32)
	event.Humidity, _ = strconv.ParseFloat(p[1], 32)
	event.SoilMoisture , _ = strconv.ParseFloat(p[2],32)
	log.Println(event)
	database.DB.Create(event)
}
func (conn *MqttConnection) Subscribe(topic string) {
	token := conn.mqttClient.Subscribe(topic, QOS, onMessageReceived())
	token.Wait()
	log.Println("Subscribed to topic : ", topic)
}
