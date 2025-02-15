package handler

import (
	"errors"
	"log"
	"net/http"
	"skill-typing-back/auth"
	"skill-typing-back/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

func (h *ApiHandler) GetLatestScore(c echo.Context) error {

	// 返却値用の構造体
	type Score struct {
		Score string `json:"score"`
	}

	// Contextからsub(userId)を取得
	user := c.Get("user").(*auth.CognitoClaims)
	userId := user.Sub

	//　ユーザーの最新スコアの取得
	score, err := h.repo.GetLatestScore(c, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusOK, Score{Score: "norecord"})
		}
		log.Printf("ERROR: failed to imprement GetLatestScore: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, Score{Score: strconv.FormatUint(uint64(score.Score), 10)})
}
