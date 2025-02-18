package repository

import (
	"product-service/internal/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint64) (*domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepository) List(page, limit int32, search string) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.Model(&domain.Product{})

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// calculate offset
	offset := (page - 1) * limit

	err = query.Offset(int(offset)).Limit(int(limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
