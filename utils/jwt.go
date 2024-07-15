package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("secretForBlogApp") //change this secret key as per your project

func GenerateJWT(userId uint) (string, error) {
	// Claim setup
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}
	return signedToken, nil
}
