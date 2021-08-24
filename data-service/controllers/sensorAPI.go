package controllers

import (
	"gihub.com/Kratos40-sba/data-service/database"
	"gihub.com/Kratos40-sba/data-service/model"
	"gihub.com/Kratos40-sba/data-service/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)
func Hello(mqtt *service.MqttConnection) func(c *gin.Context) {
	return func(x *gin.Context) {
		x.JSON(http.StatusOK,gin.H{"message":"Hello from DHT-API",
			"Status of mqtt Client":mqtt.IsClientConnected()})
	}
}
func GetAllEvents(x *gin.Context)  {
	//db := x.MustGet("db").(*gorm.DB)
	var dht []model.Event
	database.DB.Find(&dht)
	x.JSON(http.StatusOK, gin.H{"data": dht})
}
func GetLastEvents(x *gin.Context)  {
	//db := x.MustGet("db").(*gorm.DB)
	k := x.Query("last")
	kk, _ := strconv.Atoi(k)
	if kk <= 0 {
		x.JSON(http.StatusBadRequest, gin.H{"message": "check path parameter"})
		return
	}
	var dht []model.Event
	today := time.Now().Format("2006-01-02")
	database.DB.Where("date=?",today).Limit(kk).Order("time desc").Find(&dht)
	x.JSON(http.StatusOK, gin.H{"data": dht})
}
func GetEventByDate(x *gin.Context)  {
	//db := x.MustGet("db").(*gorm.DB)
	date := x.Query("date")
	var dht []model.Event
	database.DB.Where("date = ?",date).Find(&dht)
	x.JSON(http.StatusOK,gin.H{"data":dht})
}
