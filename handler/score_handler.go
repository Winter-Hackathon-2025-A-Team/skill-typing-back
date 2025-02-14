package handler

import (
	"log"
	"net/http"
	"skill-typing-back/auth"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

// 問題の新規作成
func (h *ApiHandler) CreateScore(c echo.Context) error {

	// リクエストボディを取得
	s := new(model.Score)
	if err := c.Bind(s); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Contextからsub(userId)を取得
	user := c.Get("user").(*auth.CognitoClaims)
	sub := user.Sub

	// scoreユーザーIDを設定
	score := *s
	score.UserID = sub

	//データの保存
	if err := h.repo.CreateScore(c, &score); err != nil {
		log.Printf("ERROR: failed to imprement CreateScore: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	// 成功した場合のメッセージ
	type Message struct {
		Message string `json:"message"`
	}
	success := &Message{
		Message: "registered successfully",
	}

	return c.JSON(http.StatusCreated, success)
}
