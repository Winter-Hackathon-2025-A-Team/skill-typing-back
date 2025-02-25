package model

import "time"

type Question struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`              // ID
	UserID     string    `json:"user_id" gorm:"not null;type:varchar(255);index"` // UserID (string に統一)
	CategoryID uint      `json:"category_id" gorm:"not null;index"`               // CategoryID
	// AnswerID   uint      `json:"answer_id" gorm:"not null;index"`                 // AnswerID
	Title      string    `json:"title" gorm:"size:255;not null"`                  // タイトル
	Content    string    `json:"content" gorm:"type:text;not null"`               // 問題文
	Choice1ID  uint      `json:"choice1_id" gorm:"not null"`                      // 選択肢1
	Choice2ID  uint      `json:"choice2_id" gorm:"not null"`                      // 選択肢2
	Choice3ID  uint      `json:"choice3_id" gorm:"not null"`                      // 選択肢3
	Choice4ID  uint      `json:"choice4_id" gorm:"not null"`                      // 選択肢4
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`                // 作成日時
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`                // 更新日時

	User     User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`     // ユーザーとのリレーション
	Category Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"`               // カテゴリとのリレーション
	// Answer   Answer   `gorm:"foreignKey:AnswerID;references:ID;constraint:OnDelete:CASCADE;"` // 正解とのリレーション
	// Choices  []Choice `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE;"` // 選択肢とのリレーション
	Choice1 Choice `gorm:"foreignKey:Choice1ID;references:ID"`
	Choice2 Choice `gorm:"foreignKey:Choice2ID;references:ID"`
	Choice3 Choice `gorm:"foreignKey:Choice3ID;references:ID"`
	Choice4 Choice `gorm:"foreignKey:Choice4ID;references:ID"`
}
