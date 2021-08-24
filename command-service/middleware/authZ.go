package middleware

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type JwtWrapper struct {
	SecurityKey string
	Issuer string
}
type Claim struct {
	Email string
	jwt.StandardClaims
}

func(j *JwtWrapper) ValidateToken(signedToken string) (c *Claim,err error) {
	token , err := jwt.ParseWithClaims(
		signedToken,
		&Claim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecurityKey),nil
		},
	)
	if err != nil {
		return
	}
	claims , ok := token.Claims.(*Claim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("jwt is expired")
		return
	}
	return claims , nil
}

func Authentication() gin.HandlerFunc  {
	return func(context *gin.Context) {
		clientToken := context.Request.Header.Get("Authorization")
		if clientToken == ""{
			context.JSON(403,"No Authorization header provided")
			context.Abort()
			return
		}
		extractedToken := strings.Split(clientToken,"Bearer ")
		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		}else {
			context.JSON(400,"Incorrect format for jwt")
			context.Abort()
			return
		}
		jwtWrapper :=JwtWrapper{
			SecurityKey: "secretKey",
			Issuer: "Auth-service",
		}
		claims , err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			context.JSON(401,err.Error())
			context.Abort()
			return
		}
		context.Set("email",claims.Email)
		context.Next()

	}
}

