package router

import (
	"net/http"
	"skill-typing-back/auth"
	"skill-typing-back/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {
	e := echo.New()

	// CORSミドルウェアの設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Cognitoの認証設定
	cognitoAuth, err := auth.NewCognitoAuth()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Cognitoのユーザーサービスの初期化
	cognitoService, err := auth.NewCognitoUserService()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルーティング
	api := e.Group("/api")

	api.Use(cognitoAuth.AuthMiddleware(cognitoService))

	api.GET("/auth", func(c echo.Context) error {
		user := c.Get("user").(*auth.CognitoClaims)
		return c.JSON(http.StatusOK, map[string]string{
			"sub": user.Sub,
		})
	})
	// ユーザー情報取得エンドポイント
	api.GET("/user/me", handler.GetMe)

	// 質問の作成エンドポイント
	api.POST("/questions", handler.CreateQuestion)
	//  AI 生成クイズ API
	e.GET("/generate-quiz", handler.GenerateQuizHandler)

	// 質問の作成エンドポイント
	api.POST("/scores", handler.CreateScore)

	return e
}
