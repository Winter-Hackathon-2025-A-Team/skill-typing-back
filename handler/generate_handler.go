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
	"AWSアソシエイト":  true,
}

// AI を使って問題を生成し、データベースに保存
func (h *ApiHandler) GenerateQuizHandler(c echo.Context) error {
	categoryName := strings.TrimSpace(c.QueryParam("category"))

	// バリデーション
	if categoryName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "カテゴリーを指定してください"})
	}

	// // 新しいカテゴリーがリストになければ追加
	// if !validCategories[categoryName] {
	// 	validCategories[categoryName] = true
	// }

	// カテゴリーを取得 or 作成
	var category model.Category
	category, err := h.repo.GetCategoryByTitle(c, categoryName)
	// カテゴリーの存在確認、なければ作成
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
`, categoryName)

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

	// 正解の選択肢のインデックスを特定
	answerIndex := -1
	for i, choice := range quizData.Choices {
		if choice == quizData.Answer {
			answerIndex = i
		}
	}

	if answerIndex == -1 {
		log.Println("正解の選択肢が見つかりませんでした")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "正解の選択肢を特定できませんでした"})
	}

	// フロントエンド用のレスポンス形式に変換
	type ChoiceResponse struct {
		Content string `json:"content"`
		Description string `json:"description"`
	}

	choices := make([]ChoiceResponse, 4)
	for i := 0; i < 4; i++ {
		choices[i] = ChoiceResponse{
			Content: quizData.Choices[i],
			Description: quizData.Descriptions[i],
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"title":        quizData.Title,
		"content":      quizData.Content,
		"category":     category.Title,
		"choices":      choices,
		"answer_index": answerIndex,
		"explanation":  quizData.Explanation,
	})
}

// 生成されたクイズをDBに保存する
func (h *ApiHandler) SaveQuizHandler(c echo.Context) error {
	// リクエストデータの取得
	var requestData struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		Category    string `json:"category"`
		Choices     []struct {
			Content     string `json:"content"`
			Description string `json:"description"`
		} `json:"choices"`
		AnswerIndex int `json:"answer_index"`
		Explanation string `json:"explanation"`
	}

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なリクエストフォーマットです"})
	}

	// バリデーション
	if requestData.Title == "" || requestData.Content == "" || requestData.Category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "タイトル、問題文、カテゴリーは必須です"})
	}
	if len(requestData.Choices) != 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "選択肢は4つ必要です"})
	}
	if requestData.AnswerIndex < 0 || requestData.AnswerIndex > 3 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "正解のインデックス範囲が無効です"})
	}

	// ユーザーIDの取得
	userId := c.Get("user").(*auth.CognitoClaims).Sub

	var category model.Category
	category, err := h.repo.GetCategoryByTitle(c, requestData.Category)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		newCategory := model.Category{Title: requestData.Category}
		if err := h.repo.CreateCategory(c, &newCategory); err != nil {
			log.Println("❌ カテゴリ作成エラー:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "カテゴリの作成に失敗しました"})
		}
		category = newCategory
	} else if err != nil {
		log.Println("❌ カテゴリ取得エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "カテゴリの取得に失敗しました"})
	}
	// 選択肢を作成
	var choice1, choice2, choice3, choice4 model.Choice
	var choices = []model.Choice{}
	var answerChoiceID uint

	for i, choice := range requestData.Choices {
		// 既存の選択肢を検索
		existingChoice, err := h.repo.GetChoiceByContent(c, choice.Content)

		var newChoice model.Choice
		if err == nil {
			// 既存の選択肢があれば使用する
			newChoice = existingChoice
			log.Printf("既存の選択肢を利用: ID=%d, Content=%s", newChoice.ID, newChoice.Content)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// 見つからなかった場合は新しく作成する
			newChoice := model.Choice{
				Content: choice.Content,
				Description: choice.Description,
			}
				
			if err := h.repo.CreateChoice(c, &newChoice); err != nil {
				log.Println("選択肢の保存エラー:", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の保存に失敗しました"})
			}
			log.Printf("新しいあたらしい選択肢を作成: ID=%d, Content=%s", newChoice.ID, newChoice.Content)
		} else {
			// その他のエラーの場合
			log.Println("選択肢検索エラー:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "選択肢の検索に失敗しました"})
		}

		choices = append(choices, newChoice)

		switch i {
		case 0:
			choice1 = newChoice
		case 1:
			choice2 = newChoice
		case 2:
			choice3 = newChoice
		case 3:
			choice4 = newChoice
		}
	}

	answerChoiceID = choices[requestData.AnswerIndex].ID

	// 問題を保存
	question := model.Question{
		Title:      requestData.Title,
		Content:    requestData.Content,
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
		Explanation: requestData.Explanation,
	}
	if err := h.repo.CreateAnswer(c, &answer); err != nil {
		log.Println("正解データの保存エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "正解データの保存に失敗しました"})
	}

	// **レスポンス**
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "問題が正常に保存されました",
		"question": question,
		"category": category,
		"choices":  choices,
		"answer":   answer,
	})
}
