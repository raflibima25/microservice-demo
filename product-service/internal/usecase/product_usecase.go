package usecase

import (
	"errors"
	"product-service/internal/domain"
)

type productUseCase struct {
	productRepo domain.ProductRepository
}

func NewProductUseCase(productRepo domain.ProductRepository) domain.ProductUseCase {
	return &productUseCase{productRepo: productRepo}
}

func (u *productUseCase) Create(name, description string, price float64, stock int32) (*domain.Product, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	product := &domain.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	err := u.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUseCase) GetByID(id uint64) (*domain.Product, error) {
	return u.productRepo.FindByID(id)
}

func (u *productUseCase) Update(id uint64, name, description string, price float64, stock int32) (*domain.Product, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		product.Name = name
	}
	if description != "" {
		product.Description = description
	}
	if price > 0 {
		product.Price = price
	}
	if stock >= 0 {
		product.Stock = stock
	}

	err = u.productRepo.Update(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUseCase) Delete(id uint64) error {
	_, err := u.productRepo.FindByID(id)
	if err != nil {
		return err
	}

	return u.productRepo.Delete(id)
}

func (u *productUseCase) List(page, limit int32, search string) ([]domain.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return u.productRepo.List(page, limit, search)
}
