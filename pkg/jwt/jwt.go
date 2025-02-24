package jwt

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4" 
)

var (
	ErrInvalidToken = errors.New("invalid token")
	secretKey       []byte
)

// SetSecretKey sets the JWT secret key
func SetSecretKey(key []byte) {
	secretKey = key
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Valid implements the jwt.Claims interface
func (c *Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}

func GenerateToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Go Training App",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Changed to HS256
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}
