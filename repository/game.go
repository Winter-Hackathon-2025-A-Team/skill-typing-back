package repository

import (
	"fmt"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

// ゲーム用の問題を取得する
func (r *GameRepository) GetGameQuestions(c echo.Context) ([]model.Question, error) {
	var questions []model.Question
	if err := r.db.Preload("Choices").Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to get game questions: %v", err)
	}
	return questions, nil
}

// 追加のゲーム用問題を取得する
func (r *GameRepository) GetAdditionalGameQuestions(c echo.Context) ([]model.Question, error) {
	var questions []model.Question
	if err := r.db.Preload("Choices").Order("created_at DESC").Limit(10).Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to get additional game questions: %v", err)
	}
	return questions, nil
}
