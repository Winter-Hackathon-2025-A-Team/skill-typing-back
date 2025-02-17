package repository

import (
	"errors"
	"fmt"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// スコアをDBに登録する
func (r *DbRepository) CreateScore(c echo.Context, score *model.Score) error {

	if err := r.db.Create(&score).Error; err != nil {
		return fmt.Errorf("failed to create score : %v", err)
	}
	return nil
}

// 指定ユーザーの最新スコアをDBから取得する
func (r *DbRepository) GetLatestScore(c echo.Context, userId string) (*model.Score, error) {

	var score model.Score

	if err := r.db.Where("user_id = ?", userId).Last(&score).Error; err != nil {
		// 404 のときだけ明示的にエラーを指定する。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get latest score : %v", err)
	}
	return &score, nil
}
