// Cognito認証関連
package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"skill-typing-back/repository"

	// "github.com/99designs/gqlgen/codegen/config"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jhosan7/cognito-jwt-verify/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
	repo     *repository.DbRepository
}

// JWTクレームの構造体
type CognitoClaims struct {
	jwt.RegisteredClaims
	Sub       string `json:"sub"`
	Iss       string `json:"iss"`
	Version   int    `json:"version"`
	ClientID  string `json:"client_id"`
	OriginJti string `json:"origin_jti"`
	EventID   string `json:"event_id"`
	TokenUse  string `json:"token_use"` // トークンタイプ（"access" of "id"）
	Scope     string `json:"scope"`
	AuthTime  int64  `json:"auth_time"`
	Username  string `json:"username"`
}

// CognitoJwtVerifierの構造体
type CognitoJwtVerifier struct {
	issuer   string
	jwksUri  string
	tokenUse string
	cache    *utils.Cache
}

type Config struct {
	UserPoolId string
	TokenUse   string
	ClientId   string
}

func NewCognitoUserService() (*CognitoUserService, error) {
	// SDkの設定を初期化
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}
	// Cognitoクライアントを作成
	client := cognitoidentityprovider.NewFromConfig(cfg)

	return &CognitoUserService{
		client: client,
	}, nil
}

// 新しいCognitoAuth インスタンス作成
func NewCognitoAuth(repo *repository.DbRepository) (*CognitoAuth, error) {
	config := Config{
		UserPoolId: os.Getenv("COGNITO_USER_POOL_ID"),
		TokenUse:   "access",
		ClientId:   os.Getenv("COGNITO_CLIENT_ID"),
	}

	verifier, err := Create(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cognito verifier: %w", err)
	}

	return &CognitoAuth{
		verifier: verifier,
		repo:     repo,
	}, nil
}

// JWTVerifierの作成
func Create(config Config) (CognitoJwtVerifier, error) {
	issuer, jwksUri, err := utils.ParseUserPoolId(config.UserPoolId)
	if err != nil {
		return CognitoJwtVerifier{}, err
	}
	return CognitoJwtVerifier{
		issuer:   issuer,
		jwksUri:  jwksUri,
		tokenUse: config.TokenUse,
		cache:    utils.NewCache(),
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
				Sub:      mapClaims["sub"].(string),
				TokenUse: mapClaims["token_use"].(string),
			}

			sub := cognitoClaims.Sub

			// DBでユーザーを検索
			user, err := a.repo.GetUser(sub)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// トークンからユーザー情報を取得
					userInfo, err := cognitoService.GetUserInfo(c.Request().Context(), token)
					if err != nil {
						log.Printf("Failed to get user info from Cognito: %v", err)
						return echo.ErrInternalServerError
					}
					var userName string
					var isAdmin bool
					for _, attr := range userInfo.UserAttributes {
						switch *attr.Name {
						case "name":
							userName = *attr.Value
						case "custom:isAdmin":
							isAdmin = *attr.Value == "true"
						}
					}
					// 新規ユーザーを作成
					user, err = a.repo.CreateUser(sub, userName, isAdmin)
					if err != nil {
						log.Printf("Failed to create user: %v", err)
						return echo.ErrInternalServerError
					}
				} else {
					log.Printf("Failed to get user: $v", err)
					return echo.ErrInternalServerError
				}
			}
			// Cognitoからユーザー情報を取得してログ出力
			// userInfo, err := cognitoService.GetUserInfo(c.Request().Context(), token)
			// if err != nil {
			// 	log.Printf("Failed to get user info from Cognito: %v", err)
			// } else {
			// 	log.Printf("Cognito User Attributes:")
			// 	for _, attr :=range userInfo.UserAttributes {
			// 		if attr.Name != nil && attr.Value != nil {
			// 			log.Printf(" %s: %s", *attr.Name, *attr.Value)
			// 		}
			// 	}
			// }

			// log.Printf("Authentication successful - UserSub: %s", cognitoClaims.Sub)

			c.Set("user", cognitoClaims)
			c.Set("dbUser", user)

			return next(c)
		}
	}
}
