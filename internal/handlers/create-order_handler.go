package handlers

import (
	"back/internal/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateOrder handles incoming orders and saves them to the database
func (h *Handler) CreateOrder(c echo.Context) error {
	var orderRequest struct {
		CustomerInfo models.CustomerInfo `json:"customer_info"`
		TotalAmount  float64             `json:"total_amount"`
		Discount     float64             `json:"discount"`
		Items        []models.OrderItem  `json:"items"`
	}

	// Parse the incoming JSON request
	if err := c.Bind(&orderRequest); err != nil {
		fmt.Println("Error binding request:", err)
		fmt.Printf("Received body: %+v\n", c.Request().Body)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Step 1: Save the CustomerInfo
	customerInfo := orderRequest.CustomerInfo
	if err := h.DB.Create(&customerInfo).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save customer info"})
	}

	// Step 2: Save the Order with a reference to CustomerInfo
	order := models.Order{
		CustomerInfoID: customerInfo.ID,
		TotalAmount:    orderRequest.TotalAmount,
		OrderDate:      customerInfo.CreatedAt,
	}
	if err := h.DB.Create(&order).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save order"})
	}

	// Step 3: Save OrderItems and link them to the Order
	for _, item := range orderRequest.Items {
		orderItem := models.OrderItem{
			OrderID:    order.ID,
			MoyskladID: item.MoyskladID,
			Quantity:   item.Quantity,
			Price:      item.Price,
		}
		if err := h.DB.Create(&orderItem).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save order items"})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Order created successfully",
		"order_id":      order.ID,
		"customer_info": customerInfo.ID,
	})
}
