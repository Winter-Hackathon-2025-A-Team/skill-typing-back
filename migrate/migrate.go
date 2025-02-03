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

	// マイグレーション実行
	fmt.Println("Starting migration...")
	err = dbConn.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Score{},
		&model.Question{},
		&model.Choice{},
		&model.Answer{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Successfully Migrated")
}
