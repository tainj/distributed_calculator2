package auth

import (
    "errors"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type Config struct {
    SecretKey      string
    ExpireDuration time.Duration
    Issuer         string
}

type JWTService struct {
    config Config
}

func NewJWTService(config Config) *JWTService {
    return &JWTService{config: config}
}

func (s *JWTService) GenerateToken(userID string) (string, error) {
    if _, err := uuid.Parse(userID); err != nil {
        return "", errors.New("invalid user ID")
    }

    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.ExpireDuration)),
            Issuer:    s.config.Issuer,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.config.SecretKey))
}

func (s *JWTService) ParseToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(s.config.SecretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}