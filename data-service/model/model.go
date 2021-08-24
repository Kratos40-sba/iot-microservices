package model

type Event struct {

	Time         string    `json:"time"`
	Date         string    `json:"date"`
	Temperature  float64   `json:"temperature"`
	Humidity     float64   `json:"humidity"`
	SoilMoisture float64   `json:"soil_moisture_percentage"`
}

