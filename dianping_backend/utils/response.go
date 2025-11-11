package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 统一响应结构
// EN: Unified API response envelope
type Result struct {
	Success bool        `json:"success"`
	ErrorMsg string     `json:"errorMsg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

// SuccessResult 成功响应
// EN: Success response with message
func SuccessResult(message string) *Result {
	return &Result{
		Success: true,
		Data:    message,
	}
}

// SuccessResultWithData 成功响应带数据
// EN: Success response with data payload
func SuccessResultWithData(data interface{}) *Result {
	return &Result{
		Success: true,
		Data:    data,
	}
}

// ErrorResult 错误响应
// EN: Error response with message
func ErrorResult(errorMsg string) *Result {
	return &Result{
		Success:  false,
		ErrorMsg: errorMsg,
	}
}

// Response 统一响应处理
// EN: Write unified response with HTTP 200 for both success/error
func Response(c *gin.Context, result *Result) {
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// SuccessResponse 成功响应
// EN: Shorthand success response
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse 错误响应
// EN: Shorthand error response with HTTP status code
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Result{
		Success:  false,
		ErrorMsg: message,
	})
}

// PageResult 分页响应
// EN: Paginated list response helper
func PageResult(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"list":  data,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}
