package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourusername/storefront-saas-go/internal/domain"
	"github.com/yourusername/storefront-saas-go/internal/service"
)

type mockStoreRepo struct {
	mock.Mock
}

func (m *mockStoreRepo) ListActiveStores(page, limit int) ([]domain.Store, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]domain.Store), args.Get(1).(int64), args.Error(2)
}

func (m *mockStoreRepo) FindStoreBySlug(slug string) (*domain.Store, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Store), args.Error(1)
}

func (m *mockStoreRepo) FindStoreByID(id uint) (*domain.Store, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Store), args.Error(1)
}

func (m *mockStoreRepo) CreateStore(store *domain.Store) error {
	args := m.Called(store)
	return args.Error(0)
}

func (m *mockStoreRepo) UpdateStoreDetails(store *domain.Store) error {
	args := m.Called(store)
	return args.Error(0)
}

func (m *mockStoreRepo) DeactivateStore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestListActiveStores(t *testing.T) {
	mockRepo := new(mockStoreRepo)
	svc := service.NewStoreService(mockRepo)

	expected := []domain.Store{
		{ID: 1, Name: "Test Store", Slug: "test-store", IsActive: true},
	}
	mockRepo.On("ListActiveStores", 1, 20).Return(expected, int64(1), nil)

	stores, total, err := svc.ListActiveStores(1, 20)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, stores, 1)
	assert.Equal(t, "Test Store", stores[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestGetStoreBySlug_Found(t *testing.T) {
	mockRepo := new(mockStoreRepo)
	svc := service.NewStoreService(mockRepo)

	expected := &domain.Store{ID: 1, Name: "Test Store", Slug: "test-store", IsActive: true}
	mockRepo.On("FindStoreBySlug", "test-store").Return(expected, nil)

	store, err := svc.GetStoreBySlug("test-store")
	assert.NoError(t, err)
	assert.Equal(t, "test-store", store.Slug)
	mockRepo.AssertExpectations(t)
}

func TestGetStoreBySlug_NotFound(t *testing.T) {
	mockRepo := new(mockStoreRepo)
	svc := service.NewStoreService(mockRepo)

	mockRepo.On("FindStoreBySlug", "missing").Return(nil, errors.New("record not found"))

	store, err := svc.GetStoreBySlug("missing")
	assert.Error(t, err)
	assert.Nil(t, store)
	mockRepo.AssertExpectations(t)
}

func TestGetStoreBySlug_EmptySlug(t *testing.T) {
	mockRepo := new(mockStoreRepo)
	svc := service.NewStoreService(mockRepo)

	store, err := svc.GetStoreBySlug("")
	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestCreateStore(t *testing.T) {
	mockRepo := new(mockStoreRepo)
	svc := service.NewStoreService(mockRepo)

	req := service.CreateStoreRequest{
		Name:        "New Store",
		Slug:        "new-store",
		Description: "A new store",
	}

	mockRepo.On("CreateStore", mock.AnythingOfType("*domain.Store")).Return(nil)

	store, err := svc.CreateStore(req)
	assert.NoError(t, err)
	assert.Equal(t, "New Store", store.Name)
	assert.Equal(t, "new-store", store.Slug)
	mockRepo.AssertExpectations(t)
}
