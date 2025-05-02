package xerr

import "net/http"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	NotLoginErr = NewError(http.StatusUnauthorized, 200401, "未登录或登录过期")
	PreSignErr  = NewError(http.StatusBadRequest, 200400, "预签到失败")
	SignErr     = NewError(http.StatusBadRequest, 200400, "签到失败")

	HttpErr   = NewError(http.StatusBadRequest, 200500, "请求失败")
	ServerErr = NewError(http.StatusInternalServerError, 200500, "服务器错误")
)

func OtherError(message string) *Error {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}
