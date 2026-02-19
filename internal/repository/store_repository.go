package repository

import (
	"github.com/yourusername/storefront-saas-go/internal/domain"
	"gorm.io/gorm"
)

type StoreRepository interface {
	ListActiveStores(page, limit int) ([]domain.Store, int64, error)
	FindStoreBySlug(slug string) (*domain.Store, error)
	FindStoreByID(id uint) (*domain.Store, error)
	CreateStore(store *domain.Store) error
	UpdateStoreDetails(store *domain.Store) error
	DeactivateStore(id uint) error
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) ListActiveStores(page, limit int) ([]domain.Store, int64, error) {
	var stores []domain.Store
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&domain.Store{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("is_active = ?", true).Offset(offset).Limit(limit).Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

func (r *storeRepository) FindStoreBySlug(slug string) (*domain.Store, error) {
	var store domain.Store
	if err := r.db.Where("slug = ? AND is_active = ?", slug, true).First(&store).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) FindStoreByID(id uint) (*domain.Store, error) {
	var store domain.Store
	if err := r.db.First(&store, id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) CreateStore(store *domain.Store) error {
	return r.db.Create(store).Error
}

func (r *storeRepository) UpdateStoreDetails(store *domain.Store) error {
	return r.db.Save(store).Error
}

func (r *storeRepository) DeactivateStore(id uint) error {
	return r.db.Model(&domain.Store{}).Where("id = ?", id).Update("is_active", false).Error
}
