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

func (r *ProductRepository) Update(ctx context.Context, id uint64, product *entity.Product) error {
	return r.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Updates(product).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

func (r *ProductRepository) BelongsToUser(ctx context.Context, productID uint64, userID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Product{}).
		Where("id = ? AND user_id = ?", productID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ProductRepository) ListByUserID(ctx context.Context, userID uint64, page, perPage int32) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64

	if err := r.db.WithContext(ctx).Model(&entity.Product{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Offset(int(offset)).
		Limit(int(perPage)).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
