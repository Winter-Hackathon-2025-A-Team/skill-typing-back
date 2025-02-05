package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"skill-typing-back/auth"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".envファイルが読み込めません: %v", err)
	}
	// 環境変数の値を確認
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	clientID := os.Getenv("COGNITO_CLIENT_ID")
	
	fmt.Printf("COGNITO_USER_POOL_ID: %s\n", userPoolID)
	fmt.Printf("COGNITO_CLIENT_ID: %s\n", clientID)
}

func main() {
	e := echo.New()

	// CORSミドルウェアの設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Cognitoの認証設定
	cognitoAuth, err := auth.NewCognitoAuth()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// 認証が必要なルートグループ
	api := e.Group("/api")
	api.Use(cognitoAuth.AuthMiddleware())

	
	api.GET("/auth", func(c echo.Context) error {
		user := c.Get("user").(*auth.CognitoClaims)
		return c.JSON(http.StatusOK, map[string]string{
			"sub": user.Sub,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
