package utils

import (
	"errors"
	"fmt"
	"time"
	"todolist/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret              []byte
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
)

// InitJWT 初始化 JWT 配置
func InitJWT() {

	jwtSecret = config.JWTSecret

	// AccessToken 有效期：15 分钟
	accessTokenExpiration = config.JWTExpiration

	// RefreshToken 有效期：7 天
	refreshTokenExpiration = config.JWTExpiration + time.Hour
}

// Claims JWT 声明
type Claims struct {
	UserID    uint   `json:"userId"`
	Username  string `json:"username"`
	TokenType string `json:"tokenType"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

// TokenPair Token 对
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // AccessToken 过期时间（秒）
	RefreshIn    int64  `json:"refreshIn"` // RefreshToken 过期时间（秒）
}

// GenerateTokenPair 生成 Token 对（AccessToken + RefreshToken）
func GenerateTokenPair(userID uint, username string) (*TokenPair, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	now := time.Now()

	// 生成 AccessToken
	accessToken, err := generateToken(userID, username, "access", now, accessTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成 RefreshToken
	refreshToken, err := generateToken(userID, username, "refresh", now, refreshTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTokenExpiration.Seconds()),
		RefreshIn:    int64(refreshTokenExpiration.Seconds()),
	}, nil
}

// JWT Claims

func generateToken(userID uint, username, tokenType string, now time.Time, expiration time.Duration) (string, error) {
	claims := Claims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "todo-api",
			Subject:   fmt.Sprintf("%d", userID),
			ID:        fmt.Sprintf("%s-%d-%d", tokenType, userID, now.Unix()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 检查 Token 是否过期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshAccessToken 使用 RefreshToken 刷新 AccessToken
func RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := ParseToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 验证是 RefreshToken
	if claims.TokenType != "refresh" {
		return nil, errors.New("token is not a refresh token")
	}

	// 生成新的 Token 对
	return GenerateTokenPair(claims.UserID, claims.Username)
}

// ValidateAccessToken 验证 AccessToken
func ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, errors.New("token is not an access token")
	}

	return claims, nil
}

// GetTokenExpiration 获取 Token 过期时间
func GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, errors.New("token has no expiration")
	}

	return claims.ExpiresAt.Time, nil
}

// IsTokenExpiringSoon 检查 Token 是否即将过期（5 分钟内）
func IsTokenExpiringSoon(tokenString string) bool {
	expiresAt, err := GetTokenExpiration(tokenString)
	if err != nil {
		return true
	}

	// 如果在 5 分钟内过期，返回 true
	return time.Until(expiresAt) < time.Minute*5
}
