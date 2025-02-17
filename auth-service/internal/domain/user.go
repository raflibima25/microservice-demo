package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id uint64) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint64) error
}

type AuthUseCase interface {
	Register(username, email, password string) (*User, string, error)
	Login(username, password string) (*User, string, error)
	ValidateToken(token string) (*User, error)
	Logout(token string) error
}

type TokenService interface {
	GenerateToken(userID uint64) (string, error)
	ValidateToken(token string) (uint64, error)
	BlacklistToken(token string) error
	IsTokenBlacklisted(token string) bool
}
