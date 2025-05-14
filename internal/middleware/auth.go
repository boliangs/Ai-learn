package middleware

import (
	"net/http"
	"strings"

	"ai-interview/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "未提供认证token")
			ctx.Abort()
			return
		}

		// 检查token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "token格式错误")
			ctx.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "无效的token: "+err.Error())
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先进行普通认证
		Auth()(c)
		if c.IsAborted() {
			return
		}

		// 检查是否是管理员
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
