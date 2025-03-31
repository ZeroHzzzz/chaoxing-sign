package services

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"chaoxing/internal/utils"
	"context"
	"log"
)

func GetCourses(ctx context.Context, username string) ([]models.Course, error) {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return nil, err
	}

	cookies := cookieData.ToCookies()

	formData := map[string]string{
		"courseType":       "1",
		"courseFolderId":   "0",
		"courseFolderSize": "0",
	}

	resp, err := svc.Rty.R().
		SetHeaders(map[string]string{
			"Accept":          "text/html, */*; q=0.01",
			"Accept-Encoding": "gzip, deflate",
			"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
			"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
		}).
		SetCookies(cookies).
		SetFormData(formData).
		Post(globals.GET_COURSELIST_URL)

	if resp.StatusCode() == 302 {
		log.Println("获取课程列表失败，可能是 Cookie 过期")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	data := resp.String()

	courses := utils.ParseCourse(data)

	if len(courses) == 0 {
		log.Println("无课程可查")
		return nil, nil
	}

	return courses, nil
}
