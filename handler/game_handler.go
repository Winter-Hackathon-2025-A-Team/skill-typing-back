package handler

import (
	"errors"
	"log"
	"net/http"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetGameQuestions はランダムなゲーム問題を取得する
func (h *ApiHandler) GetGameQuestions(c echo.Context) error {
	questions, err := h.repo.GetRandomQuestions(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // `errors` を使用
			return c.JSON(http.StatusNotFound, map[string]string{"error": "問題が見つかりませんでした"})
		}
		log.Printf("ERROR: failed to get random questions: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の取得に失敗しました"})
	}

	var response []map[string]interface{}

	for _, question := range questions {
		choices, err := h.repo.GetChoicesByQuestionID(c, question.ID) // c を追加
		if err != nil {
			log.Printf("ERROR: failed to get choices: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の取得に失敗しました"})
		}

		// 選択肢が4つ未満の場合、ダミーデータを追加
		for len(choices) < 4 {
			choices = append(choices, model.Choice{
				ID:          0,
				QuestionID:  question.ID,
				Content:     "N/A",
				Description: "この選択肢は利用できません",
			})
		}

		questionData := map[string]interface{}{
			"id":                question.ID,
			"title":             question.Title,
			"content":           question.Content,
			"choices":           choices[:4],       // 必ず4つに制限
			"correct_choice_id": question.AnswerID, // 正解の選択肢ID
		}

		response = append(response, questionData)
	}

	return c.JSON(http.StatusOK, response)
}
