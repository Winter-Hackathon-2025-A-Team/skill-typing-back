package handler

import (
	"net/http"
	"skill-typing-back/auth"
	"skill-typing-back/repository"

	"github.com/labstack/echo/v4"
)

// ログインユーザーの情報を取得する
func GetMe(c echo.Context) error {
	
	user := c.Get("user").(*auth.CognitoClaims)

	dbUser, err := repository.GetUser(user.Sub)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user info",
		})
	}

	return c.JSON(http.StatusOK, dbUser)
}