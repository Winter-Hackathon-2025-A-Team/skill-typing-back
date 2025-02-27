package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"skill-typing-back/auth"
	"skill-typing-back/model"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock作成
type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) CreateAnswer(c echo.Context, ansewer *model.Answer) error {
	args := m.Called()
	return args.Error(0)
}

func (m *repositoryMock) CreateCategory(c echo.Context, category *model.Category) error {
	args := m.Called()
	return args.Error(0)
}

func (m *repositoryMock) CreateChoice(c echo.Context, choice *model.Choice) error {
	args := m.Called()
	return args.Error(0)
}

func (m *repositoryMock) CreateQuestion(c echo.Context, question *model.Question) error {
	args := m.Called()
	return args.Error(0)
}

func (m *repositoryMock) CreateScore(c echo.Context, score *model.Score) error {
	args := m.Called(c, score)
	return args.Error(0)
}

func (m *repositoryMock) CreateUser(id string, name string, isAdmin bool) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *repositoryMock) GetAllQuestions(c echo.Context) (*[]model.Question, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Question), args.Error(1)
}

func (m *repositoryMock) GetCategoryByTitle(c echo.Context, title string) (model.Category, error) {
	args := m.Called()
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *repositoryMock) GetChoiceByContentAndQuestionId(c echo.Context, content string, questionId uint) (*model.Choice, error) {
	args := m.Called()
	return args.Get(0).(*model.Choice), args.Error(1)
}

func (m *repositoryMock) GetLatestScore(c echo.Context, userId string) (*model.Score, error) {
	args := m.Called(c, userId)
	return args.Get(0).(*model.Score), args.Error(1)
}

func (m *repositoryMock) GetQuestion(c echo.Context, id string) (*model.Question, error) {
	args := m.Called()
	return args.Get(0).(*model.Question), args.Error(1)
}

func (m *repositoryMock) GetUser(id string) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *repositoryMock) UpdateQuestion(c echo.Context, question *model.Question, columnName string, updateValue interface{}) error {
	args := m.Called()
	return args.Error(0)
}
func (m *repositoryMock) GetRandomQuestions(c echo.Context) ([]model.Question, error) {
	args := m.Called()
	return args.Get(0).([]model.Question), args.Error(1)

}

func (m *repositoryMock) GetQuestionsByCategory(c echo.Context, categoryID uint, limit int) ([]model.Question, map[uint]uint, error) {
	args := m.Called(categoryID, limit)
	return args.Get(0).([]model.Question), args.Get(1).(map[uint]uint), args.Error(2)
}

func (m *repositoryMock) GetChoicesByQuestionID(c echo.Context, questionID uint) ([]model.Choice, error) {
	args := m.Called(questionID)
	return args.Get(0).([]model.Choice), args.Error(1)
}

func (m *repositoryMock) GetChoiceByContent(c echo.Context, content string) (*model.Choice, error) {
    args := m.Called(c, content)
    if choice, ok := args.Get(0).(*model.Choice); ok {
        return choice, args.Error(1)
    }
    return nil, args.Error(1)
}

// ユーザー取得テスト
func TestGetMeSuccess(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// userIDを設定
	userID := "testSub"
	cognitoClaims := &auth.CognitoClaims{Sub: userID}
	c.Set("user", cognitoClaims)

	// モックリポジトリの設定
	expectedUser := &model.User{
		ID:        userID,
		Name:      "Test User",
		IsAdmin:   true,
		CreatedAt: time.Now().UTC(),
	}
	m := new(repositoryMock)
	m.On("GetUser", userID).Return(expectedUser, nil)

	// ハンドラーの作成
	h := New(m)

	// ハンドラーの実行
	err := h.GetMe(c)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// レスポンスの検証
	expectedCreatedAt := expectedUser.CreatedAt.Format("2006-01-02T15:04:05Z")
	expectedJSON := `{"name":"Test User","is_admin":true,"created_at":"` + expectedCreatedAt + `"}`
	assert.JSONEq(t, expectedJSON, strings.TrimSpace(rec.Body.String()))

	// モックが期待通りに呼び出されたか検証
	m.AssertExpectations(t)
}

// ユーザー取得時コンテキストが無い時のエラーテスト
func TestGetMeUserNotInContext(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	m := new(repositoryMock)
	h := New(m)

	h.GetMe(c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	expectedJSON := `{"error":"User information not found in context"}`
	assert.JSONEq(t, expectedJSON, strings.TrimSpace(rec.Body.String()))
}

// ユーザー取得が失敗する時のテスト
func TestGetMeUserFails(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// userIDを設定
	userID := "testSub"
	cognitoClaims := &auth.CognitoClaims{Sub: userID}
	c.Set("user", cognitoClaims)

	m := new(repositoryMock)
	m.On("GetUser", userID).Return(
		(*model.User)(nil), 
		errors.New("Database error"),
		)

	h := New(m)

	err := h.GetMe(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	expectedJSON := `{"error": "Failed to get user info"}`
	assert.JSONEq(t, expectedJSON, strings.TrimSpace(rec.Body.String()))
}

func TestGetLatestScoreSucess(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/scores/latest", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// userIDを設定
	userID := "testSub"

	// mockが返却するscore
	maxScore := model.Score{
		ID:        33,
		UserID:    userID,
		Score:     77,
		CreatedAt: time.Date(2025, 1, 1, 12, 00, 00, 00, time.Local),
	}

	// 期待されるレスポンス
	resp := `{"score":"77"}`

	// authMiddleware の処理
	cognitoClaims := &auth.CognitoClaims{
		Sub:      userID,
		TokenUse: "access",
	}
	c.Set("user", cognitoClaims)

	// mock作成
	m := new(repositoryMock)
	// handler
	h := New(m)

	m.On("GetLatestScore", c, userID).Return(&maxScore, nil)

	if assert.NoError(t, h.GetLatestScore(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, resp, strings.TrimSpace(rec.Body.String()))
	}

}

func TestCreateScoreSuccess(t *testing.T) {
	// Setup
	e := echo.New()
	recScore := `{"score": 88}`
	req := httptest.NewRequest(http.MethodPost, "/scores", strings.NewReader(recScore))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// userIDを設定
	userID := "testSub"

	// 期待されるレスポンス
	resp := `{"message":"registered successfully"}`

	// authMiddleware の処理
	cognitoClaims := &auth.CognitoClaims{
		Sub:      userID,
		TokenUse: "access",
	}
	c.Set("user", cognitoClaims)

	// mock作成
	m := new(repositoryMock)

	// mockに入力されるScore
	score := model.Score{
		UserID: userID,
		Score:  88,
	}

	m.On("CreateScore", c, &score).Return(nil)

	// handler
	h := New(m)

	if assert.NoError(t, h.CreateScore(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, resp, strings.TrimSpace(rec.Body.String()))
	}

}
