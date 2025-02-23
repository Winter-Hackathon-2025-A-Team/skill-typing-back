package repository

import (
	"errors"
	"fmt"
	"skill-typing-back/model"

	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

// NewGameRepository は新しい GameRepository インスタンスを作成
func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

// GetRandomQuestions はランダムに指定された数の質問を取得
func (r *GameRepository) GetRandomQuestions(limit int) ([]model.Question, error) {
	var questions []model.Question
	if err := r.db.Preload("Choices").Preload("Category").Order("RAND()").Limit(limit).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

// GetQuestionsByCategory は指定したカテゴリーの質問を取得
func (r *GameRepository) GetQuestionsByCategory(categoryID uint, limit int) ([]model.Question, error) {
	var questions []model.Question
	if err := r.db.Where("category_id = ?", categoryID).Limit(limit).Find(&questions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("指定されたカテゴリの問題が見つかりません")
		}
		return nil, fmt.Errorf("問題の取得に失敗しました: %w", err)
	}
	return questions, nil
}

// GetChoicesByQuestionID は指定された質問の選択肢を取得
func (r *GameRepository) GetChoicesByQuestionID(questionID uint) ([]model.Choice, error) {
	var choices []model.Choice
	if err := r.db.Where("question_id = ?", questionID).Find(&choices).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("選択肢の取得に失敗しました: %w", err)
	}

	// 選択肢が4つ未満の場合、ダミーデータを追加
	for len(choices) < 4 {
		choices = append(choices, model.Choice{
			ID:          0,
			QuestionID:  questionID,
			Content:     "N/A",
			Description: "この選択肢は利用できません",
		})
	}

	return choices, nil
}
