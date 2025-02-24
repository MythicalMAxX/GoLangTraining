package jwt

import (
    "errors"
    "mypackage/internal/models"
    "time"

    jwt "github.com/golang-jwt/jwt/v4"
)

var (
    ErrInvalidToken = errors.New("invalid token")
    secretKey       []byte
)

func SetSecretKey(key []byte) {
    secretKey = key
}

type Claims struct {
    UserID string     `json:"user_id"`
    Role   models.Role `json:"role"`
    jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
    return c.RegisteredClaims.Valid()
}

func GenerateToken(userID string, role models.Role) (string, error) {
    claims := &Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "Auth Service",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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