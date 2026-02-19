package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourusername/storefront-saas-go/internal/domain"
	"github.com/yourusername/storefront-saas-go/internal/repository"
)

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock"`
	StoreID     uint    `json:"store_id" binding:"required"`
}

type ProductService interface {
	ListHotProducts(page, limit int) ([]domain.Product, int64, error)
	ListProductsByStore(storeID uint, page, limit int) ([]domain.Product, int64, error)
	SearchActiveProducts(query string, page, limit int) ([]domain.Product, int64, error)
	GetProductByID(id uint) (*domain.Product, error)
	CreateProduct(req CreateProductRequest) (*domain.Product, error)
}

type productService struct {
	repo  repository.ProductRepository
	redis *redis.Client
}

func NewProductService(repo repository.ProductRepository, redisClient *redis.Client) ProductService {
	return &productService{repo: repo, redis: redisClient}
}

func (s *productService) ListHotProducts(page, limit int) ([]domain.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	if s.redis != nil {
		cacheKey := fmt.Sprintf("hot_products:%d:%d", page, limit)
		ctx := context.Background()

		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var result struct {
				Products []domain.Product `json:"products"`
				Total    int64            `json:"total"`
			}
			if json.Unmarshal([]byte(cached), &result) == nil {
				return result.Products, result.Total, nil
			}
		}

		products, total, err := s.repo.ListActiveProducts(page, limit)
		if err != nil {
			return nil, 0, err
		}

		data, err := json.Marshal(struct {
			Products []domain.Product `json:"products"`
			Total    int64            `json:"total"`
		}{products, total})
		if err != nil {
			log.Printf("warn: failed to marshal cache data: %v", err)
		} else if err := s.redis.Set(ctx, cacheKey, data, 5*time.Minute).Err(); err != nil {
			log.Printf("warn: failed to set redis cache key %s: %v", cacheKey, err)
		}

		return products, total, nil
	}

	return s.repo.ListActiveProducts(page, limit)
}

func (s *productService) ListProductsByStore(storeID uint, page, limit int) ([]domain.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.ListProductsByStore(storeID, page, limit)
}

func (s *productService) SearchActiveProducts(query string, page, limit int) ([]domain.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.SearchActiveProducts(query, page, limit)
}

func (s *productService) GetProductByID(id uint) (*domain.Product, error) {
	return s.repo.FindProductByID(id)
}

func (s *productService) CreateProduct(req CreateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		StoreID:     req.StoreID,
		IsActive:    true,
	}
	if err := s.repo.CreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}
