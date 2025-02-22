package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetGameQuestions は指定したカテゴリーのゲーム問題を取得する
func (h *ApiHandler) GetGameQuestions(c echo.Context) error {
	log.Println("GetGameQuestions 関数が呼び出されました")

	// クエリパラメータから category_id を取得
	categoryIDStr := c.QueryParam("category_id")
	var categoryID uint

	if categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			log.Println("category_id のパースに失敗:", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効な category_id です"})
		}
		categoryID = uint(id)
	}

	// 質問の取得数を指定（最初は10件取得）
	limit := 10

	// 指定したカテゴリーの質問を取得
	questions, err := h.repo.GetQuestionsByCategory(c, categoryID, limit)
	if err != nil {
		log.Println("データベースから問題を取得できませんでした:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の取得に失敗しました"})
	}

	// JSON レスポンス用の構造体
	type Choice struct {
		ID          uint   `json:"id"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}

	type Question struct {
		ID       uint   `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category struct {
			ID    uint   `json:"id"`
			Title string `json:"title"`
		} `json:"category"`
		AnswerID uint     `json:"answer_id"`
		Choices  []Choice `json:"choices"`
	}

	response := struct {
		Questions []Question `json:"questions"`
	}{}

	// 取得した質問データを JSON 構造に整形
	for _, question := range questions {
		// 質問に紐づく選択肢を取得
		choices, err := h.repo.GetChoicesByQuestionID(c, question.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("選択肢が見つかりません:", question.ID)
				return c.JSON(http.StatusNotFound, map[string]string{"error": "選択肢が見つかりません"})
			}
			log.Println("選択肢の取得に失敗しました:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の取得に失敗しました"})
		}

		// 選択肢リストを作成
		choiceList := []Choice{}
		for _, choice := range choices {
			choiceList = append(choiceList, Choice{
				ID:          choice.ID,
				Content:     choice.Content,
				Description: choice.Description,
			})
		}

		// 選択肢が4つ未満の場合、ダミーデータを追加
		for len(choiceList) < 4 {
			choiceList = append(choiceList, Choice{
				ID:          0,
				Content:     "N/A",
				Description: "この選択肢は利用できません",
			})
		}

		// 質問データを作成
		q := Question{
			ID:       question.ID,
			Title:    question.Title,
			Content:  question.Content,
			AnswerID: question.AnswerID,
			Choices:  choiceList,
		}
		q.Category.ID = question.Category.ID
		q.Category.Title = question.Category.Title

		// レスポンスに追加
		response.Questions = append(response.Questions, q)
	}

	log.Println("GetGameQuestions のレスポンスを正常に返却します")
	return c.JSON(http.StatusOK, response)
}
