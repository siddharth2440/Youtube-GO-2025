package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsmonday/youtube/configs"
)

func GenerateJWTToken(userId, name, email string) (string, error) {
	env, _ := configs.Config()
	fmt.Printf("\n%v - %v - %v\n", userId, name, email)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    userId,
			"name":  name,
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	tokenstring, err := token.SignedString([]byte(env.JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

func JWT_Verification(tokenString string) (*jwt.Token, error) {
	env, _ := configs.Config()
	secretkey := env.JWT_SECRET_KEY
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})

	if err != nil {
		fmt.Printf("\n%v\n", err)
		return nil, err
	}
	fmt.Printf("\n User Token After Verification  %v\n", token)
	return token, nil
}
