package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/WeChatPadPro/WeChatPadPro/pkg/logger"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Request.URL.Path

		logger.Info("Request: %s %s from %s", c.Request.Method, start, c.ClientIP())

		c.Next()

		status := c.Writer.Status()
		if status >= 400 {
			logger.Error("Response: %s %s - Status: %d", c.Request.Method, start, status)
		}
	}
}

// Auth 认证中间件
func Auth(adminKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要认证的路径
		skipPaths := []string{
			"/health",
			"/ping",
			"/login",
			"/static",
			"/",
			"/sse",
		}

		path := c.Request.URL.Path
		for _, skip := range skipPaths {
			if strings.HasPrefix(path, skip) {
				c.Next()
				return
			}
		}

		// 检查认证
		key := c.Query("key")
		if key == "" {
			key = c.GetHeader("X-Auth-Key")
		}

		// 对于Webhook配置等操作，需要admin key
		if strings.Contains(path, "/webhook") && (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE") {
			if key != adminKey {
				logger.Warn("Unauthorized webhook operation attempt from %s", c.ClientIP())
				c.JSON(401, gin.H{"code": -1, "message": "unauthorized"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: %v", err)
				c.JSON(500, gin.H{
					"code":    -1,
					"message": "internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}