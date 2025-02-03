package model

import "time"

type Question struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"` // ID
	UserID     int       `json:"user_id" gorm:"not null;index"`      // userID
	CategoryID uint      `json:"category_id" gorm:"not null;index"`  // CategoryID
	Title      string    `json:"title" gorm:"size:255;not null"`     // タイトル
	AnswerID   string    `json:"answer_id" gorm:"not null"`          // answerID
	Content    string    `json:"content" gorm:"type:text;not null"`  // 問題文
	Choice1ID  string    `json:"choice1_id" gorm:"not null"`         // 選択肢1
	Choice2ID  string    `json:"choice2_id" gorm:"not null"`         // 選択肢2
	Choice3ID  string    `json:"choice3_id" gorm:"not null"`         // 選択肢3
	Choice4ID  string    `json:"choice4_id" gorm:"not null"`         // 選択肢4
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`   // 作成日時
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`   // 更新日時

	User     User     `gorm:"foreignKey:UserID"`     // ユーザーとのリレーション
	Category Category `gorm:"foreignKey:CategoryID"` // カテゴリとのリレーション
	Answer   Answer   `gorm:"foreignKey:AnswerID"`   // 正解とのリレーション
	Choices  []Choice `gorm:"foreignKey:QuestionID"` // 選択肢とのリレーション
}
