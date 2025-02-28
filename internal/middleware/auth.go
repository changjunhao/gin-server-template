package middleware

import (
	"gin-server/internal/config"
	"gin-server/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authorization, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 解析JWT令牌
		tokenString := parts[1]
		
		// 从应用配置获取JWT密钥
		cfg, err := config.LoadConfig("configs/config.yaml")
		if err != nil {
			response.ServerError(c, "服务器配置错误")
			c.Abort()
			return
		}

		// 解析和验证令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "无效的认证令牌")
			c.Abort()
			return
		}

		// 从令牌中提取声明
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c, "无效的令牌声明")
			c.Abort()
			return
		}

		// 将用户ID存储在上下文中
		userID, ok := claims["user_id"].(float64)
		if !ok {
			response.Unauthorized(c, "无效的用户信息")
			c.Abort()
			return
		}

		// 将用户ID设置到上下文中，供后续处理器使用
		c.Set("userID", uint(userID))
		c.Next()
	}
}