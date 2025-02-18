package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Stock       int32          `gorm:"not null" json:"stock"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductRepository interface {
	Create(product *Product) error
	FindByID(id uint64) (*Product, error)
	Update(product *Product) error
	Delete(id uint64) error
	List(page, limit int32, search string) ([]Product, int64, error)
}

type ProductUseCase interface {
	Create(name, description string, price float64, stock int32) (*Product, error)
	GetByID(id uint64) (*Product, error)
	Update(id uint64, name, description string, price float64, stock int32) (*Product, error)
	Delete(id uint64) error
	List(page, limit int32, search string) ([]Product, int64, error)
}
