package main

import (
	"fmt"
	"gihub.com/Kratos40-sba/auth-service/config"
	"gihub.com/Kratos40-sba/auth-service/controllers"
	"gihub.com/Kratos40-sba/auth-service/database"
	"gihub.com/Kratos40-sba/auth-service/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

const (
	userName = "DB_USER_NAME"
	dbName   = "DB_NAME"
	portDB   = "DB_PORT"
	authPort = "AUTH_PORT"
	HostName = "POSTGRES_HOST"
)

// todo cookies and stuff
func setupRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	// todo add refresh tokens
	api := r.Group("/api/v1/auth")
	{
		//localhost:9090/api/v1/auth/login
		//                          /signup
		// 8585/login
		api.POST("/login", controllers.Login)
		api.POST("/signup", controllers.Signup)
		api.GET("/logout", controllers.Logout)
	}
	return r
}
func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := database.InitDB(os.Getenv(HostName), os.Getenv(userName), os.Getenv(dbName), os.Getenv(portDB))
	if err != nil {
		log.Fatalln("Could not init database", err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalln("Error while migrating the user table : ", err)
	}
	// todo initialisation of redis session
	router := setupRouter()
	log.Fatalln(router.Run(fmt.Sprintf(":%s", os.Getenv(authPort))))

}
