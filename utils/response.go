package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess       = 0
	CodeError         = 10000
	CodeInvalidParams = 10001
	CodeAuthFailed    = 10002
	CodeNotFound      = 10003
	CodeServerError   = 10004
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "Success",
		Data:    data,
	})
}

func Error(c *gin.Context, httpCode, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func SuccessWithData(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    nil,
	})
}

func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    nil,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeInvalidParams, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, CodeAuthFailed, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeServerError, message)
}

func HandleAppError(c *gin.Context, err error) {
	if err == nil {
		Success(c, nil)
		return
	}

	if IsAppError(err) {
		appErr := GetAppError(err)
		c.JSON(appErr.Code, Response{
			Code:    appErr.Code,
			Message: appErr.Message,
			Data:    nil,
		})
		return
	}

	InternalServerError(c, "服务器内部错误")
}

func HandleAppErrorWithData(c *gin.Context, err error, data interface{}) {
	if err == nil {
		Success(c, data)
		return
	}

	if IsAppError(err) {
		appErr := GetAppError(err)
		c.JSON(appErr.Code, Response{
			Code:    appErr.Code,
			Message: appErr.Message,
			Data:    data,
		})
		return
	}

	InternalServerError(c, "服务器内部错误")
}
