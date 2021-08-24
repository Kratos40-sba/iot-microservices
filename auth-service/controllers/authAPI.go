package controllers

import (
	"gihub.com/Kratos40-sba/auth-service/database"
	"gihub.com/Kratos40-sba/auth-service/model"
	"gihub.com/Kratos40-sba/auth-service/security"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

// LoginPayload login body
type LoginPayload struct {
	Email     string `json:"email"`
	RfidValue string `json:"rfid_value"`
}
// LoginResponse token response
type LoginResponse struct {
	Resp string `json:"response"`
}

// Signup but in our iot app we dont need it
func Signup(c *gin.Context)  {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(400,gin.H{
			"message":"invalid json",
		})
		c.Abort()
		return
	}
	err = user.HashRfidNumber(user.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500,gin.H{
			"message":"error hashing password",
		})
		c.Abort()
		return
	}
	err = user.CreateUserRecord()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"message": "error hashing password",
		})
		c.Abort()
		return
	}
	c.JSON(200,user)
}
// Login logs user in
func Login(c *gin.Context)  {
	var payload LoginPayload
	var user model.User
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(400,gin.H{
			"message":"invalid json",
		})
		c.Abort()
		return
	}
	// in our case we replace email with rfid number
	result := database.DB.Where("email = ?",payload.Email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(401,gin.H{
			"message":"invalid user credentials",
		})
		c.Abort()
		return
	}
	err = user.CheckPassword(payload.RfidValue)
	if err != nil {
		log.Println(err)
		c.JSON(401,gin.H{
			"message":"invalid RfidValue",
		})
		c.Abort()
		return
	}
	jwtWatpper := security.JwtWrapper{
		SecretKey: "secretKey",
		Issuer: "Auth-service",
		ExpirationHours: 24 ,
	}
	signedToken , err := jwtWatpper.GenerateToken(user.Email)
	if err != nil {
		log.Println(err)
		c.JSON(500,gin.H{
			"message":"error signing token",
		})
		c.Abort()
		return
	}

	tokenResponse := LoginResponse{
		Resp: signedToken,
	}
	c.JSON(200,tokenResponse)
	return
}
func Logout(c *gin.Context)  {
    c.Set("Authorization","")
	c.JSON(200,gin.H{
		"message":"User sign out ",
	})
}