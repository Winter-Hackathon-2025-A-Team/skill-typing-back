package model

type Answer struct {
	ID         int `json:"id" gorm:"PrimaryKey;Autoincrement"`
	QuestionID int `json:"question_id" gorm:"not null"`
	ChoiceID   int `json:"choice_id" gorm:"not null"`

	// リレーション（ポインタ型で再帰を避ける）
	Question *Question `gorm:"foreignKey:QuestionID"` // 質問とのリレーション：ポインタ型に変更
	Choice   *Choice   `gorm:"foreignKey:ChoiceID"`   // 選択肢とのリレーション：ポインタ型に変更
}
