package service

import (
	"context"

	"github.com/gauss2302/testcommm/product/internal/domain/entity"
	"github.com/gauss2302/testcommm/product/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, name, description string, price float64, userID uint64) (*entity.Product, error) {
	product := &entity.Product{
		Name:        name,
		Description: description,
		Price:       price,
		UserID:      userID,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uint64) (*entity.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context, page, perPage int32) ([]*entity.Product, int64, error) {
	return s.productRepo.List(ctx, page, perPage)
}
