package handler

import (
	"skill-typing-back/repository"
)

type ApiHandler struct {
	repo *repository.DbRepository
}

// ApiHandler 初期化処理
func New(repo *repository.DbRepository) *ApiHandler {
	h := ApiHandler{
		repo: repo,
	}
	return &h
}
