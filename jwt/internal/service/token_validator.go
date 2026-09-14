package service

import (
	"github.com/golang-jwt/jwt/v5"

	jwtErrors "github.com/SigmarWater/crm/jwt/errors"
	"github.com/SigmarWater/crm/jwt/model"
)

// validateRefreshToken - валидирует refresh токен
func (s *JWTService) validateRefreshToken(tokenString string) (*model.Claims, error) {
	// Parse парсит токен
	// tokenString - сам токен
	// func - какой секрет использовать для проверки подписи

	// callback вызывается после того, как библиотека разобрала JWT и увидела его alg (алгоритм подписи).

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяется алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwtErrors.ErrInvalidToken
		}
		// Возвращаем секрет, которым была создана подпись JWT
		return []byte(refreshTokenSecret), nil
	})

	// err != nil - ошибка при парсинке
	// или токен невалидный
	if err != nil || !token.Valid {
		return nil, jwtErrors.ErrInvalidToken
	}

	// Получаем payload, если токен валиден
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
