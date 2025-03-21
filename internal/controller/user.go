package controller

import (
	"gin-server-template/internal/config"
	"gin-server-template/internal/entity"
	"gin-server-template/internal/service"
	"gin-server-template/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Avatar   string `json:"avatar"`
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "无效的请求参数")
		return
	}

	// 创建用户实体
	user := &entity.User{
		Username: req.Username,
		Password: req.Password, // 实际应用中应该对密码进行哈希处理
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1,
	}

	// 调用服务层注册用户
	if err := c.userService.Register(user); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "注册失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{"user_id": user.ID})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "无效的请求参数")
		return
	}

	// 验证用户凭证
	user, err := c.userService.VerifyCredentials(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(ctx, "用户名或密码错误")
		return
	}

	// 生成JWT令牌
	token, err := c.generateToken(user)
	if err != nil {
		response.ServerError(ctx, "生成令牌失败")
		return
	}

	response.Success(ctx, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"email":    user.Email,
		},
	})
}

// GetProfile 获取用户个人资料
func (c *UserController) GetProfile(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "未认证的请求")
		return
	}

	// 获取用户信息
	user, err := c.userService.GetUserByID(userID.(uint))
	if err != nil {
		response.NotFound(ctx, "用户不存在")
		return
	}

	response.Success(ctx, user)
}

// UpdateProfile 更新用户个人资料
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "未认证的请求")
		return
	}

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "无效的请求参数")
		return
	}

	// 获取用户信息
	user, err := c.userService.GetUserByID(userID.(uint))
	if err != nil {
		response.NotFound(ctx, "用户不存在")
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// 保存更新
	if err := c.userService.UpdateUser(user); err != nil {
		response.ServerError(ctx, "更新失败: "+err.Error())
		return
	}

	response.Success(ctx, user)
}

// generateToken 生成JWT令牌
func (c *UserController) generateToken(user *entity.User) (string, error) {
	// 从配置中获取JWT密钥
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		return "", err
	}

	// 创建JWT声明
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(cfg.JWT.Expire)).Unix(), // 从配置中获取过期时间
		"iat":      time.Now().Unix(),
		"iss":      cfg.JWT.Issuer, // 从配置中获取发行者
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret)) // 从配置中获取密钥
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
