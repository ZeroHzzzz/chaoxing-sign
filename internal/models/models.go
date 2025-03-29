package models

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type UserCookieType struct {
	Name  string `json:"name"`
	UID   string `json:"uid"`
	Fid   string `json:"fid"`
	Pid   string `json:"pid"`
	Refer string `json:"refer"`
	Blank string `json:"_blank"`
	T     bool   `json:"t"`
	Vc3   string `json:"vc3"`
	Uid   string `json:"_uid"`
	D     string `json:"_d"`
	Uf    string `json:"uf"`
	Lv    string `json:"lv"`
}

func (uc *UserCookieType) ToCookies() []*http.Cookie {
	var cookies []*http.Cookie
	val := reflect.ValueOf(uc).Elem()
	typ := val.Type()

	for i := range val.NumField() {
		field := typ.Field(i)
		value := val.Field(i)

		// 获取 json 标签作为 Cookie 名
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		cookieName := strings.Split(tag, ",")[0] // 处理可能的多标签，如 `json:"name,omitempty"`

		// 根据字段类型获取值（转为字符串）
		var cookieValue string
		switch value.Kind() {
		case reflect.String:
			cookieValue = value.String()
		case reflect.Bool:
			cookieValue = strconv.FormatBool(value.Bool())
		default:
			continue
		}

		if cookieValue == "" && value.Kind() != reflect.Bool {
			continue
		}

		cookies = append(cookies, &http.Cookie{
			Name:  cookieName,
			Value: cookieValue,
			// 可选：设置 Domain、Path 等（根据实际需求）
			// Domain: "example.com",
			// Path:   "/",
		})
	}
	return cookies
}
