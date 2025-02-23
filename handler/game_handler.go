package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// getGameQuestions は指定したカテゴリーのゲーム問題を取得する共通関数
func (h *ApiHandler) getGameQuestions(c echo.Context, limit int) error {
	log.Println("getGameQuestions 関数が呼び出されました")

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

	//  修正：categoryID と limit のみ渡す
	questions, err := h.repo.GetQuestionsByCategory(categoryID, limit)
	if err != nil {
		log.Println("データベースから問題を取得できませんでした:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の取得に失敗しました"})
	}

	// JSON レスポンスを作成
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

	// データをレスポンス形式に変換
	for _, question := range questions {
		choices, err := h.repo.GetChoicesByQuestionID(c, question.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("選択肢が見つかりません:", question.ID)
				return c.JSON(http.StatusNotFound, map[string]string{"error": "選択肢が見つかりません"})
			}
			log.Println("選択肢の取得に失敗しました:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の取得に失敗しました"})
		}

		choiceList := []Choice{}
		for _, choice := range choices {
			choiceList = append(choiceList, Choice{
				ID:          choice.ID,
				Content:     choice.Content,
				Description: choice.Description,
			})
		}

		for len(choiceList) < 4 {
			choiceList = append(choiceList, Choice{
				ID:          0,
				Content:     "N/A",
				Description: "この選択肢は利用できません",
			})
		}

		q := Question{
			ID:       question.ID,
			Title:    question.Title,
			Content:  question.Content,
			AnswerID: question.AnswerID,
			Choices:  choiceList,
		}
		q.Category.ID = question.Category.ID
		q.Category.Title = question.Category.Title

		response.Questions = append(response.Questions, q)
	}

	log.Println("getGameQuestions のレスポンスを正常に返却します")
	return c.JSON(http.StatusOK, response)
}

// GetGameQuestions は最初の8問を取得するAPI
func (h *ApiHandler) GetGameQuestions(c echo.Context) error {
	return h.getGameQuestions(c, 8)
}

// GetAdditionalGameQuestionsは追加のゲーム問題を取得するAPI
func (h *ApiHandler) GetAdditionalGameQuestions(c echo.Context) error {
	return h.getGameQuestions(c, 5)
}
