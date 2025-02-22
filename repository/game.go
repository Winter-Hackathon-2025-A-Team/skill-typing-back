package repository

import (
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

// NewGameRepository は新しい GameRepository インスタンスを作成
func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

// GetQuestionsByCategory は指定されたカテゴリーの質問を取得
func (r *GameRepository) GetQuestionsByCategory(c echo.Context, categoryID uint) ([]model.Question, error) {
	var questions []model.Question
	query := r.db

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Order("RAND()").Limit(5).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}
