package handler

import (
	"net/http"
	"skill-typing-back/auth"

	"github.com/labstack/echo/v4"
)

// ログインユーザーの情報を取得する
func (h *ApiHandler) GetMe(c echo.Context) error {

	user, ok := c.Get("user").(*auth.CognitoClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "User information not found in context")
	}

	dbUser, err := h.repo.GetUser(user.Sub)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user info",
		})
	}

	type UserResponse struct {
		Name string `json:"name"`
		IsAdmin bool `json:"is_admin"`
		CreatedAt string `json:"created_at"`
	}

	response := UserResponse {
		Name: dbUser.Name,
		IsAdmin: dbUser.IsAdmin,
		CreatedAt: dbUser.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return c.JSON(http.StatusOK, response)
}
