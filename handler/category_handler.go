package handler

import (
	"net/http"
	"skill-typing-back/db"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

// カテゴリー一覧を取得
func GetCategoriesHandler(c echo.Context) error {
	var categories []model.Category
	dbConn := db.GetDB()
	dbConn.Find(&categories)

	return c.JSON(http.StatusOK, categories)
}
