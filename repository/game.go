package repository

import (
	"errors"
	"fmt"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetRandomQuestions はランダムに5つの質問を取得
func (r *DbRepository) GetRandomQuestions(c echo.Context) ([]model.Question, error) {
	var questions []model.Question
	if err := r.db.Order("RAND()").Limit(5).Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to get random questions: %v", err)
	}
	return questions, nil
}

// GetChoicesByQuestionID は指定された質問の選択肢を取得
func (r *DbRepository) GetChoicesByQuestionID(c echo.Context, questionID uint) ([]model.Choice, error) {
	var choices []model.Choice
	if err := r.db.Where("question_id = ?", questionID).Find(&choices).Error; err != nil {
		// 404 のときだけ明示的にエラーを指定する。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get choices for question: %v", err)
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
