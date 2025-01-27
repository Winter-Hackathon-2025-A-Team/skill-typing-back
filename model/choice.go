package model

import "time"

tyoe Choice struct {
	ID          int       `json:"id" gorm:"primaryKey;Autoincrement"`          // 主キー
	Content     string    `json:"content" gorm:"type:text;not null"`           // 問題の内容
	Description string    `json:"description" gorm:"type:text;not null"`       // 問題の解説
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`            // 作成日時
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`            // 更新日時

	// リレーション
    Question Question `gorm:"foreignKey:QuestionID"`
}
