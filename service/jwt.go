package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaim struct {
	ID string `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var jwtSecret = []byte(getJwtSecret())

func getJwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	//ERROR: Check this and return error
	if secret == "" {
		return "testing"
	}
	return secret
}

func JwtGenerate(ctx context.Context, userID string, userName string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		ID: userID,
		Name: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	})

	token, err := t.SignedString(jwtSecret)
	if err != nil {
		fmt.Println()
		return "", err
	}
	return token, nil
}

func JwtValidate(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &JwtCustomClaim{}, func(t *jwt.Token)(interface{}, error){
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("problem signing method")
		}
		return jwtSecret, nil
	})
}
