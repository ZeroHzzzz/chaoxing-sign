package services

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"chaoxing/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

func GetPanToken(ctx context.Context, username string) (string, error) {
	cookie, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return "", err
	}

	cookies := cookie.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		Get(globals.GET_PANTOKEN_URL)

	if err != nil {
		log.Printf("获取网盘 Token 失败: %v\n", err)
		return "", err
	}

	var result map[string]any
	err = json.Unmarshal(r.Body(), &result)
	if err != nil {
		log.Printf("解析网盘 Token 响应失败: %v\n", err)
		return "", err
	}

	token, ok := result["_token"].(string)
	if !ok {
		log.Println("获取网盘 Token 失败")
		return "", err
	}

	return token, nil
}

func GetCourses(ctx context.Context, username string) ([]models.CourseType, error) {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return nil, err
	}

	formData := map[string]string{
		"courseType":       "1",
		"courseFolderId":   "0",
		"courseFolderSize": "0",
	}

	r, err := svc.Rty.R().
		SetHeaders(map[string]string{
			"Accept":          "text/html, */*; q=0.01",
			"Accept-Encoding": "gzip, deflate",
			"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
			"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
			"Cookie":          fmt.Sprintf("_uid=%s; _d=%s; vc3=%s", cookieData.Uid, cookieData.D, cookieData.Vc3),
		}). // 这里cookie格式特殊，因此使用了SetHeaders直接拼接
		SetFormData(formData).
		Post(globals.GET_COURSELIST_URL)

	if r.StatusCode() == 302 {
		log.Println("获取课程列表失败，可能是 Cookie 过期")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	data := r.String()

	courses := utils.ParseCourse(data)

	if len(courses) == 0 {
		log.Println("无课程可查")
		return nil, nil
	}

	return courses, nil
}

// 获取IM参数（登录用）
func GetIMParams(ctx context.Context, username string) (*models.IMParamsType, error) {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return nil, err
	}

	cookies := cookieData.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		Get(globals.GET_WEBIM_URL)

	if r.StatusCode() == 302 {
		log.Println("获取IM参数失败，可能是 Cookie 过期")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	data := r.String()

	imParams := utils.ParseIMParams(data)
	// Puid为uid
	imParams.MyPuid = cookieData.Uid
	return &imParams, nil
}
