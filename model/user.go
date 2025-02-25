package model

import "time"

type User struct {
	ID string `json:"id" gorm:"primaryKey;type:varchar(255)"` // トークンから取得したSubの値
	// CognitoID string    `json:"cognitoID" gorm:"size:255;unique"`   // AWS CognitoID（修正）
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`    // 管理者かどうか
	Name      string    `json:"name" gorm:"size:100;not null"`    // ユーザー名
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // 作成日時
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // 更新日時

	// リレーション
	Questions []Question `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // ユーザーが作成した質問
	Scores    []Score    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // ユーザーのスコア
}
