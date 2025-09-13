package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"zf-server/pkg/common/auth"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"operationID": "auth_middleware",
				"code":        401,
				"msg":         "Authorization header is required",
				"data":        nil,
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:] // 去掉"Bearer "前缀
		} else {
			tokenString = authHeader
		}

		// 解析token
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"operationID": "auth_middleware",
				"code":        401,
				"msg":         "Invalid token: " + err.Error(),
				"data":        nil,
			})
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中供后续使用
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
func GetUserId(c *gin.Context) string {
	return c.GetString("userID")
}
