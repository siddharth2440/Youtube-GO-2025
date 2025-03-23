package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hPassword), nil
}

func VerifyPassword(password, hPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
