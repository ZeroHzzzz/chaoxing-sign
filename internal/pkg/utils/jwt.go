package utils

import (
	"chaoxing/internal/globals/config"
	"chaoxing/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT自定义声明结构
type CustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(user *models.User) (string, error) {
	// 获取JWT配置
	secretKey := config.Config.GetString("jwt.secret_key")
	expireTime := config.Config.GetInt("jwt.expire_time")

	// 创建自定义声明
	claims := CustomClaims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expireTime))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_token",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名令牌并获取完整的编码令牌
	return token.SignedString([]byte(secretKey))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 获取JWT密钥
	secretKey := config.Config.GetString("jwt.secret_key")

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
