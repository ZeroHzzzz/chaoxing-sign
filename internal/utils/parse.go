package utils

import (
	"chaoxing/internal/models"
	"net/http"
	"reflect"
	"strings"
)

func ParseAnalysis(data string) string {
	codeStart := strings.Index(data, "code='+''")
	if codeStart == -1 {
		return ""
	}

	codeStart += 8
	code := data[codeStart:]

	codeEnd := strings.Index(code, "'")
	if codeEnd == -1 {
		return ""
	}

	code = code[:codeEnd]
	return code
}

func ParseCookies(cookies []*http.Cookie) models.UserCookieType {
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

func ParseIMParams(data string) models.IMParamsType {
	var imParams models.IMParamsType

	index := strings.Index(data, "id=\"myName\"")
	if index == -1 {
		return imParams
	}
	start := index + 35
	end := strings.Index(data[start:], "<")
	if end == -1 {
		return imParams
	}
	imParams.MyName = data[start : start+end]

	// 提取 myToken
	index = strings.Index(data, "id=\"myToken\"")
	if index == -1 {
		return imParams
	}
	start = index + 36
	end = strings.Index(data[start:], "<")
	if end == -1 {
		return imParams
	}
	imParams.MyToken = data[start : start+end]

	// 提取 myTuid
	index = strings.Index(data, "id=\"myTuid\"")
	if index == -1 {
		return imParams
	}
	start = index + 35
	end = strings.Index(data[start:], "<")
	if end == -1 {
		return imParams
	}
	imParams.MyTuid = data[start : start+end]

	return imParams
}

func ParseCourse(data string) []models.CourseType {
	var courses []models.CourseType
	i := 0 // 全局索引

	for {
		courseIndex := strings.Index(data[i:], "course_")
		if courseIndex == -1 {
			break
		}

		i += courseIndex
		endOfCourseId := strings.Index(data[i+len("course_"):], "_")
		if endOfCourseId == -1 {
			break
		}

		endOfCourseId += i + len("course_")
		courseId := data[i+len("course_") : endOfCourseId]

		classIdStart := endOfCourseId + 1
		classIdEnd := strings.Index(data[classIdStart:], `"`)
		if classIdEnd == -1 {
			break
		}

		classIdEnd += classIdStart
		classId := data[classIdStart:classIdEnd]

		courses = append(courses, models.CourseType{
			CourseID: courseId,
			ClassID:  classId,
		})

		i = classIdEnd + 1
	}

	return courses
}
