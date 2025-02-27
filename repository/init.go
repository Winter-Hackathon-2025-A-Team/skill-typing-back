package repository

import (
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type DbRepository struct {
	db *gorm.DB
}

type DbRepositoryInterface interface {
	CreateAnswer(c echo.Context, ansewer *model.Answer) error
	CreateCategory(c echo.Context, category *model.Category) error
	CreateChoice(c echo.Context, choice *model.Choice) error
	CreateQuestion(c echo.Context, question *model.Question) error
	CreateScore(c echo.Context, score *model.Score) error
	CreateUser(id string, name string, isAdmin bool) (*model.User, error)
	GetAllQuestions(c echo.Context) (*[]model.Question, error)
	GetCategoryByTitle(c echo.Context, title string) (model.Category, error)
	GetChoiceByContentAndQuestionId(c echo.Context, content string, questionId uint) (*model.Choice, error)
	GetChoiceByContent(c echo.Context, content string) (*model.Choice, error)
	GetLatestScore(c echo.Context, userId string) (*model.Score, error)
	GetQuestion(c echo.Context, id string) (*model.Question, error)
	GetUser(id string) (*model.User, error)
	UpdateQuestion(c echo.Context, question *model.Question, columnName string, updateValue interface{}) error
	GetQuestionsByCategory(c echo.Context, categoryID uint, limit int) ([]model.Question, map[uint]uint, error)
	GetChoicesByQuestionID(c echo.Context, questionID uint) ([]model.Choice, error)
}

func New(db *gorm.DB) *DbRepository {

	r := DbRepository{
		db: db,
	}
	return &r
}
