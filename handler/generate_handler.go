package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"skill-typing-back/auth"
	"skill-typing-back/model"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// 固定のカテゴリリスト
var validCategories = map[string]bool{
	"基本情報技術者":   true,
	"AWSアソシエイト": true,
}

// AI を使って問題を生成し、データベースに保存
func (h *ApiHandler) GenerateQuizHandler(c echo.Context) error {
	topic := c.QueryParam("topic")
	categoryName := strings.TrimSpace(c.QueryParam("category"))

	// バリデーション
	if topic == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "トピックを指定してください"})
	}
	if categoryName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "カテゴリーを指定してください"})
	}

	// 新しいカテゴリーがリストになければ追加
	if !validCategories[categoryName] {
		validCategories[categoryName] = true
	}

	// カテゴリーを取得 or 作成
	var category model.Category
	category, err := h.repo.GetCategoryByTitle(c, categoryName)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		newCategory := model.Category{Title: categoryName}
		if err := h.repo.CreateCategory(c, &newCategory); err != nil {
			log.Println("❌ カテゴリ作成エラー:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "カテゴリの作成に失敗しました"})
		}
		category = newCategory
	} else if err != nil {
		log.Println("❌ カテゴリ取得エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "カテゴリの取得に失敗しました"})
	}

	// OpenAI API キー取得
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("OpenAI APIキーが設定されていません")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "OpenAI APIキーが設定されていません"})
	}

	client := openai.NewClient(apiKey)

	// AI へリクエスト
	prompt := fmt.Sprintf(`
「%s」に関するクイズを作成してください。
- 問題のタイトル（短いフレーズ）
- 問題文
- 4つの選択肢
- 各選択肢の解説
- 正解の選択肢（選択肢の中から1つ）
- 正解の詳細な解説

出力形式は下記のjsonの形式でお願いいたします。

{
	"title" : "問題のタイトル",
    "question":"問題文",
    "choices":["選択肢1","選択肢2","選択肢3","選択肢4"],
    "descriptions":["選択肢1の解説","選択肢2の解説","選択肢3の解説","選択肢4の解説"],
	"answer": "正解の選択肢",
	"explanation": "正解の詳細な解説"
}
`, topic)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "あなたは試験問題作成の専門家です。"},
				{Role: "user", Content: prompt},
			},
		},
	)
	if err != nil || len(resp.Choices) == 0 {
		log.Println("OpenAIからのレスポンスエラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AIからの応答が無効です"})
	}

	// AI のレスポンスをパース
	var quizData struct {
		Title        string   `json:"title"`
		Content      string   `json:"question"`
		Choices      []string `json:"choices"`
		Descriptions []string `json:"descriptions"`
		Answer       string   `json:"answer"`
		Explanation  string   `json:"explanation"`
	}

	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &quizData); err != nil {
		log.Println("AIのレスポンス解析エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AIのレスポンス解析に失敗しました"})
	}

	userId := c.Get("user").(*auth.CognitoClaims).Sub

	// 4つの選択肢を保存
	var choice1, choice2, choice3, choice4 model.Choice
	var choices = []model.Choice{}
	var answerChoiceID uint

	// 選択肢を作成
	for i, choiceContent := range quizData.Choices {
		choice := model.Choice{
			Content: choiceContent,
			Description: quizData.Descriptions[i],
		}

		if err := h.repo.CreateChoice(c, &choice); err != nil {
			log.Println("選択肢の保存エラー:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の保存に失敗しました"})
		}

		choices = append(choices, choice)

		// 正解の選択肢IDを記録
		if choiceContent == quizData.Answer {
			answerChoiceID = choice.ID
		}

		// 各選択肢をそれぞれの変数に格納
		switch i {
		case 0:
			choice1 = choice
		case 1:
			choice2 = choice
		case 2:
			choice3 = choice
		case 3:
			choice4 = choice
		}
	}

	// 問題を保存
	question := model.Question{
		Title:      quizData.Title,
		Content:    quizData.Content,
		CategoryID: category.ID,
		UserID:     userId,
		Choice1ID: choice1.ID,
		Choice2ID: choice2.ID,
		Choice3ID: choice3.ID,
		Choice4ID: choice4.ID,
	}

	if err := h.repo.CreateQuestion(c, &question); err != nil {
		log.Println("問題の保存エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の保存に失敗しました"})
	}

	// 正解データを `answers` テーブルに保存
	answer := model.Answer{
		QuestionID:  question.ID,
		ChoiceID:    answerChoiceID,
		Explanation: quizData.Explanation,
	}
	if err := h.repo.CreateAnswer(c, &answer); err != nil {
		log.Println("正解データの保存エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "正解データの保存に失敗しました"})
	}

	// **レスポンス**
	return c.JSON(http.StatusOK, map[string]interface{}{
		"question": question,
		"category": category,
		"choices":  choices,
		"answer":   answer,
	})
}
