package services

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"context"
	"log"
	"strconv"
)

func GetActivityLogic(ctx context.Context, course models.CourseType, username string) ([]models.ActivityType, error) {
	// 本系统将以活动为主而不是课程
	acts, err := GetActivity(ctx, course, username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, act := range acts {
		err = GetPPTActivityInfo(ctx, username, &act)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return acts, nil
}

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

// 获取课程活动
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

	if r.StatusCode() == 302 {
		log.Println("获取活动列表失败，可能是 Cookie 过期")
		return nil, nil
	}
	// fmt.Println(r.String())
	// fmt.Println(resp)

	var activity []models.ActivityType
	for _, data := range resp.Data.ActiveList {
		otherID, _ := strconv.Atoi(data.OtherID)
		if data.Status == 1 && otherID >= 0 && otherID <= 5 {
			activity = append(activity, models.ActivityType{
				ActivityID: strconv.Itoa(data.ID),
				OtherID:    otherID,
				Name:       data.NameOne,
				Course: models.CourseType{
					CourseID: course.CourseID,
					ClassID:  course.ClassID,
				},
			})
		}
	}

	if len(activity) == 0 {
		log.Println("此课程无活动可查")
		return nil, nil
	}

	return activity, nil
	// if len(resp.Data.ActiveList) != 0 {
	// 	data := resp.Data.ActiveList[0]
	// 	otherID, _ := strconv.Atoi(data.OtherID)
	// 	if data.Status == 1 && otherID >= 0 && otherID <= 5 {
	// 		activity = models.ActivityType{
	// 			ActivityID: strconv.Itoa(data.ID),
	// 			OtherID:    otherID,
	// 			Name:       data.NameOne,
	// 			CourseID:   course.CourseID,
	// 			ClassID:    course.ClassID,
	// 		}
	// 	} else {
	// 		log.Println("活动已结束或不支持")
	// 		return nil, nil
	// 	}
	// } else {
	// 	log.Println("无活动可查")
	// 	return nil, nil
	// }

	// return &activity, nil
}

type GetPPTActivityInfoResp struct {
	ErrorMsg string `json:"errorMsg"`
	Data     struct {
		Ifphoto int `json:"ifphoto"`
		// 这三个参数启用验证码的时候全都是1，可以考虑只保留一个
		OpenPreventCheatFlag int `json:"openPreventCheatFlag"`
		ShowVCode            int `json:"showVCode"`
		IfNeedVCode          int `json:"ifNeedVCode"`
	} `json:"data"`
}

// 获取活动信息（验证码、图片）
func GetPPTActivityInfo(ctx context.Context, username string, activity *models.ActivityType) error {
	var resp GetPPTActivityInfoResp
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return err
	}
	cookies := cookieData.ToCookies()

	r, err := svc.Rty.R().
		SetCookies(cookies).
		SetQueryParam("activeId", activity.ActivityID).
		SetResult(&resp).
		Get(globals.GET_ACTIVITY_INFO_URL)

	if err != nil && r.StatusCode() == 302 {
		log.Println("获取活动信息失败，可能是 Cookie 过期")
		return nil
	}

	activity.OpenPreventCheatFlag = resp.Data.OpenPreventCheatFlag
	activity.IfPhoto = resp.Data.Ifphoto
	return nil
}
