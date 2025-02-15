package repository

import (
	"fmt"
	"skill-typing-back/model"

	"github.com/labstack/echo/v4"
)

func (r *DbRepository) CreateQuestion(c echo.Context, question *model.Question) error {

	if err := r.db.Create(question).Error; err != nil {
		return fmt.Errorf("failed to create question : %v", err)
	}
	return nil
}

func (r *DbRepository) GetQuestion(c echo.Context, id string) (*model.Question, error) {

	var question model.Question

	if err := r.db.First(&question, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get question : %v", err)
	}

	return &question, nil
}

func (r *DbRepository) GetAllQuestions(c echo.Context) (*[]model.Question, error) {

	var questions []model.Question

	if err := r.db.Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to get all question : %v", err)
	}

	return &questions, nil
}

func (r *DbRepository) UpdateQuestion(c echo.Context, question *model.Question, columnName string, updateValue interface{}) error {

	if err := r.db.Model(question).Update(columnName, updateValue).Error; err != nil {
		return fmt.Errorf("failed to update question :%v", err)

	}

	return nil

}
