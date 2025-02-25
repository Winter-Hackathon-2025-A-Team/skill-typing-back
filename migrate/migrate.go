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
		log.Fatal("❌ データベースの初期化に失敗しました")
	}
	log.Println("✅ データベースの初期化に成功")

	// DB を確実に閉じる
	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Fatal("Failed to get DB connection:", err)
	}
	defer sqlDB.Close()

	// マイグレーション実行（外部キーの順番に注意）
	fmt.Println("Starting migration...")
	err = dbConn.AutoMigrate(
		&model.User{},     
		&model.Category{},
		&model.Choice{}, 
		&model.Question{}, 
		&model.Answer{},   
		&model.Score{},    
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Successfully Migrated")
}
