package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/storefront-saas-go/internal/service"
	"github.com/yourusername/storefront-saas-go/internal/util"
)

type StoreHandler struct {
	svc service.StoreService
}

func NewStoreHandler(svc service.StoreService) *StoreHandler {
	return &StoreHandler{svc: svc}
}

// ListActiveStores godoc
// @Summary List active stores
// @Tags stores
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} util.PaginatedResponse
// @Router /stores [get]
func (h *StoreHandler) ListActiveStores(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	stores, total, err := h.svc.ListActiveStores(page, limit)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, util.PaginatedResponse{
		Data:       stores,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: util.CalcTotalPages(total, limit),
	})
}

// GetStoreBySlug godoc
// @Summary Get a store by slug
// @Tags stores
// @Produce json
// @Param slug path string true "Store slug"
// @Success 200 {object} domain.Store
// @Failure 404 {object} map[string]interface{}
// @Router /stores/{slug} [get]
func (h *StoreHandler) GetStoreBySlug(c *gin.Context) {
	slug := c.Param("slug")
	store, err := h.svc.GetStoreBySlug(slug)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "store not found")
		return
	}
	util.RespondWithSuccess(c, http.StatusOK, store)
}

// CreateStore godoc
// @Summary Create a new store
// @Tags stores
// @Accept json
// @Produce json
// @Param store body service.CreateStoreRequest true "Store data"
// @Success 201 {object} domain.Store
// @Failure 400 {object} map[string]interface{}
// @Router /stores [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req service.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	store, err := h.svc.CreateStore(req)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.RespondWithSuccess(c, http.StatusCreated, store)
}
