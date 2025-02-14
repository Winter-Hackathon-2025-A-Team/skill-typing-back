package handler

import (
	"net/http"
	"skill-typing-back/db"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

// 問題の新規作成（ユーザーが手動で登録）
func CreateQuestion(c echo.Context) error {
	q := new(model.Question)
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// 必須フィールドのバリデーション
	if q.Title == "" || q.Content == "" || q.UserID == "" || q.CategoryID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	// DB 接続
	dbConn := db.GetDB()
	if dbConn == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
	}

	//データの保存
	if err := dbConn.Create(q).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create question"})
	}

	return c.JSON(http.StatusCreated, q)
}

// 特定の問題を取得
func GetQuestion(c echo.Context) error {
	id := c.Param("id")
	var question model.Question

	dbConn := db.GetDB()
	if err := dbConn.First(&question, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Question not found"})
	}

	return c.JSON(http.StatusOK, question)
}

// すべての問題を取得
func GetAllQuestions(c echo.Context) error {
	var questions []model.Question
	dbConn := db.GetDB()
	dbConn.Find(&questions)

	return c.JSON(http.StatusOK, questions)
}
