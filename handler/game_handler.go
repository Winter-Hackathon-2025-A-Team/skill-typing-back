package handler

import (
	"net/http"
	"skill-typing-back/db"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

// ゲーム用の問題をランダムに取得する
func GetGameQuestions(c echo.Context) error {
	var questions []model.Question
	dbConn := db.GetDB()

	// ランダムに5問取得
	if err := dbConn.Order("RAND()").Limit(5).Find(&questions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の取得に失敗しました"})
	}

	var response []map[string]interface{}

	// 各問題に選択肢を取得して追加
	for _, question := range questions {
		var choices []model.Choice
		if err := dbConn.Where("question_id = ?", question.ID).Find(&choices).Error; err != nil {
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

		// JSON レスポンス用の構造
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
