// back/internal/validators/feedback_validator.go
package validators

import (
	"back/internal/models"
	"errors"
)

func ValidateFeedback(feedback models.Feedback) error {
	if feedback.Name == "" {
		return errors.New("name is required")
	}
	if feedback.Telephone == "" {
		return errors.New("telephone is required")
	}
	if feedback.City == "" {
		return errors.New("city is required")
	}
	return nil
}
