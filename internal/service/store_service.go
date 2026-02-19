package service

import (
	"errors"

	"github.com/yourusername/storefront-saas-go/internal/domain"
	"github.com/yourusername/storefront-saas-go/internal/repository"
)

type CreateStoreRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
}

type StoreService interface {
	ListActiveStores(page, limit int) ([]domain.Store, int64, error)
	GetStoreBySlug(slug string) (*domain.Store, error)
	CreateStore(req CreateStoreRequest) (*domain.Store, error)
}

type storeService struct {
	repo repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository) StoreService {
	return &storeService{repo: repo}
}

func (s *storeService) ListActiveStores(page, limit int) ([]domain.Store, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.ListActiveStores(page, limit)
}

func (s *storeService) GetStoreBySlug(slug string) (*domain.Store, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}
	return s.repo.FindStoreBySlug(slug)
}

func (s *storeService) CreateStore(req CreateStoreRequest) (*domain.Store, error) {
	store := &domain.Store{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IsActive:    true,
	}
	if err := s.repo.CreateStore(store); err != nil {
		return nil, err
	}
	return store, nil
}
