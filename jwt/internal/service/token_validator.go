package service

import (
	jwtErrors "github.com/SigmarWater/crm/jwt/errors"
	"github.com/SigmarWater/crm/jwt/model"
	"github.com/golang-jwt/jwt/v5"
)

// validateRefreshToken - валидирует refresh токен
func (s *JWTService) validateRefreshToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwtErrors.ErrInvalidToken
		}
		return []byte(refreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, jwtErrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwtErrors.ErrInvalidToken
	}

	// Проверяем тип токена
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, jwtErrors.ErrInvalidToken
	}

	userID, ok := claims["user_id"].(float64) // JWT парсит числа как float64
	if !ok {
		return nil, jwtErrors.ErrInvalidToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, jwtErrors.ErrInvalidToken
	}

	return &model.Claims{
		UserID:   int64(userID),
		Username: username,
	}, nil
}
