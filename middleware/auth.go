package middware

import (
	"errors"
	"net/http"

	"todolist/utils"

	"github.com/gin-gonic/gin"
)

// JWT 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Token
		authHeader := c.GetHeader("x-access-token")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "未提供认证令牌",
				"error":   errors.New("未提供认证令牌"),
			})
			c.Abort()
			return
		}

		// 检查 Token 格式：Bearer <token>
		// parts := strings.SplitN(authHeader, " ", 2)
		// log.Printf("parts %v", parts)
		// if !(len(parts) == 2 && parts[0] == "Bearer") {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"success": false,
		// 		"message": "认证令牌格式错误",
		// 		"error":   errors.New("认证令牌格式错误"),
		// 	})
		// 	c.Abort()
		// 	return
		// }

		// 解析 Token
		claims, err := utils.ParseToken(authHeader)
		if err != nil {

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "认证令牌无效或已过期",
			})

			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("userName", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

// 从上下文获取用户 ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

// 从上下文获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	userName, exists := c.Get("userName")
	if !exists {
		return "", false
	}

	name, ok := userName.(string)
	return name, ok
}
