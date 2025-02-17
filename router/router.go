package router

import (
	"net/http"
	"skill-typing-back/auth"
	"skill-typing-back/db"
	"skill-typing-back/handler"
	"skill-typing-back/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter() *echo.Echo {

	e := echo.New()

	// DBコネクション作成
	db := db.NewDB()
	// DBリポジトリ初期化
	repo := repository.New(db)
	// APIハンドラ構造体初期化
	apiHandler := handler.New(repo)

	// Cognitoの認証設定
	cognitoAuth, err := auth.NewCognitoAuth(repo)
	if err != nil {
		e.Logger.Fatal(err)
	}
	// Cognitoのユーザーサービスの初期化
	cognitoService, err := auth.NewCognitoUserService()
	if err != nil {
		e.Logger.Fatal(err)
	}
	// authMiddlewareの初期化
	authMiddleware := cognitoAuth.AuthMiddleware(cognitoService)

	/*
	* Middleware設定
	 */
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	/*
	* ルーティング(cognito認証なし)
	 */
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	/*
	* ルーティング(cognito認証あり)
	 */

	api := e.Group("/api")

	// api配下にのみauthMiddlwareを適用
	api.Use(authMiddleware)

	// authエンドポイント
	api.GET("/auth", func(c echo.Context) error {
		user := c.Get("user").(*auth.CognitoClaims)
		return c.JSON(http.StatusOK, map[string]string{
			"sub": user.Sub,
		})
	})
	// ユーザー情報取得エンドポイント
	api.GET("/users/me", apiHandler.GetMe)
	// 質問の作成エンドポイント
	api.POST("/questions", apiHandler.CreateQuestion)
	// スコアの作成エンドポイント
	api.POST("/scores", apiHandler.CreateScore)
	// 最新スコア取得エンドポイント
	api.GET("/scores/latest", apiHandler.GetLatestScore)
	//  AI 生成クイズ API
	api.GET("/generate-quiz", apiHandler.GenerateQuizHandler)
	//ゲーム画面の取得エンドポイント
	api.GET("/game/questions", handler.GetGameQuestions)

	return e
}
