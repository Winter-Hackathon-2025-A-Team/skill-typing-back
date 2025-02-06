// Cognito認証関連
package auth

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jhosan7/cognito-jwt-verify/utils"
	"github.com/labstack/echo/v4"
)

// Cognito API クライアントを保持する構造体
type CognitoUserService struct {
	client *cognitoidentityprovider.Client
}

// ユーザー情報取得メソッド
func (s *CognitoUserService) GetUserInfo(ctx context.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, error) {
	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: &accessToken,
	}
	return s.client.GetUser(ctx, input)
}

// Cognitoの設定やクレームを扱う構造体
type CognitoAuth struct {
	verifier CognitoJwtVerifier
}

// JWTクレームの構造体
type CognitoClaims struct {
	jwt.RegisteredClaims
	Sub string `json:"sub"`
	Iss string `json:"iss"`
	Version int `json:"version"`
	ClientID string `json:"client_id"`
	OriginJti string `json:"origin_jti"`
	EventID string `json:"event_id"`
	TokenUse string `json:"token_use"` // トークンタイプ（"access" of "id"）
	Scope string `json:"scope"`
	AuthTime int64 `json:"auth_time"`
	Username string `json:"username"`
}

// CognitoJwtVerifierの構造体
type CognitoJwtVerifier struct {
	issuer string
	jwksUri string
	tokenUse string
	cache *utils.Cache
}

type Config struct {
	UserPoolId string
	TokenUse string
	ClientId string
}

func NewCognitoUserService() (*CognitoUserService, error) {
	// 一時的な実装としてからのクライアントを返す
	return &CognitoUserService{
		client: &cognitoidentityprovider.Client{},
	}, nil 
}

// 新しいCognitoAuth インスタンス作成
func NewCognitoAuth() (*CognitoAuth, error) {
	config := Config{
		UserPoolId: os.Getenv("COGNITO_USER_POOL_ID"),
		TokenUse: "access",
		ClientId: os.Getenv("COGNITO_CLIENT_ID"),
	}

	verifier, err := Create(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cognito verifier: %w", err)
	}

	return &CognitoAuth{
		verifier: verifier,
	}, nil
}

// JWTVerifierの作成
func Create(config Config) (CognitoJwtVerifier, error) {
	issuer, jwksUri, err := utils.ParseUserPoolId(config.UserPoolId)
	if err != nil {
		return CognitoJwtVerifier{}, err
	}
	return CognitoJwtVerifier{
		issuer: issuer,
		jwksUri: jwksUri,
		tokenUse: config.TokenUse,
		cache: utils.NewCache(),
	}, nil
}

// トークン検証
func (c CognitoJwtVerifier) Verify(token string) (jwt.Claims, error) {
	decomposeUnverifiedJwt, err := utils.DecomposeUnverifiedJwt(token)
	if err != nil {
		return nil, err
	}

	jwk, err := utils.GetJwk(decomposeUnverifiedJwt, c.jwksUri, c.cache)
	if err != nil {
		return nil, err
	}

	err = utils.VerifyDecomposedJwt(decomposeUnverifiedJwt, c.issuer, c.tokenUse, jwk.Alg)
	if err != nil {
		return nil, err
	}

	validToken, err := utils.ValidateJwt(token, jwk)
	if err != nil {
		return nil, err
	}

	return validToken.Claims, nil
}

// 認証ミドルウェア
func (a *CognitoAuth) AuthMiddleware(cognitoService *CognitoUserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorizationヘッダーからトークンを取得
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.ErrUnauthorized
			}

			// "Bearer"を除去
			token := auth[7:]

			// トークンを検証
			claims, err := a.verifier.Verify(token)
			if err != nil {
				return echo.ErrUnauthorized
			}

			mapClaims, ok := claims.(jwt.MapClaims)
			if !ok {
				return echo.ErrInternalServerError
			}

			// クレームをCognitoClaimsにマッピング
			cognitoClaims := &CognitoClaims{
				Sub: mapClaims["sub"].(string),
				TokenUse: mapClaims["token_use"].(string),
				Username: mapClaims["username"].(string),
			}

			// Cognitoからユーザー情報を取得してログ出力
			userInfo, err := cognitoService.GetUserInfo(c.Request().Context(), token)
			if err != nil {
				log.Printf("Failed to get user info from Cognito: %v", err)
			} else {
				log.Printf("Cognito User Attributes:")
				for _, attr :=range userInfo.UserAttributes {
					if attr.Name != nil && attr.Value != nil {
						log.Printf(" %s: %s", *attr.Name, *attr.Value)
					}
				}
			}

			log.Printf("Authentication successful - User: %s", cognitoClaims.Username)

			// コンテキストにユーザー情報を保存
			c.Set("user", cognitoClaims)

			return next(c)
		}
	}
}