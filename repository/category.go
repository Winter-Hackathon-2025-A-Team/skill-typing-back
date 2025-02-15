package repository

import (
	"errors"
	"fmt"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (r *DbRepository) GetCategoryByTitle(c echo.Context, title string) (model.Category, error) {

	var category model.Category

	if err := r.db.Where("title = ?", title).First(&category).Error; err != nil {
		// 404 のときだけ明示的にエラーを指定する。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return category, gorm.ErrRecordNotFound
		}
		return category, fmt.Errorf("failed to get category by title : %v", err)
	}

	return category, nil
}

func (r *DbRepository) CreateCategory(c echo.Context, category *model.Category) error {

	if err := r.db.Create(category).Error; err != nil {
		return fmt.Errorf("failed to create category: %v", err)
	}

	return nil
}
