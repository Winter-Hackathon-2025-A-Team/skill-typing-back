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

	// `answerMap` も取得するように修正
	questions, answerMap, err := h.repo.GetQuestionsByCategory(c, categoryID, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("指定されたカテゴリの問題が見つかりません")
			return c.JSON(http.StatusNotFound, map[string]string{"error": "指定されたカテゴリの問題が見つかりません"})
		}
		log.Println("データベースから問題を取得できませんでした:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の取得に失敗しました"})
	}

	// JSON レスポンス作成
	type ChoiceResponse struct {
		ID          uint   `json:"id"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}

	type QuestionResponse struct {
		ID       uint   `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		AnswerID uint   `json:"answer_id"` // ✅ `answers` テーブルから取得
		Category struct {
			ID    uint   `json:"id"`
			Title string `json:"title"`
		} `json:"category"`
		Choices []ChoiceResponse `json:"choices"`
	}

	var response struct {
		Questions []QuestionResponse `json:"questions"`
	}

	for _, question := range questions {
		// 選択肢を取得
		choices, err := h.repo.GetChoicesByQuestionID(c, question.ID)
		if err != nil {
			log.Println("選択肢の取得に失敗しました:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の取得に失敗しました"})
		}

		var choiceResponses []ChoiceResponse
		for _, choice := range choices {
			choiceResponses = append(choiceResponses, ChoiceResponse{
				ID:          choice.ID,
				Content:     choice.Content,
				Description: choice.Description,
			})
		}

		// ✅ `answers` テーブルから `answer_id` を取得
		answerID := answerMap[question.ID]

		//  `category` の情報を正しく格納
		q := QuestionResponse{
			ID:       question.ID,
			Title:    question.Title,
			Content:  question.Content,
			AnswerID: answerID, //  `answers` テーブルから取得
			Choices:  choiceResponses,
		}

		// ✅ カテゴリが 0 でない場合に情報をセット
		if question.Category.ID != 0 {
			q.Category.ID = question.Category.ID
			q.Category.Title = question.Category.Title
		} else {
			q.Category.ID = 0
			q.Category.Title = "未分類" // 例えば「未分類」などのデフォルト値を設定
		}

		response.Questions = append(response.Questions, q)
	}

	log.Println("getGameQuestions のレスポンスを正常に返却します")
	return c.JSON(http.StatusOK, response)
}

// GetGameQuestions は最初の8問を取得するAPI
func (h *ApiHandler) GetGameQuestions(c echo.Context) error {
	return h.getGameQuestions(c, 8)
}

// GetAdditionalGameQuestions は追加のゲーム問題を取得するAPI
func (h *ApiHandler) GetAdditionalGameQuestions(c echo.Context) error {
	return h.getGameQuestions(c, 5)
}
