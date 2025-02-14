package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"skill-typing-back/db"
	"skill-typing-back/model"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

// 固定のカテゴリリスト
var validCategories = map[string]bool{
	"基本情報技術者":   true,
	"AWSアソシエイト": true,
}

// AI を使って問題を生成し、データベースに保存
func GenerateQuizHandler(c echo.Context) error {
	topic := c.QueryParam("topic")
	categoryName := strings.TrimSpace(c.QueryParam("category"))

	// バリデーション
	if topic == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "トピックを指定してください"})
	}
	if categoryName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "カテゴリーを指定してください"})
	}

	// DBインスタンスを取得
	dbConn := db.GetDB()
	if dbConn == nil {
		log.Println("データベース接続が確立されていません")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "データベース接続エラー"})
	}

	// 新しいカテゴリーがリストになければ追加
	if !validCategories[categoryName] {
		validCategories[categoryName] = true
	}

	// カテゴリーを取得 or 作成
	var category model.Category
	if err := dbConn.Where("title = ?", categoryName).First(&category).Error; err != nil {
		category = model.Category{Title: categoryName}
		if err := dbConn.Create(&category).Error; err != nil {
			log.Println("❌ カテゴリ作成エラー:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "カテゴリの作成に失敗しました"})
		}
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
- 正解（選択肢の中から1つ）
- 正解の詳細な解説
JSON 形式で出力してください。
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

	// 問題を保存
	question := model.Question{
		Title:      quizData.Title,
		Content:    quizData.Content,
		CategoryID: category.ID,
		UserID:     "system",
	}
	if err := dbConn.Create(&question).Error; err != nil {
		log.Println("問題の保存エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "問題の保存に失敗しました"})
	}

	// 選択肢を保存
	var answerID uint
	var choices []model.Choice
	for i, choice := range quizData.Choices {
		choiceModel := model.Choice{
			QuestionID:  question.ID,
			Content:     choice,
			Description: quizData.Descriptions[i],
		}
		if err := dbConn.Create(&choiceModel).Error; err != nil {
			log.Println("選択肢の保存エラー:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の保存に失敗しました"})
		}
		choices = append(choices, choiceModel)

		// 正解の選択肢IDを取得
		if choice == quizData.Answer {
			answerID = choiceModel.ID
		}
	}

	// 正解データを `answers` テーブルに保存
	answer := model.Answer{
		QuestionID:  question.ID,
		ChoiceID:    answerID,
		Explanation: quizData.Explanation,
	}
	if err := dbConn.Create(&answer).Error; err != nil {
		log.Println("正解データの保存エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "正解データの保存に失敗しました"})
	}

	// **問題の正解IDを更新**
	if err := dbConn.Model(&question).Update("AnswerID", answerID).Error; err != nil {
		log.Println("正解IDの更新エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "正解の設定に失敗しました"})
	}

	// **レスポンス**
	return c.JSON(http.StatusOK, map[string]interface{}{
		"question": question,
		"category": category,
		"choices":  choices,
		"answer":   answer,
	})
}
