package handler

import (
	"net/http"
	"skill-typing-back/db"
	"skill-typing-back/model"
	"skill-typing-back/repository"

	"github.com/labstack/echo/v4"
)

// ゲーム用の問題を取得するハンドラー
func GetGameQuestions(c echo.Context) error {
	dbConn := db.GetDB()
	gameRepo := repository.NewGameRepository(dbConn)

	questions, err := gameRepo.GetRandomQuestions(5)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の取得に失敗しました"})
	}

	var response []map[string]interface{}

	for _, question := range questions {
		choices, err := gameRepo.GetChoicesByQuestionID(question.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の取得に失敗しました"})
		}

		// 選択肢が4つ未満の場合はダミーデータを追加
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
			"choices":           choices[:4],
			"correct_choice_id": question.AnswerID,
		}

		response = append(response, questionData)
	}

	return c.JSON(http.StatusOK, response)
}
