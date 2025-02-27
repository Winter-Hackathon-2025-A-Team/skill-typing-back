package repository

import (
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

// GetChoicesByQuestionID implements DbRepositoryInterface.
func (r *DbRepository) GetChoicesByQuestionID(c echo.Context, questionID uint) ([]model.Choice, error) {
	var question model.Question
	if err := r.db.Where("id = ?", questionID).First(&question).Error; err != nil {
		return nil, err
	}

	// クエリを question の choice1_id ～ choice4_id で取得する
	var choices []model.Choice
	if err := r.db.Where("id IN (?, ?, ?, ?)", question.Choice1ID, question.Choice2ID, question.Choice3ID, question.Choice4ID).Find(&choices).Error; err != nil {
		return nil, err
	}

	return choices, nil
}

// GetQuestionsByCategory implements DbRepositoryInterface.
func (r *DbRepository) GetQuestionsByCategory(c echo.Context, categoryID uint, limit int) ([]model.Question, map[uint]uint, error) {
	var questions []model.Question
	// answerMap を初期化
	answerMap := make(map[uint]uint)

	// カテゴリ情報も含めて質問を取得
	if err := r.db.Preload("Category").Where("category_id = ?", categoryID).Limit(limit).Find(&questions).Error; err != nil {
		return questions, answerMap, err
	}

	// 各質問に対して、answers テーブルから正しい回答を取得
	for _, question := range questions {
		var answer model.Answer
		if err := r.db.Where("question_id = ?", question.ID).First(&answer).Error; err == nil {
			answerMap[question.ID] = answer.ChoiceID
		}
	}

	return questions, answerMap, nil
}
