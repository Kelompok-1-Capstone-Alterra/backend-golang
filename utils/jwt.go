package utils

import (
	"errors"

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

func GetUserIDFromToken(tokenString string) (uint, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		// Return the secret key used for signing
		return []byte(constant.SECRET_JWT), nil
	})

	if err != nil {
		return 0, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the user_id claim
		if userID, ok := claims["user_id"].(float64); ok {
			return uint(userID), nil
		}
	}

	return 0, errors.New("unable to retrieve user_id from token")
}
