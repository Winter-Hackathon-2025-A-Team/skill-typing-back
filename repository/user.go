package repository

import (
	"fmt"

	"skill-typing-back/db"
	"skill-typing-back/model"
)

// DBにユーザーが存在するか確認し存在すればuserを返す
func GetUser(id string) (*model.User, error) {

	dbConn := db.NewDB()
	if dbConn == nil {
		return nil, fmt.Errorf("Failed to connect to database")
	}

	var user model.User
	result := dbConn.Where("ID = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 新規ユーザー登録
func CreateUser(id string, name string, isAdmin bool) (*model.User, error) {

	dbConn := db.NewDB()
	if dbConn == nil {
		return nil, fmt.Errorf("failed to connect to database")
	}

	newUser := model.User{
		ID:      id,
		Name:    name,
		IsAdmin: isAdmin,
	}

	if err := dbConn.Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}
