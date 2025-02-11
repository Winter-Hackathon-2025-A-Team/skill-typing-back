package handler

import (
	"net/http"
	"skill-typing-back/auth"
	"skill-typing-back/db"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

// 問題の新規作成
func CreateScore(c echo.Context) error {

	// リクエストボディを取得
	q := new(model.Score)
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Contextからsub(userId)を取得
	user := c.Get("user").(*auth.CognitoClaims)
	sub := user.Sub

	// scoreユーザーIDを設定
	score := *q
	score.UserID = sub

	// DB 接続
	dbConn := db.NewDB()
	if dbConn == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
	}

	//データの保存
	if err := dbConn.Create(&score).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create score"})
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
