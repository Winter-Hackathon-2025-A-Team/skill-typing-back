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
	log.Println("Starting token verification...")

	decomposeUnverifiedJwt, err := utils.DecomposeUnverifiedJwt(token)
	if err != nil {
		log.Printf("Failed to decompose JWT: %v", err)
		return nil, err
	}

	jwk, err := utils.GetJwk(decomposeUnverifiedJwt, c.jwksUri, c.cache)
	if err != nil {
		log.Printf("Failed to get JWT: %v", err)
		return nil, err
	}

	err = utils.VerifyDecomposedJwt(decomposeUnverifiedJwt, c.issuer, c.tokenUse, jwk.Alg)
	if err != nil {
		log.Printf("Failed to verify decomposed JWT: %v", err)
		return nil, err
	}

	// validToken, err := utils.ValidateJwt(token, jwk)
	// if err != nil {
	// 	log.Printf("Failed to validate JWT: %v", err)
	// 	return nil, err
	// }

	// log.Printf("Claims type after validation: %T", validToken.Claims)
	// log.Printf("Claims content: %+v", validToken.Claims)

	// return validToken.Claims, nil

	validToken, err := utils.ValidateJwt(token, jwk)
	if err != nil {
		log.Printf("Failed to validate JWT: %v", err)
		return nil, err
	}
	log.Printf("Claims type after validation: %T", validToken.Claims)
	log.Printf("Claims content: %+v", validToken.Claims)

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

			// デバッグ用: トークンの内容を確認
			log.Printf("received token: %s", token)

			// トークンを検証
			claims, err := a.verifier.Verify(token)
			if err != nil {
				log.Printf("Token verification failed: %v", err)
				return echo.ErrUnauthorized
			}

			log.Printf("Claims type before cast: %T", claims)
			log.Printf("Claims content before cast: %+v", claims)

			mapClaims, ok := claims.(jwt.MapClaims)
			if !ok {
				log.Printf("Failed to cast to MapClaims")
				return echo.ErrInternalServerError
			}

			// クレームをCognitoClaimsにマッピング
			cognitoClaims := &CognitoClaims{
				Sub: mapClaims["sub"].(string),
				TokenUse: mapClaims["token_use"].(string),
				Username: mapClaims["username"].(string),
			}

			// コンテキストにユーザー情報を保存
			c.Set("user", cognitoClaims)

			return next(c)
		}
	}
}