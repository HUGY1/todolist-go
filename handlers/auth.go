package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"todolist/config"
	"todolist/models"
	"todolist/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

// 登录请求
type LoginRequest struct {
	Username string `json:"userName" binding:"required"` // 可以是用户名或邮箱
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresIn    int64       `json:"expiresIn"` // AccessToken 过期时间（秒）
	RefreshIn    int64       `json:"refreshIn"` // RefreshToken 过期时间（秒）
	TokenType    string      `json:"tokenType"` // "Bearer"
	User         models.User `json:"user"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenResponse 刷新 Token 响应
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	RefreshIn    int64  `json:"refreshIn"`
	TokenType    string `json:"tokenType"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	var user models.User

	err := config.GetDB().Model(&models.User{}).Where("user_name = ?", req.Username).First((&user)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "用户名或密码错误",
				"error":   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "查询用户失败",
				"error":   err.Error(),
			})
			return
		}
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户名或密码错误",
			"error":   errors.New("用户名或密码错误"),
		})
		return
	}

	// 生成 Token 对
	tokenPair, err := utils.GenerateTokenPair(user.UserID, user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成令牌失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("✅ 登录成功: UserID=%d", user.UserID)

	var data = LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		RefreshIn:    tokenPair.RefreshIn,
		TokenType:    "Bearer",
		User:         user,
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// 刷新 Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未提供认证令牌",
			"error":   errors.New("未提供认证令牌"),
		})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "认证令牌格式错误",
			"error":   errors.New("认证令牌格式错误"),
		})

		return
	}
	newToken, err := utils.RefreshAccessToken(parts[1])
	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "刷新令牌失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    newToken,
	})
}

// 登出（可选，前端删除 Token 即可）
func (h *AuthHandler) Logout(c *gin.Context) {
	// 这里可以实现 Token 黑名单机制
	// 简单实现：前端删除 Token 即可
	// utils.SuccessMsg(c, nil, "登出成功")
}
