package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsmonday/youtube/configs"
)

func GenerateJWTToken(userId, name, email string) (string, error) {
	env, _ := configs.Config()
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    userId,
			"name":  name,
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24),
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
		return nil, err
	}
	return token, nil
}
