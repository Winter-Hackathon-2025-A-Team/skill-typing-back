package handler

import (
	"gorm.io/gorm"

	"skill-typing-back/model"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) FindByID(ID string) (*model.User, error) {
	var user model.User
	err := h.db.Where("ID = ?", ID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (h *UserHandler) Create(user *model.User) error {
	return h.db.Create(user).Error
}
