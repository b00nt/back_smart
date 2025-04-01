// back/internal/handlers/feedback_handler.go
package handlers

import (
	"back/internal/models"
	"back/internal/services"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func (h *Handler) SubmitFeedback(c echo.Context) error {
	// Parse and validate incoming JSON data
	var input struct {
		models.Feedback        // Embed your Feedback model
		ContextCity     string `json:"contextCity"` // Include context-provided city
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Use input.ContextCity if needed
	log.Printf("Context-provided city: %s", input.ContextCity)

	// Save feedback to the database
	if err := h.DB.Create(&input.Feedback).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save feedback"})
	}

	// Send confirmation email asynchronously
	go func() {
		if err := services.SendFeedback(input.ContextCity, input.Feedback); err != nil {
			log.Printf("Email sending error: %v", err)
		} else {
			log.Printf("Email sending success\n")
		}
	}()

	return c.JSON(http.StatusOK, map[string]string{"message": "Feedback submitted and email sent successfully"})
}
