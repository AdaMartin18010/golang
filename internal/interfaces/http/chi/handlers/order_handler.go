package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	app "github.com/yourusername/golang/internal/application/order"
	"github.com/yourusername/golang/pkg/errors"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	service app.Service
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(service app.Service) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder 创建订单
// @Summary 创建订单
// @Tags orders
// @Accept json
// @Produce json
// @Param order body app.CreateOrderRequest true "订单信息"
// @Success 201 {object} response.Response{data=app.OrderDTO}
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req app.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid request body"))
		return
	}

	order, err := h.service.CreateOrder(r.Context(), req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusCreated, order)
}

// GetOrder 获取订单
// @Summary 获取订单
// @Tags orders
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response{data=app.OrderDTO}
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Order ID is required"))
		return
	}

	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err)
		return
	}

	Success(w, http.StatusOK, order)
}

// GetUserOrders 获取用户订单列表
// @Summary 获取用户订单列表
// @Tags orders
// @Produce json
// @Param user_id query string true "用户ID"
// @Param limit query int false "限制数量" default(10)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.Response{data=[]app.OrderDTO}
// @Router /orders [get]
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("user_id is required"))
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

	orders, err := h.service.GetUserOrders(r.Context(), userID, limit, offset)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusOK, orders)
}

// PayOrder 支付订单
// @Summary 支付订单
// @Tags orders
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response{data=app.OrderDTO}
// @Router /orders/{id}/pay [post]
func (h *OrderHandler) PayOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order id is required")
		return
	}

	order, err := h.service.PayOrder(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, order)
}

// ShipOrder 发货
// @Summary 发货
// @Tags orders
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response{data=app.OrderDTO}
// @Router /orders/{id}/ship [post]
func (h *OrderHandler) ShipOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order id is required")
		return
	}

	order, err := h.service.ShipOrder(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, order)
}

// DeliverOrder 送达
// @Summary 送达
// @Tags orders
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response{data=app.OrderDTO}
// @Router /orders/{id}/deliver [post]
func (h *OrderHandler) DeliverOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order id is required")
		return
	}

	order, err := h.service.DeliverOrder(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, order)
}

// CancelOrder 取消订单
// @Summary 取消订单
// @Tags orders
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response{data=app.OrderDTO}
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order id is required")
		return
	}

	order, err := h.service.CancelOrder(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, order)
}

// RefundOrder 退款
// @Summary 退款
// @Tags orders
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response{data=app.OrderDTO}
// @Router /orders/{id}/refund [post]
func (h *OrderHandler) RefundOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order id is required")
		return
	}

	order, err := h.service.RefundOrder(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, order)
}

// UpdateOrder 更新订单
// @Summary 更新订单
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param order body app.UpdateOrderRequest true "订单信息"
// @Success 200 {object} response.Response{data=app.OrderDTO}
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order id is required")
		return
	}

	var req app.UpdateOrderRequest
	if err := decodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	order, err := h.service.UpdateOrder(r.Context(), id, req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Success(w, http.StatusOK, order)
}

