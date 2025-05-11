package middleware

import (
	"chaoxing/internal/pkg/utils"
	"chaoxing/internal/pkg/xerr"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 向下执行请求
		c.Next()

		// 如果存在错误，则处理错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if err != nil {
				var apiErr *xerr.Error

				// 尝试将错误转换为 apiException
				ok := errors.As(err, &apiErr)

				// 如果转换失败，则使用 ServerError
				if !ok {
					apiErr = xerr.ServerErr
					zap.L().Error("Unknown Error Occurred", zap.Error(err))
				}

				utils.JsonErrorResponse(c, apiErr.Code, apiErr.Msg)
				return
			}
		}
	}
}

// HandleNotFound 处理 404 错误。
func HandleNotFound(c *gin.Context) {
	err := xerr.NotFound
	// 记录 404 错误日志
	zap.L().Warn("404 Not Found",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
	)
	utils.JsonResponse(c, http.StatusNotFound, err.Code, err.Msg, nil)
}
