package repository

import (
	"fmt"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

func (r DbRepository) CreateScore(c echo.Context, score *model.Score) error {

	if err := r.db.Create(&score).Error; err != nil {
		return fmt.Errorf("failed to create score : %v", err)
	}
	return nil
}
