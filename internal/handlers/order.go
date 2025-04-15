package handlers

import (
	"back/internal/models"
	// "back/internal/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// CreateOrder handles incoming orders and saves them to the database
func (h *Handler) CreateOrder(c echo.Context) error {
	// Define a struct to bind the incoming JSON request
	type OrderRequest struct {
		CustomerInfo struct {
			FullName               string `json:"full_name"`
			TelephoneNumber        string `json:"telephone_number"`
			Email                  string `json:"email"`
			Comment                string `json:"comment"`
			City                   string `json:"city"`
			Street                 string `json:"street"`
			House                  string `json:"house"`
			Entrance               string `json:"entrance"`
			Floor                  string `json:"floor"`
			Apartment              string `json:"apartment"`
			AnotherFullName        string `json:"another_full_name"`
			AnotherTelephoneNumber string `json:"another_telephone_number"`
		} `json:"customer_info"`
		TotalAmount float64 `json:"total_amount"`
		Discount    float64 `json:"discount"`
		Items       []struct {
			Name                        string `json:"name"`
			MoyskladID                  string `json:"MoyskladID"`
			Quantity                    uint   `json:"quantity"`
			ModificationCharacteristics []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"modification_characteristics"`
			Price float64 `json:"price"`
		} `json:"items"`
		CurrentCity string `json:"currentCity"`
	}

	// Bind the request body to our struct
	orderReq := new(OrderRequest)
	if err := c.Bind(orderReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request format",
		})
	}

	// Validation
	if orderReq.CustomerInfo.FullName == "" || orderReq.CustomerInfo.TelephoneNumber == "" ||
		orderReq.CustomerInfo.Email == "" || len(orderReq.Items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Missing required fields",
		})
	}

	// Start a transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to start transaction",
		})
	}

	// Create CustomerInfo record
	customerInfo := models.CustomerInfo{
		FullName:               orderReq.CustomerInfo.FullName,
		TelephoneNumber:        orderReq.CustomerInfo.TelephoneNumber,
		Email:                  orderReq.CustomerInfo.Email,
		Comment:                orderReq.CustomerInfo.Comment,
		City:                   orderReq.CustomerInfo.City,
		Street:                 orderReq.CustomerInfo.Street,
		House:                  orderReq.CustomerInfo.House,
		Entrance:               orderReq.CustomerInfo.Entrance,
		Floor:                  orderReq.CustomerInfo.Floor,
		Apartment:              orderReq.CustomerInfo.Apartment,
		AnotherFullName:        orderReq.CustomerInfo.AnotherFullName,
		AnotherTelephoneNumber: orderReq.CustomerInfo.AnotherTelephoneNumber,
	}

	if err := tx.Create(&customerInfo).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to create customer info",
		})
	}

	// Create Order record
	order := models.Order{
		CustomerInfoID: customerInfo.ID,
		TotalAmount:    orderReq.TotalAmount,
		OrderDate:      time.Now(),
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to create order",
		})
	}

	// Create OrderItems with their ModificationCharacteristics
	for _, item := range orderReq.Items {
		orderItem := models.OrderItem{
			OrderID:    order.ID,
			Name:       item.Name,
			MoyskladID: item.MoyskladID,
			Quantity:   item.Quantity,
			Price:      item.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to create order item",
			})
		}

		// Create ModificationCharacteristics for this item
		for _, char := range item.ModificationCharacteristics {
			modChar := models.ModificationCharacteristicOrder{
				OrderItemID: orderItem.ID,
				ModID:       "", // Note: You might need to get this from somewhere
				Name:        char.Name,
				Value:       char.Value,
			}

			if err := tx.Create(&modChar).Error; err != nil {
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"message": "Failed to create modification characteristic",
				})
			}
		}

		// Optional: Update product stock if needed
		// This would reduce the stock of each product by the quantity ordered
		// Uncomment and implement if needed
		/*
		   var product models.Product
		   if err := tx.Where("moysklad_id = ?", item.MoyskladID).First(&product).Error; err == nil {
		       if err := tx.Model(&product).Update("stock", product.Stock - int(item.Quantity)).Error; err != nil {
		           tx.Rollback()
		           return c.JSON(http.StatusInternalServerError, map[string]string{
		               "message": "Failed to update product stock",
		           })
		       }
		   }
		*/
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to commit transaction",
		})
	}

	// Return success
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Order created successfully",
		"order_id": order.ID,
	})
}
