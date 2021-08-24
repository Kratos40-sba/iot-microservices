package controllers

import (
	"fmt"
	"gihub.com/Kratos40-sba/visualisation-service/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)
// WebSocket part
type WsConnection struct {
	*websocket.Conn
}
var eventChannel = make(chan model.Event)
var up = websocket.Upgrader{
	ReadBufferSize: 1024 ,
	WriteBufferSize: 1024 ,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(x *gin.Context) {
	socket, err := up.Upgrade(x.Writer, x.Request, nil)
	if err != nil {

		log.Println(">>> here: broken pipe <<<")
		log.Println(err)
	}
	//var event model.Event
    go func() {
		for {
			select {
			case event := <-eventChannel:
				err = socket.WriteJSON(event)
				if err != nil {
					//log.Println("here")
					log.Println(err)
				}
			}

		}
	}()

}
	// MQTT part
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
func parseDhtMessage(msg mqtt.Message) model.Event {
	event := model.Event{Time: time.Now() }
	p := strings.Split(string(msg.Payload()), "||")
	event.Temp, _ = strconv.ParseFloat(p[0], 32)
	event.Hum,  _ = strconv.ParseFloat(p[1], 32)
	event.Soil, _ = strconv.ParseFloat(p[2],32)
	return event

}
func (conn *MqttConnection) Subscribe(topic string) {
	token := conn.mqttClient.Subscribe(topic, QOS, onMessageReceived())
	token.Wait()
	log.Println("Subscribed to topic : ", topic)
}

func onMessageReceived() func(client mqtt.Client,msg mqtt.Message) {
	return func(client mqtt.Client ,msg mqtt.Message) {
		switch msg.Topic() {
		case "esp/sensor":
				e:=  parseDhtMessage(msg)
				// blocking instruction
				go func() {eventChannel <- e}()
				//log.Println(e)
		default:
			log.Println("This topic is unknown : ",msg.Topic())
		}
	}
}