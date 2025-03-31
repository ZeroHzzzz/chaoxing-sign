package services

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"context"
	"fmt"
	"log"
)

type GetActivityChaoxingResp struct {
	Result int `json:"result"`
	Data   struct {
		ActiveList []struct {
			Status     int    `json:"status"`
			NameHour   string `json:"nameHour"`
			ID         int    `json:"id"`
			OtherID    string `json:"otherId"`
			NameOne    string `json:"nameOne"`
			ActiveType int    `json:"activeType"`
		} `json:"activeList"`
	} `json:"data"`
	ErrorMsg string `json:"errorMsg"`
}

func GetActivity(ctx context.Context, course models.CourseType, username string) ([]models.ActivityType, error) {
	var resp GetActivityChaoxingResp
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return nil, err
	}
	formData := map[string]string{
		"fid":      "0",
		"courseId": course.CourseID,
		"classId":  course.ClassID,
		// "_":        strconv.FormatInt(time.Now().Unix(), 10),
	}

	cookies := cookieData.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(formData).
		SetResult(&resp).
		Get(globals.GET_ACTIVITY_URL)

	if err != nil {
		log.Printf("获取活动列表失败: %v\n", err)
		return nil, err
	}

	fmt.Println(r.String())
	fmt.Println(resp)
	return nil, nil
}
