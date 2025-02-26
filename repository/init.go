package repository

import (
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type DbRepository struct {
	db *gorm.DB
}

// GetChoicesByQuestionID implements DbRepositoryInterface.
func (r *DbRepository) GetChoicesByQuestionID(c echo.Context, questionID uint) ([]model.Choice, error) {
	var choices []model.Choice
	if err := r.db.Where("question_id = ?", questionID).Find(&choices).Error; err != nil {
		return nil, err
	}
	return choices, nil
}

// GetQuestionsByCategory implements DbRepositoryInterface.
func (r *DbRepository) GetQuestionsByCategory(c echo.Context, categoryID uint, limit int) ([]model.Question, error) {
	var questions []model.Question
	if err := r.db.Where("category_id = ?", categoryID).Limit(limit).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

// GetRandomQuestions implements DbRepositoryInterface.
func (r *DbRepository) GetRandomQuestions(c echo.Context) ([]model.Question, error) {
	panic("unimplemented")
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
	GetRandomQuestions(c echo.Context) ([]model.Question, error)
	GetQuestionsByCategory(c echo.Context, categoryID uint, limit int) ([]model.Question, error)
	GetChoicesByQuestionID(c echo.Context, questionID uint) ([]model.Choice, error)
}

func New(db *gorm.DB) *DbRepository {

	r := DbRepository{
		db: db,
	}
	return &r
}
