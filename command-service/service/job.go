package service

import "time"

func (mqtt *MqttConnection) AddTimedPub(d time.Duration, topic string) {
	mqtt.Publish(topic, "1")
	time.AfterFunc(d*time.Second, func() {
		mqtt.Publish(topic, "0")
	})
}
