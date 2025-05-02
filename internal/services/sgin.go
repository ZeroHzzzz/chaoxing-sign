package services

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"chaoxing/internal/pkg/xerr"
	"chaoxing/internal/utils"
	"context"
	"fmt"
	"log"
)

func SignLogic(ctx context.Context, act models.ActivityType, signCfg models.SignConfigType, enc, username string) error {
	status := PreSign(ctx, act, username)
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
				status = GeneralSign(ctx, act, username)
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
			name, err := GetUserName(ctx, username)
			if err != nil {
				fmt.Println(err)
			}

			// 暂时先传空值
			status = QrcodeSign(ctx, models.LocationType{}, enc, name, act.ActivityID, username)
			if !status {
				return xerr.SignErr
			}
			break
		}
	case 3:
		{
			// 手势签到
			// Todo：这里有些问题，需要后续修改
			status = GeneralSign(ctx, act, username)
			if !status {
				return xerr.SignErr
			}
			break
		}
	case 4:
		{
			// 定位签到
			name, err := GetUserName(ctx, username)
			if err != nil {
				fmt.Println(err)
				return err
			}

			var signFlag = false
			for _, location := range signCfg.Locations {
				status = LocationSign(ctx, location, name, act.ActivityID, username)
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
			// Todo：这里有些问题，需要后续修改
			status = GeneralSign(ctx, act, username)
			if !status {
				return xerr.SignErr
			}
			break
		}
	}

	return nil
}
func PreSign(ctx context.Context, act models.ActivityType, username string) bool {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("[PreSign] 获取 Cookie 失败: %v\n", err)
		return false
	}

	cookies := cookieData.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"courseId":        act.Course.CourseID,
			"classId":         act.Course.ClassID,
			"activePrimaryId": act.ActivityID,
			"general":         "1",
			"sys":             "1",
			"ls":              "1",
			"appType":         "15",
			"uid":             cookieData.Uid,
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
	r, err = svc.Rty.R().
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
	r, err = svc.Rty.R().
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

func GeneralSign(ctx context.Context, act models.ActivityType, username string) bool {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("[通用] 获取 Cookie 失败: %v\n", err)
		return false
	}

	cookies := cookieData.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"activeId":  act.ActivityID,
			"uid":       cookieData.Uid,
			"latitude":  "-1",
			"longitude": "-1",
			"appType":   "15",
			"fid":       cookieData.Fid,
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

func QrcodeSign(ctx context.Context, location models.LocationType, enc, name, activeId, username string) bool {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("[Qrcode] 获取 Cookie 失败: %v\n", err)
		return false
	}

	formated_location := fmt.Sprintf("{\"result\":\"1\",\"address\":\"%s\",\"latitude\":%s,\"longitude\":%s,\"altitude\":%s}", location.Address, location.Latitude, location.Longitude, location.Altitude)
	cookies := cookieData.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"enc":       enc,
			"activeId":  activeId,
			"uid":       cookieData.Uid,
			"location":  formated_location,
			"appType":   "15",
			"fid":       cookieData.Fid,
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

func LocationSign(ctx context.Context, location models.LocationType, name, activeId, username string) bool {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("[Location] 获取 Cookie 失败: %v\n", err)
		return false
	}

	cookies := cookieData.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"activeId":  activeId,
			"address":   location.Address,
			"uid":       cookieData.Uid,
			"appType":   "15",
			"fid":       cookieData.Fid,
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
