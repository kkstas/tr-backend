package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

type UserToken struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
	TokenType string `json:"tokenType"`
}

func CreateToken(userID string) (*UserToken, error) {
	expiresIn := time.Now().Add(time.Hour * 24 * 7).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  userID,
			"exp": expiresIn,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &UserToken{Token: tokenString, ExpiresIn: expiresIn, TokenType: "Bearer"}, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
