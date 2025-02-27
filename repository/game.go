package repository

import (
	"errors"
	"fmt"
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

// GetRandomQuestions はランダムに指定された数の質問を取得
func (r *GameRepository) GetRandomQuestions(limit int) ([]model.Question, error) {
	var questions []model.Question
	if err := r.db.Preload("Choices").Preload("Category").Order("RAND()").Limit(limit).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

// GetQuestionsByCategory は指定したカテゴリーの質問を取得
func (r *GameRepository) GetQuestionsByCategory(c echo.Context, categoryID uint, limit int) ([]model.Question, map[uint]uint, error) {
	var questions []model.Question
	answerMap := make(map[uint]uint)

	// 質問データを取得
	if err := r.db.Where("category_id = ?", categoryID).Limit(limit).Find(&questions).Error; err != nil {
		return nil, nil, fmt.Errorf("問題の取得に失敗しました: %w", err)
	}

	// 各 `question` の `answerID` を取得
	for _, question := range questions {
		var answer model.Answer
		if err := r.db.Where("question_id = ?", question.ID).First(&answer).Error; err == nil {
			answerMap[question.ID] = answer.ChoiceID
		}
	}

	return questions, answerMap, nil
}

// GetChoicesByQuestionID は指定された質問の選択肢を取得
func (r *GameRepository) GetChoicesByQuestionID(questionID uint) ([]model.Choice, error) {
	var question model.Question

	// 質問情報を取得し、選択肢IDを取得
	if err := r.db.First(&question, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("質問の取得に失敗しました: %w", err)
	}

	// 選択肢IDリストを作成
	choiceIDs := []uint{question.Choice1ID, question.Choice2ID, question.Choice3ID, question.Choice4ID}

	// 選択肢を取得
	var choices []model.Choice
	if err := r.db.Where("id IN ?", choiceIDs).Find(&choices).Error; err != nil {
		return nil, fmt.Errorf("選択肢の取得に失敗しました: %w", err)
	}

	// 選択肢が4つ未満の場合、ダミーデータを追加
	for len(choices) < 4 {
		choices = append(choices, model.Choice{
			ID:          0,
			Content:     "N/A",
			Description: "この選択肢は利用できません",
		})
	}

	return choices, nil
}
