package model

type Answer struct {
	ID         uint `json:"id" gorm:"PrimaryKey;Autoincrement"`
	QuestionID uint `json:"question_id" gorm:"not null"`
	ChoiceID   uint `json:"choice_id" gorm:"not null"`

	// リレーション（ポインタ型で再帰を避ける）
	Question *Question `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE;"` // 質問とのリレーション
	Choice   *Choice   `gorm:"foreignKey:ChoiceID;references:ID;constraint:OnDelete:CASCADE;"`   // 選択肢とのリレーション
}
