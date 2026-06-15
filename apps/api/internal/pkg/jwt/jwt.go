package jwt

import (
	"errors"
	"time"

	"file-manager-service/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType Token 类型
type TokenType string

const (
	TokenTypeUser  TokenType = "user"  // 用户 Token
	TokenTypeApp   TokenType = "app"   // 应用 Token
)

// Claims JWT 声明
type Claims struct {
	Type          TokenType `json:"type"`            // Token 类型：user 或 app
	UserID        int       `json:"user_id"`         // 用户 ID（用户 Token）
	Username      string    `json:"username"`        // 用户名（用户 Token）
	AppID         int       `json:"app_id"`          // 应用 ID（应用 Token）
	AppName       string    `json:"app_name"`        // 应用名称（应用 Token）
	AppIdentifier string    `json:"app_identifier"`  // 应用标识（应用 Token，用于文件存储路径）
	jwt.RegisteredClaims
}

var jwtSecret []byte

// Init 初始化 JWT
func Init(cfg *config.JWTConfig) {
	jwtSecret = []byte(cfg.Secret)
}

// GenerateUserToken 生成用户 Token
func GenerateUserToken(userID int, username string) (string, error) {
	claims := Claims{
		Type:     TokenTypeUser,
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(config.GlobalConfig.JWT.ExpireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateAppToken 生成应用 Token
func GenerateAppToken(appID int, appName string, appIdentifier string) (string, error) {
	claims := Claims{
		Type:          TokenTypeApp,
		AppID:         appID,
		AppName:       appName,
		AppIdentifier: appIdentifier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(config.GlobalConfig.JWT.ExpireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateToken 兼容旧代码的 Token 生成方法（用户 Token）
func GenerateToken(userID int, username string) (string, error) {
	return GenerateUserToken(userID, username)
}
