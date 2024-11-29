// back/internal/handlers/feedback_handler.go
package handlers

import (
	"back/internal/models"
	"back/internal/validators"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) SubmitFeedback(c echo.Context) error {
	// Parse and validate incoming JSON data
	var input models.Feedback
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := validators.ValidateFeedback(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Save feedback to the database using the handler's DB instance
	if err := h.DB.Create(&input).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save feedback"})
	}

	// Send confirmation email after saving feedback
	// if err := services.SendEmail(input); err != nil {
	// 	log.Printf("Email sending error: %v", err)
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Feedback saved but failed to send email"})
	// }

	return c.JSON(http.StatusOK, map[string]string{"message": "Feedback submitted and email sent successfully"})
}
