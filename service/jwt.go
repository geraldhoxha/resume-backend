package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/geraldhoxha/resume-backend/graph/model"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaim struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// jwt.StandardClaims
	jwt.RegisteredClaims
}
type JwtRefreshToken struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// jwt.StandardClaims
	jwt.RegisteredClaims
}
type CustomClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// jwt.StandardClaims
	jwt.RegisteredClaims
}

var (
	jwtSecret        = []byte(getJwtSecret("ACCESS_JWT"))
	refreshJwtSecret = []byte(getJwtSecret("REFRESH_JWT"))
)

func getJwtSecret(secret_for string) string {
	secret := os.Getenv(secret_for)
	// ERROR: Check this and return error
	if secret == "" {
		return "testing"
	}
	return secret
}

func JwtGenerate(ctx context.Context, userID string, userName string, userEmail string) (*model.JwtToken, error) {
	a_Token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		ID:    userID,
		Name:  userName,
		Email: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	r_Token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtRefreshToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	accessToken, err := a_Token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	refreshToken, err := r_Token.SignedString(refreshJwtSecret)
	if err != nil {
		return nil, err
	}

	return &model.JwtToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func JwtValidate(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &JwtCustomClaim{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("problem signing method")
		}
		return jwtSecret, nil
	})
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Wrong body", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(request.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing problem")
		}
		return refreshJwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Wth", http.StatusUnauthorized)
		return
	}
	user := claims["id"].(string)
	name := claims["name"].(string)
	email := claims["email"].(string)

	tokenPair, err := JwtGenerate(ctx, user, name, email)
	if err != nil {
		http.Error(w, "SOmething went wrong", http.StatusNotFound)
		return
	}

	response := map[string]*model.JwtToken{
		"response": tokenPair,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}
