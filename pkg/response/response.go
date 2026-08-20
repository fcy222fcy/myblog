package response

import (
	"net/http"

	bizErr "blog/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    bizErr.CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    bizErr.CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// BadRequest 请求错误响应
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    bizErr.CodeBadRequest,
		Message: message,
		Data:    nil,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    bizErr.CodeUnauthorized,
		Message: message,
		Data:    nil,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    bizErr.CodeForbidden,
		Message: message,
		Data:    nil,
	})
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    bizErr.CodeNotFound,
		Message: message,
		Data:    nil,
	})
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    bizErr.CodeInternalServer,
		Message: message,
		Data:    nil,
	})
}

// BizError 业务错误响应
func BizError(c *gin.Context, err error) {
	if bizError, ok := err.(*bizErr.BizError); ok {
		httpStatus := http.StatusOK
		// 根据错误码范围返回合适的 HTTP 状态码
		switch {
		case bizError.Code >= 1000 && bizError.Code < 2000:
			// 用户/认证业务错误 (1000-1999) -> 401
			httpStatus = http.StatusUnauthorized
		case bizError.Code >= 2000 && bizError.Code < 3000:
			// 参数错误 (2000-2999) -> 400
			httpStatus = http.StatusBadRequest
		case bizError.Code >= 3000 && bizError.Code < 4000:
			// 资源错误 (3000-3999) -> 404
			httpStatus = http.StatusNotFound
		case bizError.Code >= 4000 && bizError.Code < 5000:
			// 业务错误 (4000-4999) -> 400
			httpStatus = http.StatusBadRequest
		}
		c.JSON(httpStatus, Response{
			Code:    bizError.Code,
			Message: bizError.Message,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Code:    bizErr.CodeInternalServer,
		Message: err.Error(),
		Data:    nil,
	})
}

// ListData 列表数据结构
type ListData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// SuccessList 成功列表响应
func SuccessList(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, ListData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
