package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetGameQuestions はランダムなゲーム問題を取得する
func (h *ApiHandler) GetGameQuestions(c echo.Context) error {
	fmt.Println("GetGameQuestions 関数が呼び出されました")

	questions, err := h.repo.GetRandomQuestions(c)
	if err != nil {
		log.Println("データベースから問題を取得できませんでした:", err)
		fmt.Printf("エラー詳細: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の取得に失敗しました"})
	}

	// **デバッグ用ログ**
	if len(questions) == 0 {
		log.Println("データベースに問題がありません")
		fmt.Println("取得された質問の数が 0 件です")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "データベースに問題がありません"})
	}

	// JSON 構造体を作成
	response := struct {
		Questions []map[string]interface{} `json:"questions"`
	}{Questions: []map[string]interface{}{}}

	for _, question := range questions {
		// 選択肢を取得
		choices, err := h.repo.GetChoicesByQuestionID(c, question.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("選択肢が見つかりません:", question.ID)
				fmt.Printf("質問 ID: %d に関連する選択肢が見つかりません\n", question.ID)
				return c.JSON(http.StatusNotFound, map[string]string{"error": "選択肢が見つかりません"})
			}
			log.Println("選択肢の取得に失敗しました:", err)
			fmt.Printf("質問 ID: %d の選択肢取得エラー: %v\n", question.ID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の取得に失敗しました"})
		}

		// JSON 形式で整形
		choiceList := []map[string]interface{}{}
		for _, choice := range choices {
			choiceData := map[string]interface{}{
				"id":          choice.ID,
				"content":     choice.Content,
				"description": choice.Description,
			}
			choiceList = append(choiceList, choiceData)
		}

		// 選択肢が4つ未満のときダミーデータを追加
		for len(choiceList) < 4 {
			choiceList = append(choiceList, map[string]interface{}{
				"id":          0,
				"content":     "N/A",
				"description": "この選択肢は利用できません",
			})
		}

		// **Goコンパイラの未使用警告を防ぐために model を利用**
		category := model.Category{
			ID:    question.Category.ID,
			Title: question.Category.Title,
		}

		questionData := map[string]interface{}{
			"id":      question.ID,
			"title":   question.Title,
			"content": question.Content,
			"category": map[string]interface{}{
				"id":    category.ID, // model.Category のフィールドを使用
				"title": category.Title,
			},
			"answer_id": question.AnswerID,
			"choices":   choiceList,
		}

		response.Questions = append(response.Questions, questionData)
	}

	fmt.Println("GetGameQuestions のレスポンスを正常に返却します")
	return c.JSON(http.StatusOK, response)
}
