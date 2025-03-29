package utils

import (
	"chaoxing/internal/models"
	"net/http"
	"reflect"
)

func Cookie2Struct(cookies []*http.Cookie) models.UserCookieType {
	// Reflection
	var cookie models.UserCookieType
	val := reflect.ValueOf(&cookie).Elem()
	typ := val.Type()

	for i := range val.NumField() {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		for _, cookie := range cookies {
			if cookie.Name == jsonTag {
				fieldVal := val.Field(i)
				if fieldVal.CanSet() {
					switch field.Type.Kind() {
					case reflect.String:
						fieldVal.SetString(cookie.Value)
					case reflect.Bool:
						value := cookie.Value == "true" || cookie.Value == "1"
						fieldVal.SetBool(value)
					}
				}
				break
			}
		}
	}
	return cookie
}
