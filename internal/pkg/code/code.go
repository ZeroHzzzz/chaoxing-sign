package code

import "net/http"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	ServerErr = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	HttpErr   = NewError(http.StatusInternalServerError, 200500, "网络异常，请稍后重试!")

	ParamErr    = NewError(http.StatusBadRequest, 200400, "参数错误")
	ResponseErr = NewError(http.StatusBadRequest, 200400, "响应错误")

	NotLogin     = NewError(http.StatusUnauthorized, 200401, "未登录")
	NoPermission = NewError(http.StatusForbidden, 200403, "无权限")

	NotInit  = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	NotFound = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unknown  = NewError(http.StatusInternalServerError, 300500, "系统异常，请稍后重试!")
)

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}
