package service

import (
	"fmt"
	"gihub.com/Kratos40-sba/command-service/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	HostFormat = "tcp://%s:%s"
	QOS        = 1
)

var (
	connectionHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Println("Client is connected")
	}
	connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Println("Connection of Client lost caused by Error  : ", err)
	}
)

type MqttConnection struct {
	mqttClient mqtt.Client
}

func NewMqttConnection(host, port, clientName string) (mqttConnection *MqttConnection) {
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf(HostFormat, host, port))
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
	return mqttConnection
}

// IsClientConnected for the hello API
func (mqtt *MqttConnection) IsClientConnected() bool {
	connected := mqtt.mqttClient.IsConnected()
	if !connected {
		log.Println("MQTT client is not connected")
	}
	return connected
}

func (mqtt *MqttConnection) Publish(topicName string, command string) {
	token := mqtt.mqttClient.Publish(topicName, QOS, false, command)
	token.Wait()
	log.Printf("Publishing on Topic : %s  payload : %s ", topicName, command)
}

// todo subscribe to sensors
func onMessageReceived(conn *MqttConnection) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		switch msg.Topic() {
		case "esp/sensor":
			handleMessage(conn, msg)
		default:
			log.Println("Client is not subscribed to this Topic : ", msg.Topic())
		}
	}
}
func handleMessage(conn *MqttConnection, msg mqtt.Message) {
	p := strings.Split(string(msg.Payload()), "||")
	tem, _ := strconv.ParseFloat(p[0], 32)
	h, _ := strconv.ParseFloat(p[1], 32)
	s, _ := strconv.ParseFloat(p[2], 32)
	var ts []model.Trigger
	ts = model.GetTriggers()
	for _, t := range ts {
		switch t.SensorT {
		case model.HUM:
			if t.GreaterThen {
				if t.Value >= h {
					switch t.CommandT {
					case model.OFF:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","0")
					case model.ON:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","1")
					}
				}
			} else {
				if t.Value < h {
					switch t.CommandT {
					case model.OFF:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","0")
					case model.ON:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","1")
					}

				}
			}
		case model.TEMP:
			if t.GreaterThen {
				if t.Value >= tem {
					switch t.CommandT {
					case model.OFF:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","0")

					case model.ON:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","1")
					}
				}
			} else {
				if t.Value < tem {
					switch t.CommandT {
					case model.OFF:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","0")
					case model.ON:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","1")
					}

				}
			}
		case model.SOIL:
			if t.GreaterThen {
				if t.Value >= s {
					switch t.CommandT {
					case model.OFF:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","0")
					case model.ON:
						sendEmail(t, t.Time)
						//payload:= "1"
						//conn.Publish("esp/action",payload)
					}
				}
			} else {
				if t.Value < s {
					switch t.CommandT {
					case model.OFF:
						sendEmail(t, t.Time)
						//conn.Publish("esp/action","0")
					case model.ON:
						sendEmail(t, t.Time)
						//payload:= "1"
						//conn.Publish("esp/action",payload)
					}
				}
			}
		}
		time.Sleep(10 * time.Second)
	}

}
func sendEmail(t model.Trigger, tt string) {
	s := []string{t.Email}
	sub := fmt.Sprintf("Trigger activated by %s", tt)
	message := fmt.Sprintf("ALERT ON SENSOR : %s \n VALUE : %f \n COMMAND : LED %s \n TIME & DATE : %s %s", t.SensorT, t.Value, t.CommandT, t.Time, t.Date)
	n := NewNotification("microservice.iot@gmail.com",
		"587", "smtp.gmail.com", []byte("Subject: "+sub+"\r\n"+
			"\r\n"+
			message+"\r\n"), s)
	err := n.SendEmail("iot22000")
	if err != nil {
		log.Println(err)
	}
}
func (mqtt *MqttConnection) Subscribe(topic string) {
	token := mqtt.mqttClient.Subscribe(topic, QOS, onMessageReceived(mqtt))
	token.Wait()
	log.Println("Subscribed to topic : ", topic)
}

// solve publishing
// solve email message formats
// solve deploying / testing
