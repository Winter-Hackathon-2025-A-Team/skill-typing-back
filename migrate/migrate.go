package main

import (
	"fmt"
	"log"

	"skill-typing-back/db"
	"skill-typing-back/model"
)

func main() {
	// DB 接続
	dbConn := db.NewDB()
	if dbConn == nil {
		log.Fatal("Failed to connect to database")
	}

	// DB を確実に閉じる
	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Fatal("Failed to get DB connection:", err)
	}
	defer sqlDB.Close()

	// マイグレーション実行（外部キーの順番に注意）
	fmt.Println("Starting migration...")
	err = dbConn.AutoMigrate(
		&model.User{},     // 1. ユーザー（独立）
		&model.Category{}, // 2. カテゴリー（独立）
		&model.Question{}, // 3. 質問（User, Category に依存）
		&model.Answer{},   // 4. 回答（Question に依存）
		&model.Choice{},   // 5. 選択肢（Question に依存）
		&model.Score{},    // 6. スコア（User に依存）
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Successfully Migrated")
}
