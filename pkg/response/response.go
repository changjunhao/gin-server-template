package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 返回400错误响应
func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, message)
}

// Unauthorized 返回401错误响应
func Unauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, message)
}

// Forbidden 返回403错误响应
func Forbidden(c *gin.Context, message string) {
	Fail(c, http.StatusForbidden, message)
}

// NotFound 返回404错误响应
func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, message)
}

// ServerError 返回500错误响应
func ServerError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, message)
}