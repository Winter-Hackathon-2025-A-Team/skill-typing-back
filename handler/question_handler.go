package handler

import (
	"net/http"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

// 問題の新規作成
func (h *ApiHandler) CreateQuestion(c echo.Context) error {
	// リクエストボディを取得
	q := new(model.Question)
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	//データの保存
	if err := h.repo.CreateQuestion(c, q); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create question"})
	}

	return c.JSON(http.StatusCreated, q)
}
