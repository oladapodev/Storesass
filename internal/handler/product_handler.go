package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/storefront-saas-go/internal/service"
	"github.com/yourusername/storefront-saas-go/internal/util"
)

type ProductHandler struct {
	productSvc service.ProductService
	storeSvc   service.StoreService
}

func NewProductHandler(productSvc service.ProductService, storeSvc service.StoreService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc, storeSvc: storeSvc}
}

// ListHotProducts godoc
// @Summary List all active products
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} util.PaginatedResponse
// @Router /products [get]
func (h *ProductHandler) ListHotProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	products, total, err := h.productSvc.ListHotProducts(page, limit)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, util.PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: util.CalcTotalPages(total, limit),
	})
}

// ListProductsByStore godoc
// @Summary List products for a specific store
// @Tags products
// @Produce json
// @Param slug path string true "Store slug"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} util.PaginatedResponse
// @Failure 404 {object} map[string]interface{}
// @Router /stores/{slug}/products [get]
func (h *ProductHandler) ListProductsByStore(c *gin.Context) {
	slug := c.Param("slug")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	store, err := h.storeSvc.GetStoreBySlug(slug)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "store not found")
		return
	}

	products, total, err := h.productSvc.ListProductsByStore(store.ID, page, limit)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, util.PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: util.CalcTotalPages(total, limit),
	})
}

// SearchActiveProducts godoc
// @Summary Search active products
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} util.PaginatedResponse
// @Router /products/search [get]
func (h *ProductHandler) SearchActiveProducts(c *gin.Context) {
	q := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	products, total, err := h.productSvc.SearchActiveProducts(q, page, limit)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, util.PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: util.CalcTotalPages(total, limit),
	})
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.productSvc.GetProductByID(uint(id))
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "product not found")
		return
	}
	util.RespondWithSuccess(c, http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body service.CreateProductRequest true "Product data"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]interface{}
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.productSvc.CreateProduct(req)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.RespondWithSuccess(c, http.StatusCreated, product)
}
