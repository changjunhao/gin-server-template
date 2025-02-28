package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件，记录请求信息
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 请求方法和路径
		method := c.Request.Method
		path := c.Request.URL.Path

		// 状态码
		statusCode := c.Writer.Status()

		// 客户端IP
		clientIP := c.ClientIP()

		// 记录日志
		log.Printf("%s | %3d | %13v | %15s | %s",
			method,
			statusCode,
			latency,
			clientIP,
			path,
		)
	}
}