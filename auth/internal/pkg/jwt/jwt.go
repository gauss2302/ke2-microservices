package jwt

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTMaker struct {
	privateKey string
}

func NewJWTMaker(privateKey string) *JWTMaker {
	return &JWTMaker{privateKey: privateKey}
}

func (maker *JWTMaker) CreateTokenPair(userID uint64) (*TokenPair, error) {
	// Создаем access token (короткоживущий)
	accessToken, err := maker.createToken(userID, time.Hour*24)
	if err != nil {
		return nil, err
	}

	// Создаем refresh token (долгоживущий)
	refreshToken, err := maker.createToken(userID, time.Hour*24*7)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (maker *JWTMaker) createToken(userID uint64, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.privateKey))
}

type Claims struct {
	jwt.StandardClaims
	UserID uint64 `json:"user_id"`
}

func (c *Claims) GetUserID() uint64 {
	return c.UserID
}

func (maker *JWTMaker) VerifyToken(tokenString string) (uint64, error) {
	log.Printf("Verifying token: %s", tokenString)

	// Используем MapClaims вместо структуры Claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.privateKey), nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v", err)
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Получаем user_id из MapClaims
		if userID, ok := claims["user_id"].(float64); ok {
			return uint64(userID), nil
		}
		return 0, errors.New("user_id not found in token")
	}

	return 0, errors.New("invalid token")
}
