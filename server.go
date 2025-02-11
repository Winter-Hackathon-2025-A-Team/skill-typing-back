package main

import (
	"fmt"
	"log"
	"os"
	"skill-typing-back/router"

	"github.com/joho/godotenv"
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

	e := router.SetupRouter()

	e.Logger.Fatal(e.Start(":1323"))
}
