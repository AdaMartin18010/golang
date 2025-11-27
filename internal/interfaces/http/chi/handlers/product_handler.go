package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	app "github.com/yourusername/golang/internal/application/product"
	"github.com/yourusername/golang/pkg/http/response"
)

// ProductHandler 产品处理器
type ProductHandler struct {
	service app.Service
}

// NewProductHandler 创建产品处理器
func NewProductHandler(service app.Service) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct 创建产品
// @Summary 创建产品
// @Tags products
// @Accept json
// @Produce json
// @Param product body app.CreateProductRequest true "产品信息"
// @Success 201 {object} response.Response{data=app.ProductDTO}
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req app.CreateProductRequest
	if err := decodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.service.CreateProduct(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, http.StatusCreated, product)
}

// GetProduct 获取产品
// @Summary 获取产品
// @Tags products
// @Produce json
// @Param id path string true "产品ID"
// @Success 200 {object} response.Response{data=app.ProductDTO}
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
		return
	}

	response.Success(w, http.StatusOK, product)
}

// GetProductBySKU 根据SKU获取产品
// @Summary 根据SKU获取产品
// @Tags products
// @Produce json
// @Param sku path string true "产品SKU"
// @Success 200 {object} response.Response{data=app.ProductDTO}
// @Router /products/sku/{sku} [get]
func (h *ProductHandler) GetProductBySKU(w http.ResponseWriter, r *http.Request) {
	sku := chi.URLParam(r, "sku")
	if sku == "" {
		response.Error(w, http.StatusBadRequest, "product sku is required")
		return
	}

	product, err := h.service.GetProductBySKU(r.Context(), sku)
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
		return
	}

	response.Success(w, http.StatusOK, product)
}

// ListProducts 列出产品
// @Summary 列出产品
// @Tags products
// @Produce json
// @Param limit query int false "限制数量" default(10)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.Response{data=[]app.ProductDTO}
// @Router /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	products, err := h.service.ListProducts(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, http.StatusOK, products)
}

// SearchProducts 搜索产品
// @Summary 搜索产品
// @Tags products
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param limit query int false "限制数量" default(10)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.Response{data=[]app.ProductDTO}
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		response.Error(w, http.StatusBadRequest, "keyword is required")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	products, err := h.service.SearchProducts(r.Context(), keyword, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, http.StatusOK, products)
}

// GetProductsByCategory 根据分类获取产品
// @Summary 根据分类获取产品
// @Tags products
// @Produce json
// @Param category_id path string true "分类ID"
// @Param limit query int false "限制数量" default(10)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.Response{data=[]app.ProductDTO}
// @Router /products/category/{category_id} [get]
func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	if categoryID == "" {
		response.Error(w, http.StatusBadRequest, "category_id is required")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	products, err := h.service.GetProductsByCategory(r.Context(), categoryID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, http.StatusOK, products)
}

// UpdateProduct 更新产品
// @Summary 更新产品
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "产品ID"
// @Param product body app.UpdateProductRequest true "产品信息"
// @Success 200 {object} response.Response{data=app.ProductDTO}
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	var req app.UpdateProductRequest
	if err := decodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.service.UpdateProduct(r.Context(), id, req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, product)
}

// DeleteProduct 删除产品
// @Summary 删除产品
// @Tags products
// @Produce json
// @Param id path string true "产品ID"
// @Success 200 {object} response.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, http.StatusOK, nil)
}

// ActivateProduct 上架产品
// @Summary 上架产品
// @Tags products
// @Produce json
// @Param id path string true "产品ID"
// @Success 200 {object} response.Response{data=app.ProductDTO}
// @Router /products/{id}/activate [post]
func (h *ProductHandler) ActivateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	product, err := h.service.ActivateProduct(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, product)
}

// DeactivateProduct 下架产品
// @Summary 下架产品
// @Tags products
// @Produce json
// @Param id path string true "产品ID"
// @Success 200 {object} response.Response{data=app.ProductDTO}
// @Router /products/{id}/deactivate [post]
func (h *ProductHandler) DeactivateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	product, err := h.service.DeactivateProduct(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, product)
}

// UpdateStock 更新库存
// @Summary 更新库存
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "产品ID"
// @Param stock body object{stock=int} true "库存数量"
// @Success 200 {object} response.Response{data=app.ProductDTO}
// @Router /products/{id}/stock [put]
func (h *ProductHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	var req struct {
		Stock int `json:"stock"`
	}
	if err := decodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.service.UpdateStock(r.Context(), id, req.Stock)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, product)
}

