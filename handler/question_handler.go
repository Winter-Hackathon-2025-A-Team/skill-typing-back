package handler

import (
	"net/http"
	"skill-typing-back/db"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

// 問題の新規作成
func CreateQuestion(c echo.Context) error {
	// リクエストボディを取得
	q := new(model.Question)
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// DB 接続
	dbConn := db.NewDB()
	if dbConn == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
	}

	//データの保存
	if err := dbConn.Create(q).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create question"})
	}

	return c.JSON(http.StatusCreated, q)
}
