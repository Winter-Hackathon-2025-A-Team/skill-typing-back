package repository

import (
	"fmt"

	"skill-typing-back/model"
)

// DBにユーザーが存在するか確認し存在すればuserを返す
func (r *DbRepository) GetUser(id string) (*model.User, error) {

	var user model.User
	result := r.db.Where("ID = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 新規ユーザー登録
func (r *DbRepository) CreateUser(id string, name string, isAdmin bool) (*model.User, error) {

	newUser := model.User{
		ID:      id,
		Name:    name,
		IsAdmin: isAdmin,
	}

	if err := r.db.Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}
