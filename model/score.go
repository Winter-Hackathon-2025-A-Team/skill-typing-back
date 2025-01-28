package model

import "time"

type Score struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"` // 主キー
	UserID    int       `json:"user_id" gorm:"not null"`            // 関連するユーザー
	Score     int       `json:"score" gorm:"not null"`              // スコア
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`   // 問題実行日時

	// リレーション
	User User `gorm:"foreignKey:UserID"` // ユーザーとのリレーション
}
