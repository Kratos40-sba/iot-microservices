package security

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}
type JwtClaim struct {
	Email string
	jwt.StandardClaims
}

func(j *JwtWrapper) GenerateToken (email string ) (signedToken string , err error){

	claims := &JwtClaim{
		Email: email,
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour*time.Duration(j.ExpirationHours)).Unix(),
			Issuer: j.Issuer ,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signedToken ,err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "",err
	}
	return signedToken,nil
}
func(j *JwtWrapper) ValidateToken(signedToken string)(claims *JwtClaim,err error){
	token , err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return nil , err
	}
	claims , ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil , err

	}
	if claims.ExpiresAt <= time.Now().Local().Unix() {
		err = errors.New("jwt is expired")
		return nil , err
	}
	return claims , nil

}
