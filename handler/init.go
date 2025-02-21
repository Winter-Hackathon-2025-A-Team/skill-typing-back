package handler

import (
	"skill-typing-back/repository"
)

type ApiHandler struct {
	repo repository.DbRepositoryInterface
}

// ApiHandler 初期化処理
func New(repo repository.DbRepositoryInterface) *ApiHandler {
	h := ApiHandler{
		repo: repo,
	}
	return &h
}
