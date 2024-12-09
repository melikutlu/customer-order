package internal

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tesodev-korpes/OrderService/internal/types"
)

type Handler struct {
	service *Service
}

// @title Order Service API
// @version 1.0
// @description This is the Order Service API for handling CRUD operations related to orders.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8002
// @BasePath /api/v1
// @schemes http
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service}

	g := e.Group("/order")
	g.GET("/:id", handler.GetByID)
	g.POST("/:customer_id", handler.CreateOrder)
	g.PUT("/:id", handler.Update)
	g.PATCH("/:id", handler.PartialUpdate)
	g.DELETE("/:id", handler.Delete)
}

// GetByID retrieves an order by its ID.
// @Summary Get order by ID
// @Description Get order details by ID
// @Tags order
// @Produce  json
// @Param id path string true "Order ID"
// @Param Authentication header string true "JWT token"
// @Success 200 {object} types.OrderResponseModel
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /order/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {

	id := c.Param("id")
	order, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, order)
}

// CreateOrder creates a new order for a customer.
// @Summary Create a new order
// @Description Create a new order for a specific customer
// @Tags order
// @Accept  json
// @Produce  json
// @Param customer_id path string true "Customer ID"
// @Param order body types.OrderRequestModel true "Order data"
// @Param Authentication header string true "JWT token"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /order/{customer_id} [post]
func (h *Handler) CreateOrder(c echo.Context) error {
	var order types.OrderRequestModel

	customerid := c.Param("customer_id")
	token := c.Request().Header.Get("Authentication")

	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	id, err := h.service.CreateOrderService(c.Request().Context(), customerid, &order, token)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"message":   "Success!",
		"createdId": id,
	}

	return c.JSON(http.StatusCreated, response)
}

// Update modifies an existing order's details.
// @Summary Update order details
// @Description Update order details with the given data
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Param Authentication header string true "JWT token"
// @Param order body types.OrderUpdateModel true "Order data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /order/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var order types.OrderUpdateModel
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.service.Update(c.Request().Context(), id, order); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order updated successfully",
	})
}

// PartialUpdate modifies specific fields of an existing order.
// @Summary Partially update order details
// @Description Partially update order details with the given data
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Param Authentication header string true "JWT token"
// @Param order body types.OrderUpdateModel true "Order data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /order/{id} [patch]
func (h *Handler) PartialUpdate(c echo.Context) error {
	id := c.Param("id")
	var order types.OrderUpdateModel
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.service.Update(c.Request().Context(), id, order); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order partially updated successfully",
	})
}

// Delete removes an order from the database.
// @Summary Delete order
// @Description Delete an order by its ID
// @Tags order
// @Produce  json
// @Param Authentication header string true "JWT token"
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} string
// @Router /order/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order deleted successfully",
	})
}
