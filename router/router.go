package router

import (
	"skill-typing-back/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルーティング
	api := e.Group("/api")

	// 質問の作成エンドポイント
	api.POST("/questions", handler.CreateQuestion)

	return e
}
