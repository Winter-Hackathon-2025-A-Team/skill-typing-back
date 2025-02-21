package handler

import (
	"errors"
	"fmt"
	"net/http"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 問題の新規作成
func (h *ApiHandler) CreateQuestion(c echo.Context) error {
	// リクエストボディを取得
	q := new(model.Question)
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// 必須フィールドのバリデーション
	if q.Title == "" || q.Content == "" || q.UserID == "" || q.CategoryID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	// データの保存
	if err := h.repo.CreateQuestion(c, q); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create question"})
	}

	return c.JSON(http.StatusCreated, q)
}

// 特定の問題を取得
func (h *ApiHandler) GetQuestion(c echo.Context) error {
	id := c.Param("id")

	question, err := h.repo.GetQuestion(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Question not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to retrieve question: %v", err)})
	}

	return c.JSON(http.StatusOK, question)
}

// すべての問題を取得
func (h *ApiHandler) GetAllQuestions(c echo.Context) error {
	questions, err := h.repo.GetAllQuestions(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to retrieve questions: %v", err)})
	}

	return c.JSON(http.StatusOK, questions)
}
