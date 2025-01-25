package model

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CognitoID string    `json:"cognitoID" gorm:"size:255;uniqueIndex`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`    // 管理者かどうか
	Name      string    `json:"name" gorm:"size:100;not null"`    // ユーザー名
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // 作成日時
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
