package handlers

import (
	"back/internal/models"
	"back/internal/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// CreateOrder handles incoming orders and saves them to the database
func (h *Handler) CreateOrder(c echo.Context) error {
	var orderRequest struct {
		CustomerInfo models.CustomerInfo `json:"customer_info"`
		TotalAmount  float64             `json:"total_amount"`
		Discount     float64             `json:"discount"`
		// Items        []models.OrderItem  `json:"items"`
		Items []struct {
			Name                        string                                    `json:"name"`
			MoyskladID                  string                                    `json:"moysklad_id"`
			Quantity                    uint                                      `json:"quantity"`
			Price                       float64                                   `json:"price"`
			ModificationCharacteristics []models.ModificationCharacteristicsOrder `json:"modification_characteristics"`
		} `json:"items"`
		SelectedCity string `json:"selected_city"`
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
	var orderItems []models.OrderItem
	for _, item := range orderRequest.Items {
		orderItem := models.OrderItem{
			OrderID:    order.ID,
			Name:       item.Name,
			MoyskladID: item.MoyskladID,
			Quantity:   item.Quantity,
			Price:      item.Price,
		}
		// Save order item and refresh the ID
		if err := h.DB.Create(&orderItem).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save order items"})
		}

		// Now orderItem.ID is populated, use it in modification characteristics
		for _, char := range item.ModificationCharacteristics {
			modChar := models.ModificationCharacteristicsOrder{
				OrderItemID: orderItem.ID, // Ensure this ID is populated
				ModID:       orderItem.MoyskladID,
				Name:        char.Name,
				Value:       char.Value,
			}
			if err := h.DB.Create(&modChar).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save modification characteristics"})
			}
		}
	}
	// for _, item := range orderRequest.Items {
	// 	orderItem := models.OrderItem{
	// 		OrderID:    order.ID,
	// 		Name:       item.Name,
	// 		MoyskladID: item.MoyskladID,
	// 		Quantity:   item.Quantity,
	// 		Price:      item.Price,
	// 	}
	// 	if err := h.DB.Create(&orderItem).Error; err != nil {
	// 		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save order items"})
	// 	}
	// 	orderItems = append(orderItems, orderItem)
	// }

	// Send confirmation email after saving the order
	go func() {
		if err := services.SendOrder(orderRequest.SelectedCity, order, customerInfo, orderItems); err != nil {
			log.Printf("Email sending error: %v", err)
		} else {
			log.Printf("Email sending success\n")
		}
	}()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Order created successfully",
		"order_id":      order.ID,
		"customer_info": customerInfo.ID,
	})
}
