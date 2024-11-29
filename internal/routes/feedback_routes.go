// routes/feedback_routes.go
package routes

import (
	"back/internal/handlers"
	"github.com/labstack/echo/v4"
)

func InitFeedbackRoutes(e *echo.Echo, h *handlers.Handler) {
	// Define a group for feedback routes
	feedbackGroup := e.Group("/api/feedback")

	// Feedback routes
	feedbackGroup.POST("", h.SubmitFeedback)
}
