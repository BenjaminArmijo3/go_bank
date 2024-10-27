package token

import (
	"fmt"
	"log"
	"time"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/BenjaminArmijo3/bank/internal/config"
	"github.com/golang-jwt/jwt"
)

type JWTToken struct {
	config *config.Config
}

type jwtClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
	Exp    int64 `json:"exp"`
}

func NewJWTToken(config *config.Config) *JWTToken {
	return &JWTToken{
		config: config,
	}
}

func (j *JWTToken) CreateToken(user_id int64) (string, error) {
	claims := jwtClaim{
		UserID: user_id,
		Exp:    time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"user_id": claims.UserID,
	// 	"exp":     claims.Exp,
	// })

	tokenString, err := token.SignedString([]byte(j.config.SigningKey))
	if err != nil {
		return "", err
	}
	return string(tokenString), nil
}

func (j *JWTToken) VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid authentication token")
		}
		return []byte(j.config.SigningKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid authentication token")
	}

	claims, ok := token.Claims.(*jwtClaim)

	if !ok {
		return 0, fmt.Errorf("invalid authentication token")
	}

	if claims.Exp < time.Now().Unix() {
		return 0, fmt.Errorf("token has expired")
	}
	log.Println(claims)

	return claims.UserID, nil
}
