package repository

import (
	"errors"
	"fmt"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (r *DbRepository) CreateChoice(c echo.Context, choice *model.Choice) error {

	if err := r.db.Create(choice).Error; err != nil {
		return fmt.Errorf("failed to create choice: %v", err)
	}
	return nil
}

func (r *DbRepository) GetChoiceByContentAndQuestionId(c echo.Context, content string, questionId uint) (*model.Choice, error) {

	var choice model.Choice

	if err := r.db.Where("content = ?", content).Where("question_id = ?", questionId).First(&choice).Error; err != nil {
		// 404 のときだけ明示的にエラーを指定する。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &choice, gorm.ErrRecordNotFound
		}
		return &choice, fmt.Errorf("failed to get choice by content: %v", err)
	}

	return &choice, nil
}

// 選択肢の内容を確認するためのリポジトリ
func (r *DbRepository) GetChoiceByContent(c echo.Context, content string) (model.Choice, error) {
	var choice model.Choice

	if err := r.db.Where("content = ?", content).First(&choice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return choice, gorm.ErrRecordNotFound
		}
		return choice, fmt.Errorf("failed to get choice by content: %v", err)
	}

	return choice, nil 
}