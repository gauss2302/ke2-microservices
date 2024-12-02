package repository

import (
	"context"

	"github.com/gauss2302/testcommm/product/internal/domain/entity"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepository) GetByID(ctx context.Context, id uint64) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) List(ctx context.Context, page, perPage int32) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64

	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := r.db.WithContext(ctx).Offset(int(offset)).Limit(int(perPage)).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
