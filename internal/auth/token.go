package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type UserToken struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
	TokenType string `json:"tokenType"`
}

func CreateToken(secretKey []byte, userID string) (*UserToken, error) {
	expiresIn := time.Now().Add(time.Hour * 24 * 7).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"exp": expiresIn,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &UserToken{Token: tokenString, ExpiresIn: expiresIn, TokenType: "Bearer"}, nil
}

func VerifyToken(secretKey []byte, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}
