package model

type Category struct {
	ID    string `json:"id" gorm:"primaryKey;autoIncrement"` // 主キー
	Title string `json:"title" gorm:"primaryKey"`            // 問題のカテゴリー名

	//リレーション
	Questions []Question `gorm:"foreignKey:CategorieID"` // 質問とのリレーション
}
