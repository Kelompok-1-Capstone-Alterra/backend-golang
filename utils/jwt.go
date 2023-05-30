package utils

import (
	"github.com/agriplant/constant"
	"github.com/dgrijalva/jwt-go"
)

func CreateTokenAdmin(adminId uint, name string) (string, error) {
	// create the claims
	claims := jwt.MapClaims{}
	claims["admin_id"] = adminId
	claims["name"] = name
	claims["role"] = "admin"

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(constant.SECRET_JWT))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateTokenUser(userId uint, name string) (string, error) {
	// create the claims
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["name"] = name
	claims["role"] = "user"

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(constant.SECRET_JWT))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
