package xerr

import (
	"chaoxing/internal/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Error struct {
	Code  int
	Msg   string
	Level log.Level
}

var (
	NotLoginErr = NewError(200401, log.LevelInfo, "未登录或登录过期")
	PreSignErr  = NewError(200400, log.LevelInfo, "预签到失败")
	SignErr     = NewError(200400, log.LevelInfo, "签到失败")

	HttpErr   = NewError(200500, log.LevelInfo, "请求失败")
	ServerErr = NewError(200500, log.LevelInfo, "服务器错误")

	PoolClosedErr = NewError(200500, log.LevelInfo, "线程池已关闭")
	PoolFullErr   = NewError(200500, log.LevelInfo, "线程池已满")
)

func (e *Error) Error() string {
	return e.Msg
}

func NewError(code int, level log.Level, msg string) *Error {
	return &Error{
		Code:  code,
		Msg:   msg,
		Level: level,
	}
}

// AbortWithException 用于返回自定义错误信息
func AbortWithException(c *gin.Context, apiError *Error, err error) {
	logError(c, apiError, err)
	_ = c.AbortWithError(200, apiError) //nolint:errcheck
}

// logError 记录错误日志
func logError(c *gin.Context, apiErr *Error, err error) {
	// 构建日志字段
	logFields := []zap.Field{
		zap.Int("error_code", apiErr.Code),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("ip", c.ClientIP()),
		zap.Error(err), // 记录原始错误信息
	}
	log.GetLogFunc(apiErr.Level)(apiErr.Msg, logFields...)
}
