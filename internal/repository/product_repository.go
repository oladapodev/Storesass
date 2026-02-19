package repository

import (
	"github.com/yourusername/storefront-saas-go/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository interface {
	ListActiveProducts(page, limit int) ([]domain.Product, int64, error)
	ListProductsByStore(storeID uint, page, limit int) ([]domain.Product, int64, error)
	SearchActiveProducts(query string, page, limit int) ([]domain.Product, int64, error)
	FindProductByID(id uint) (*domain.Product, error)
	CreateProduct(product *domain.Product) error
	UpdateProductDetails(product *domain.Product) error
	DeactivateProduct(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) ListActiveProducts(page, limit int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64
	offset := (page - 1) * limit

	if err := r.db.Model(&domain.Product{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("is_active = ?", true).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) ListProductsByStore(storeID uint, page, limit int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64
	offset := (page - 1) * limit

	if err := r.db.Model(&domain.Product{}).Where("store_id = ? AND is_active = ?", storeID, true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("store_id = ? AND is_active = ?", storeID, true).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) SearchActiveProducts(query string, page, limit int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64
	offset := (page - 1) * limit
	like := "%" + query + "%"

	base := r.db.Model(&domain.Product{}).Where("is_active = ? AND (name LIKE ? OR description LIKE ?)", true, like, like)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := base.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindProductByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) CreateProduct(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) UpdateProductDetails(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) DeactivateProduct(id uint) error {
	return r.db.Model(&domain.Product{}).Where("id = ?", id).Update("is_active", false).Error
}
