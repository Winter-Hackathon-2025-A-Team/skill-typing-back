package repository

import (
	"fmt"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

func (r *DbRepository) CreateAnswer(c echo.Context, ansewer *model.Answer) error {

	if err := r.db.Create(ansewer).Error; err != nil {
		return fmt.Errorf("failed to create answer: %v", err)
	}
	return nil
}
