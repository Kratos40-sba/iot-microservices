package model

import "time"

type Event struct {
	Time time.Time `json:"time"`
	Temp float64 `json:"temperature"`
	Hum  float64 `json:"humidity"`
	Soil float64 `json:"soil_percentage"`
}
