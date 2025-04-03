package services

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/utils"
	"context"
	"log"
)

func PreSign(ctx context.Context, activityID, courseID, classID, username string) (string, error) {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return "", err
	}

	cookies := cookieData.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"courseId":        courseID,
			"classId":         classID,
			"activePrimaryId": activityID,
			"general":         "1",
			"sys":             "1",
			"ls":              "1",
			"appType":         "15",
			"uid":             cookieData.Uid,
			"ut":              "s",
		}).
		Get(globals.PRESIGN_URL)
	if r.StatusCode() == 302 {
		log.Println("获取预签到失败，可能是 Cookie 过期")
		return "", nil
	} else if err != nil {
		log.Printf("获取预签到失败: %v\n", err)
		return "", err
	}

	// ANALYSIS
	r, err = svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"vs":          "1",
			"DB_STRATEGY": "RANDOM",
			"aid":         activityID,
		}).
		Get(globals.ANALYSIS_URL)
	if r.StatusCode() == 302 {
		log.Println("获取analysis失败，可能是 Cookie 过期")
		return "", nil
	} else if err != nil {
		log.Printf("获取analysis失败: %v\n", err)
		return "", err
	}

	code := utils.ParseAnalysis(r.String())

	// ANALYSIS2
	r, err = svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"DB_STRATEGY": "RANDOM",
			"code":        code,
		}).
		Get(globals.ANALYSIS2_URL)
	if r.StatusCode() == 302 {
		log.Println("获取analysis2失败，可能是 Cookie 过期")
		return "", nil
	} else if err != nil {
		log.Printf("获取analysis2失败: %v\n", err)
		return "", err
	}

	log.Println("请求结果: " + r.String())
	return "", nil
}
