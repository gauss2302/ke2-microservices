package service

import (
	"context"
	"errors"

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

func (s *ProductService) UpdateProduct(ctx context.Context, id uint64, userID uint64, name, description string, price float64) (*entity.Product, error) {
	// Check if product belongs to user
	belongs, err := s.productRepo.BelongsToUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, errors.New("product does not belong to user")
	}

	product := &entity.Product{
		Name:        name,
		Description: description,
		Price:       price,
	}

	if err := s.productRepo.Update(ctx, id, product); err != nil {
		return nil, err
	}

	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint64, userID uint64) error {
	// Check if product belongs to user
	belongs, err := s.productRepo.BelongsToUser(ctx, id, userID)
	if err != nil {
		return err
	}
	if !belongs {
		return errors.New("product does not belong to user")
	}

	return s.productRepo.Delete(ctx, id)
}

func (s *ProductService) ListUserProducts(ctx context.Context, userID uint64, page, perPage int32) ([]*entity.Product, int64, error) {
	return s.productRepo.ListByUserID(ctx, userID, page, perPage)
}
