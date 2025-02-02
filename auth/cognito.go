// Cognito認証関連
package auth

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jhosan7/cognito-jwt-verify/utils"
	"github.com/labstack/echo/v4"
)

// Cognitoの設定やクレームを扱う構造体
type CognitoAuth struct {
	verifier CognitoJwtVerifier
}

// JWTクレームの構造体
type CognitoClaims struct {
	Sub string `json:"sub"`
	Email string `json:"email"`
	Name string `json:"name"`
	TokenUse string `json:"token_use"` // トークンタイプ（"access" of "id"）
	AuthTime int64 `json:"auth_time"`
	jwt.RegisteredClaims
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
func (a *CognitoAuth) AuthMiddleware() echo.MiddlewareFunc {
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
			// claims, err := a.verifier.Verify(token)
			// if err != nil {
			// 	log.Printf("Token verification failed: %v", err)
			// 	return echo.ErrUnauthorized
			// }

			// クレームをCognitoClaimsにマッピング
			// cognitoClaims, ok := claims.(*CognitoClaims)
			// if !ok {
			// 	log.Printf("Failed to cast claims to CognitoClaims")
			// 	return echo.ErrInternalServerError
			// }
			// ユーザー情報をログに出力
			// log.Printf("Authenticated user - Sub: %s, Email: %s, Name: %s",
			// 	cognitoClaims.Sub,
			// )

			// コンテキストにユーザー情報を保存
			// c.Set("user", cognitoClaims)

			// デバッグ用: トークンの内容を確認
			log.Printf("Received token: %s", token)

			return next(c)
		}
	}
}