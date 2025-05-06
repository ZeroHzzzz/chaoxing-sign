package utils

import (
	"chaoxing/internal/models"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func ParseCookies(cookies []*http.Cookie) models.ChaoxingCookieType {
	// Reflection
	var cookie models.ChaoxingCookieType
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

func ParseCourseData(data string) []models.CourseType {
	var courses []models.CourseType
	reader := strings.NewReader(data)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal("加载 HTML 失败：", err)
		return nil
	}

	doc.Find("li.course.clearfix").Each(func(i int, li *goquery.Selection) {
		// 获取 courseId 和 clazzId
		courseId, exists := li.Attr("courseid")
		if !exists {
			log.Println("警告：未找到 courseId 属性")
			return
		}
		classId, exists := li.Attr("clazzid")
		if !exists {
			log.Println("警告：未找到 clazzId 属性")
			return
		}

		// 在 <li> 中查找课程名
		span := li.Find("span.course-name.overHidden2")
		titleAttr, _ := span.Attr("title") // 虽然我们检查了是否存在，但这里假定总是存在的
		// textContent := span.Text()
		courses = append(courses, models.CourseType{
			CourseID: courseId,
			ClassID:  classId,
			Name:     titleAttr,
		})
	})

	return courses
}

func ParseUserName(data string) string {
	endOfMessageName := strings.Index(data, "messageName") + 20
	endOfName := strings.Index(data[endOfMessageName:], `"`) + endOfMessageName
	return data[endOfMessageName:endOfName]
}
