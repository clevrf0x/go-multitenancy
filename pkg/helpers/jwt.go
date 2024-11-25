package helpers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

func GetJWTAuth() (*jwtauth.JWTAuth, error) {
	secretKey, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable is not set")
	}
	return jwtauth.New("HS256", []byte(secretKey), nil), nil
}

func getAccessExpHours() (int, error) {
	accessExpStr, ok := os.LookupEnv("JWT_ACCESS_EXP_HOURS")
	if !ok {
		return 0, fmt.Errorf("JWT_ACCESS_EXP_HOURS environment variable is not set")
	}
	return strconv.Atoi(accessExpStr)
}

func getRefreshExpHours() (int, error) {
	refreshExpStr, ok := os.LookupEnv("JWT_REFRESH_EXP_HOURS")
	if !ok {
		return 0, fmt.Errorf("JWT_REFRESH_EXP_HOURS environment variable is not set")
	}
	return strconv.Atoi(refreshExpStr)
}

func GenerateAccessToken(userID, tenantID string) (string, error) {
	TokenAuth, err := GetJWTAuth()
	if err != nil {
		return "", err
	}

	accessExpHours, err := getAccessExpHours()
	if err != nil {
		return "", err
	}

	claims := map[string]interface{}{
		"user_id":   userID,
		"tenant_id": tenantID,
		"exp":       time.Now().Add(time.Hour * time.Duration(accessExpHours)).Unix(),
		"iat":       time.Now().Unix(),
		"nbf":       time.Now().Unix(),
		"type":      "access",
	}

	_, tokenString, err := TokenAuth.Encode(claims)
	return tokenString, err
}

func GenerateRefreshToken(userID, tenantID string) (string, error) {
	TokenAuth, err := GetJWTAuth()
	if err != nil {
		return "", err
	}

	refreshExpHours, err := getRefreshExpHours()
	if err != nil {
		return "", err
	}

	claims := map[string]interface{}{
		"user_id":   userID,
		"tenant_id": tenantID,
		"exp":       time.Now().Add(time.Hour * time.Duration(refreshExpHours)).Unix(),
		"iat":       time.Now().Unix(),
		"nbf":       time.Now().Unix(),
		"type":      "refresh",
	}

	_, tokenString, err := TokenAuth.Encode(claims)
	return tokenString, err
}

func GetTokenClaims(r context.Context) (map[string]interface{}, error) {

	_, claims, err := jwtauth.FromContext(r)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
