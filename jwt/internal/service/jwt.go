package service

import (
	"time"

	"github.com/SigmarWater/crm/jwt/model"
)

const (
	// В реальном проекте должны быть в переменных окружения
	accessTokenSecret  = "access-secret-key-very-long-and-secure"
	refreshTokenSecret = "refresh-secret-key-very-long-and-secure"
	accessTokenTTL     = 15 * time.Minute
	refreshTokenTTL    = 24 * time.Hour
)

// JWTService - сервис для работы с JWT токенами
type JWTService struct {
	users map[string]model.User // Хардкодные пользователи
}

// NewJWTService - создает новый JWT сервис
func NewJWTService() *JWTService {
	// Хардкодим 5 пользователей с хешами паролей (bcrypt)
	users := map[string]model.User{
		"admin": { //nolint:gosec // bcrypt hash for demo user
			ID:       1,
			Username: "admin",
			Password: "$2a$10$gyiozWTQ/GT6eboui8TaVO7ylNJXUyiHEmTwxom.zwQRoOj5JxSD2", // admin123
		},
		"user1": { //nolint:gosec // bcrypt hash for demo user
			ID:       2,
			Username: "user1",
			Password: "$2a$10$51k5c6cgO9bz3j86B2nL2O1REDi/Uh79GU42./00jiubmoIUyMUti", // password1
		},
		"user2": { //nolint:gosec // bcrypt hash for demo user
			ID:       3,
			Username: "user2",
			Password: "$2a$10$kR7cmWGYYHncFPcRFZRccOPKmuUKdOZ4o7L5NyA.oc2i0JORzJnyq", // password2
		},
		"john": { //nolint:gosec // bcrypt hash for demo user
			ID:       4,
			Username: "john",
			Password: "$2a$10$9n9uiJDbuUhfqrrXg/sj6.5veD4P3RL1y3HU//MxnhYDa7IWajFKm", // john123
		},
		"alice": { //nolint:gosec // bcrypt hash for demo user
			ID:       5,
			Username: "alice",
			Password: "$2a$10$bqqQ2pDjdCDlleIV/nyetOk2yIw.xAeez4MxhxQ8pwaBOm/XE5sau", // alice456
		},
	}

	return &JWTService{
		users: users,
	}
}
