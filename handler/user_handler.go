package handler

import (
	"net/http"
	"skill-typing-back/auth"

	"github.com/labstack/echo/v4"
)

// ログインユーザーの情報を取得する
func (h *ApiHandler) GetMe(c echo.Context) error {

	user := c.Get("user").(*auth.CognitoClaims)

	dbUser, err := h.repo.GetUser(user.Sub)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user info",
		})
	}

	return c.JSON(http.StatusOK, dbUser)
}
