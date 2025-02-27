package model

import "time"

type Choice struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`    // 主キー
	Content     string    `json:"content" gorm:"type:text;not null"`     // 問題の内容
	Description string    `json:"description" gorm:"type:text;not null"` // 問題の解説
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`      // 作成日時
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`      // 更新日時

	// リレーション
	Answers []Answer `gorm:"foreignKey:ChoiceID;references:ID;constraint:OnDelete:CASCADE;"` // ✅ Answerとのリレーション明確化
}
