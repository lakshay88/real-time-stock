package database

import (
	"github.com/lakshay88/real-time-stock/internal/user/models"
)

type Database interface {
	// Close() error
	CreateUser(user models.User) error
	GetUserList(user []models.User) ([]models.User, error)
	GetUserByEmail(email string) (models.User, error)
}
