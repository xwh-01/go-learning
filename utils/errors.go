package utils

import (
	"errors"
	"fmt"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrInvalidToken      = NewAppError(401, "无效的Token", nil)
	ErrMissingAuthHeader = NewAppError(401, "缺少授权头部", nil)
	ErrUserNotFound      = NewAppError(404, "用户不存在", nil)
	ErrInvalidPassword   = NewAppError(401, "密码错误", nil)
	ErrUserExists        = NewAppError(409, "用户已存在", nil)
	ErrArticleNotFound   = NewAppError(404, "文章不存在", nil)
	ErrInvalidParams     = NewAppError(400, "参数错误错误", nil)
	ErrRecordNotFound    = NewAppError(404, "记录未找到", nil)
	ErrInternalServer    = NewAppError(500, "服务器内部错误", nil)
	ErrDatabase          = NewAppError(500, "数据库操作失败", nil)
	ErrRedis             = NewAppError(500, "缓存操作失败", nil)
	ErrUnauthorized      = NewAppError(401, "未授权访问", nil)
	ErrForbidden         = NewAppError(403, "禁止访问", nil)
)

func FormatDBError(err error) *AppError {
	if err == nil {
		return nil
	}
	return NewAppError(ErrDatabase.Code, ErrDatabase.Message, err)
}

func FormatRedisError(err error) *AppError {
	if err == nil {
		return nil
	}
	return NewAppError(ErrRedis.Code, ErrRedis.Message, err)
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return NewAppError(500, "未知错误", err)
}
