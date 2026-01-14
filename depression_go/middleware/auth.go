package middleware

import (
	"strings"

	"depression_go/pkg/response"
	"depression_go/pkg/token"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "令牌格式错误")
			c.Abort()
			return
		}

		// 提取令牌
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证令牌
		valid, claims := token.ValidateToken(tokenString)
		if !valid {
			response.Unauthorized(c, "无效的认证令牌")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}
