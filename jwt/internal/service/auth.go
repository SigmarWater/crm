package service

import (
	"github.com/SigmarWater/crm/jwt/model"
	"time"

	"golang.org/x/crypto/bcrypt"

	authErrors "github.com/SigmarWater/crm/jwt/errors"
)

// Login - аутентификация пользователя
func (s *JWTService) Login(username, password string) (*model.TokenPair, error) {
	user, exists := s.users[username]
	if !exists {
		return nil, authErrors.ErrInvalidCredentials
	}

	// Проверяем пароль с хешем
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, authErrors.ErrInvalidCredentials
	}

	return s.generateTokenPair(user)
}

// GetAccessToken - получение нового access токена по refresh токену
func (s *JWTService) GetAccessToken(refreshToken string) (string, time.Time, error) {
	claims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	user, exists := s.users[claims.Username]
	if !exists {
		return "", time.Time{}, authErrors.ErrInvalidToken
	}

	return s.generateAccessToken(user)
}

// GetRefreshToken - получение нового refresh токена
func (s *JWTService) GetRefreshToken(refreshToken string) (string, time.Time, error) {
	claims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	user, exists := s.users[claims.Username]
	if !exists {
		return "", time.Time{}, authErrors.ErrInvalidToken
	}

	return s.generateRefreshToken(user)
}
