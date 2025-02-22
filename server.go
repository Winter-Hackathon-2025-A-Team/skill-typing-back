package main

import (
	"fmt"
	"os"
	"skill-typing-back/router"
)

func init() {
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
