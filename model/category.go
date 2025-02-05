package model

type Category struct {
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"` // 主キー
	Title string `json:"title" gorm:"size:255;not null"`     // カテゴリー名（ユニーク）

	// リレーション
	Questions []Question `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"` //  外部キーの設定
}
