package model

import "time"

type Score struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"` // 主キー
	UserID    uint      `json:"user_id" gorm:"not null;index"`      // 関連するユーザー
	Score     uint      `json:"score" gorm:"not null"`              // スコア
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`   // 問題実行日時

	// リレーション
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"` //  外部キー制約の明示
}
