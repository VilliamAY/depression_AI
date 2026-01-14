package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "操作成功",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalServerError 服务器内部错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, 500, message)
}

// ValidationError 验证错误
func ValidationError(c *gin.Context, message string) {
	Error(c, 422, message)
}

// CustomError 自定义错误
func CustomError(c *gin.Context, code int, message string) {
	Error(c, code, message)
}

// PageResponse 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// SuccessWithPage 分页成功响应
func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, PageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
} 