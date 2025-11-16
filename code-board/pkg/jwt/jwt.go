package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/ak-repo/code-board/pkg/config"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Email  string `json:"email"`
	Role   string `json:"role"`
	UserID uint   `json:"userID"`
	Type   string `json:"type"` // 👈 to differentiate access vs refresh
	jwt.RegisteredClaims
}

type JWT struct {
	Email  string `json:"email"`
	Role   string `json:"role"`
	UserID uint   `json:"userID"`
}

func AccessTokenGenerator(data JWT, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email:  data.Email,
		UserID: data.UserID,
		Role:   data.Role,
		Type:   "access", // 👈 helpful for refresh validation
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(cfg.JWT.AccessExpireMinutes))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.JWT.Issuer,
		},
	})
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func RefreshTokenGenerator(data JWT, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email:  data.Email,
		UserID: data.UserID,
		Role:   data.Role,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cfg.JWT.RefreshExpireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.JWT.Issuer,
		},
	})
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func TokenValidator(tokenStr string, cfg *config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token expired")
			}
		}
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}
	if !token.Valid {
		return nil, errors.New("token invalid")
	}

	return claims, nil
}
