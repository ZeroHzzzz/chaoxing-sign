package chaoxing

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"chaoxing/internal/pkg/xerr"
	"chaoxing/internal/utils"
	"context"
	"fmt"
	"log"
	"strconv"
)

func (c *Chaoxing) SignLogic(ctx context.Context, act models.ActivityType, signCfg models.SignConfigType, enc, signCode, username string) error {
	status := c.PreSign(ctx, act, username)
	if !status {
		return xerr.PreSignErr
	}

	switch act.OtherID {
	case 0:
		{
			if act.IfPhoto == 1 {
				// todo: 补充拍照签到逻辑
			} else {
				// 普通签到
				status = c.GeneralSign(ctx, act, username)
				if !status {
					return xerr.SignErr
				}
			}
			break
		}
	case 2:
		{
			// 二维码
			// 先获取用户名
			name, err := c.GetUserName(ctx, username)
			if err != nil {
				fmt.Println(err)
			}

			// 暂时先传空值
			status = c.QrcodeSign(ctx, models.LocationType{}, enc, name, act.ActivityID, username)
			if !status {
				return xerr.SignErr
			}
			break
		}
	case 3:
		{
			// 手势签到
			status = c.CodeSign(ctx, act, signCode, username)
			if !status {
				return xerr.SignErr
			}
			break
		}
	case 4:
		{
			// 定位签到
			name, err := c.GetUserName(ctx, username)
			if err != nil {
				fmt.Println(err)
				return err
			}

			var signFlag = false
			for _, location := range signCfg.Locations {
				status = c.LocationSign(ctx, location, name, act.ActivityID, username)
				if status {
					signFlag = true
					break
				}
			}

			if !signFlag {
				return xerr.SignErr
			}
			break
		}
	case 5:
		{
			// 签到码签到
			status = c.CodeSign(ctx, act, signCode, username)
			if !status {
				return xerr.SignErr
			}
			break
		}
	}

	return nil
}
func (c *Chaoxing) PreSign(ctx context.Context, act models.ActivityType, username string) bool {
	cookies := c.Cookie.ToCookies()
	r, err := c.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"courseId":        act.Course.CourseID,
			"classId":         act.Course.ClassID,
			"activePrimaryId": act.ActivityID,
			"general":         "1",
			"sys":             "1",
			"ls":              "1",
			"appType":         "15",
			"uid":             c.Cookie.Uid,
			"ut":              "s",
		}).
		Get(globals.PRESIGN_URL)
	if r.StatusCode() == 302 {
		log.Println("[PreSign] 获取预签到失败，可能是 Cookie 过期")
		return false
	} else if err != nil {
		log.Printf("[PreSign] error: %v\n", err)
		return false
	}

	// ANALYSIS
	r, err = c.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"vs":          "1",
			"DB_STRATEGY": "RANDOM",
			"aid":         act.ActivityID,
		}).
		Get(globals.ANALYSIS_URL)
	if r.StatusCode() == 302 {
		log.Println("[PreSign] 获取analysis失败，可能是 Cookie 过期")
		return false
	} else if err != nil {
		log.Printf("[PreSign] 获取analysis失败: %v\n", err)
		return false
	}

	code := utils.ParseAnalysis(r.String())

	// ANALYSIS2
	r, err = c.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"DB_STRATEGY": "RANDOM",
			"code":        code,
		}).
		Get(globals.ANALYSIS2_URL)
	if r.StatusCode() == 302 {
		log.Println("[PreSign] 获取analysis2失败，可能是 Cookie 过期")
		return false
	} else if err != nil {
		log.Printf("[PreSign] 获取analysis2失败: %v\n", err)
		return false
	}

	log.Println("[PreSign] " + r.String())
	return true
}

func (c *Chaoxing) GeneralSign(ctx context.Context, act models.ActivityType, username string) bool {
	cookies := c.Cookie.ToCookies()
	r, err := c.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"activeId":  act.ActivityID,
			"uid":       c.Cookie.Uid,
			"latitude":  "-1",
			"longitude": "-1",
			"appType":   "15",
			"fid":       c.Cookie.Fid,
			"name":      username,
		}).
		Get(globals.PPT_SIGN_URL)

	if r.StatusCode() == 302 {
		log.Println("[通用] 签到失败，可能是 Cookie 过期")
		return false
	} else if err != nil {
		log.Printf("[通用] 签到失败: %v\n", err)
		return false
	}

	if r.String() == "success" {
		log.Println("[通用] 签到成功")
	} else {
		log.Printf("[通用] 签到失败: %s\n", r.String())
		return false
	}

	return true
}

func (c *Chaoxing) CodeSign(ctx context.Context, act models.ActivityType, signCode, username string) bool {
	cookies := c.Cookie.ToCookies()
	r, err := c.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"activeId":  act.ActivityID,
			"uid":       c.Cookie.Uid,
			"latitude":  "-1",
			"longitude": "-1",
			"appType":   "15",
			"fid":       c.Cookie.Fid,
			"signCode":  signCode,
			"name":      username,
		}).
		Get(globals.PPT_SIGN_URL)

	if r.StatusCode() == 302 {
		log.Println("[通用] 签到失败，可能是 Cookie 过期")
		return false
	} else if err != nil {
		log.Printf("[通用] 签到失败: %v\n", err)
		return false
	}

	if r.String() == "success" {
		log.Println("[通用] 签到成功")
	} else {
		log.Printf("[通用] 签到失败: %s\n", r.String())
		return false
	}

	return true
}

func (c *Chaoxing) QrcodeSign(ctx context.Context, location models.LocationType, enc, name, activeId, username string) bool {
	formated_location := fmt.Sprintf("{\"result\":\"1\",\"address\":\"%s\",\"latitude\":%s,\"longitude\":%s,\"altitude\":%s}", location.Address, location.Latitude, location.Longitude, location.Altitude)
	cookies := c.Cookie.ToCookies()
	r, err := c.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"enc":       enc,
			"activeId":  activeId,
			"uid":       c.Cookie.Uid,
			"location":  formated_location,
			"appType":   "15",
			"fid":       c.Cookie.Fid,
			"name":      name,
			"clientip":  "",
			"latitude":  "-1",
			"longitude": "-1",
		}).
		Get(globals.PPT_SIGN_URL)

	if r.StatusCode() == 302 {
		log.Println("[Qrcode] 签到失败，可能是 Cookie 过期")
		return false
	} else if err != nil {
		log.Printf("[Qrcode] 签到失败: %v\n", err)
		return false
	}

	if r.String() == "success" {
		log.Println("[Qrcode] 签到成功")
	} else {
		log.Printf("[Qrcode] 签到失败: %s\n", r.String())
		return false
	}

	log.Println("[Qrcode] " + r.String())
	return true
}

func (c *Chaoxing) LocationSign(ctx context.Context, location models.LocationType, name, activeId, username string) bool {
	cookies := c.Cookie.ToCookies()
	r, err := c.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"activeId":  activeId,
			"address":   location.Address,
			"uid":       c.Cookie.Uid,
			"appType":   "15",
			"fid":       c.Cookie.Fid,
			"name":      name,
			"clientip":  "",
			"latitude":  location.Latitude,
			"longitude": location.Longitude,
			"ifTiJiao":  "1",
		}).
		Get(globals.PPT_SIGN_URL)
	if r.StatusCode() == 302 {
		log.Println("[Location] 签到失败，可能是 Cookie 过期")
		return false
	} else if err != nil {
		log.Printf("[Location] 签到失败: %v\n", err)
		return false
	}

	if r.String() == "success" {
		log.Println("[Location] 签到成功")
	} else {
		log.Printf("[Location] 签到失败: %s\n", r.String())
		return false
	}

	return true
}

func (c *Chaoxing) GetActivityLogic(ctx context.Context, course models.CourseType, username string) ([]models.ActivityType, error) {
	// 本系统将以活动为主而不是课程
	acts, err := c.GetActivity(ctx, course, username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, act := range acts {
		err = c.GetPPTActivityInfo(ctx, username, &act)
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
func (c *Chaoxing) GetActivity(ctx context.Context, course models.CourseType, username string) ([]models.ActivityType, error) {
	var resp GetActivityChaoxingResp

	formData := map[string]string{
		"fid":      "0",
		"courseId": course.CourseID,
		"classId":  course.ClassID,
		// "_":        strconv.FormatInt(time.Now().Unix(), 10),
	}

	cookies := c.Cookie.ToCookies()
	r, err := c.Rty.R().
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
func (c *Chaoxing) GetPPTActivityInfo(ctx context.Context, username string, activity *models.ActivityType) error {
	var resp GetPPTActivityInfoResp

	cookies := c.Cookie.ToCookies()

	r, err := c.Rty.R().
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
