package handler

import (
	"net/http"
	"skill-typing-back/db"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

// 過去の問題を取得
func GetQuizHistoryHandler(c echo.Context) error {
	var questions []model.Question
	db.DB.Preload("Choice").Find(&questions)
	return c.JSON(http.StatusOK, questions)
}
