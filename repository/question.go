package repository

import (
	"fmt"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

func (r DbRepository) CreateQuestion(c echo.Context, question *model.Question) error {

	if err := r.db.Create(question).Error; err != nil {
		return fmt.Errorf("failed to create question : %v", err)
	}
	return nil
}
