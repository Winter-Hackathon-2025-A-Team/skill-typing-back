package handler

import (
	"errors" // ✅ errors を追加
	"log"
	"net/http"
	"strconv"

	"skill-typing-back/model" // ✅ model をインポート

	"github.com/labstack/echo/v4"
	"gorm.io/gorm" // ✅ gorm を追加
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

	// ✅ 修正: `c` を引数に追加
	questions, err := h.repo.GetQuestionsByCategory(c, categoryID, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // ✅ gorm を使ってエラーハンドリング
			log.Println("指定されたカテゴリの問題が見つかりません")
			return c.JSON(http.StatusNotFound, map[string]string{"error": "指定されたカテゴリの問題が見つかりません"})
		}
		log.Println("データベースから問題を取得できませんでした:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の取得に失敗しました"})
	}

	// JSON レスポンス作成
	response := struct {
		Questions []model.Question `json:"questions"` // ✅ model.Question を使用
	}{Questions: questions}

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
