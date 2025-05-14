package utils

import (
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}
