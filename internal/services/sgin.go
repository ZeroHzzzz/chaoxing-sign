package services

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"chaoxing/internal/utils"
	"context"
	"fmt"
	"log"
)

func SignLogic(ctx context.Context, act models.ActivityType, signCfg models.SignConfigType, enc, username string) error {
	err := PreSign(ctx, act.ActivityID, act.CourseID, act.ClassID, username)
	if err != nil {
		log.Println(err)
		return err
	}

	switch act.OtherID {
	case 0:
		{
			if act.IfPhoto == 1 {
				// todo: 补充拍照签到逻辑
			} else {
				// 普通签到
				err = GeneralSign(ctx, act.ActivityID, act.CourseID, act.ClassID, username)
				if err != nil {
					log.Println(err)
					return err
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
			err = QrcodeSign(ctx, enc, name, act.ActivityID, "", "lat", "lon", "0", username)
			if err != nil {
				log.Println(err)
				return err
			}
			break
		}
	case 3:
		{
			// 手势签到
			err = GeneralSign(ctx, act.ActivityID, act.CourseID, act.ClassID, username)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	case 4:
		{
			// 定位签到
			name, err := GetUserName(ctx, username)
			if err != nil {
				fmt.Println(err)
			}
			err = LocationSign(ctx, signCfg, name, act.ActivityID, username)
			if err != nil {
				log.Println(err)
				return err
			}
			break
		}
	case 5:
		{
			// 签到码签到
			err = GeneralSign(ctx, act.ActivityID, act.CourseID, act.ClassID, username)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}
func PreSign(ctx context.Context, activityID, courseID, classID, username string) error {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return err
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
		return nil
	} else if err != nil {
		log.Printf("获取预签到失败: %v\n", err)
		return err
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
		return nil
	} else if err != nil {
		log.Printf("获取analysis失败: %v\n", err)
		return err
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
		return nil
	} else if err != nil {
		log.Printf("获取analysis2失败: %v\n", err)
		return err
	}

	log.Println("请求结果: " + r.String())
	return nil
}

type GeneralSignResp struct {
	Data string `json:"data"`
}

func GeneralSign(ctx context.Context, activityID, courseID, classID, username string) error {
	var resp GeneralSignResp
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return err
	}

	cookies := cookieData.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"activeId":  activityID,
			"uid":       cookieData.Uid,
			"latitude":  "-1",
			"longitude": "-1",
			"appType":   "15",
			"fid":       cookieData.Fid,
			"name":      username,
		}).
		SetResult(&resp).
		Get(globals.PPT_SIGN_URL)

	if r.StatusCode() == 302 {
		log.Println("签到失败，可能是 Cookie 过期")
		return nil
	} else if err != nil {
		log.Printf("签到失败: %v\n", err)
		return err
	}

	if resp.Data == "success" {
		log.Println("[通用]签到成功")
	} else {
		log.Printf("[通用]签到失败: %s\n", resp.Data)
		return nil
	}

	log.Println("[PreSign]: " + r.String())
	return nil
}

func QrcodeSign(ctx context.Context, enc, name, activeId, address, lat, lon, altitude, username string) error {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return err
	}

	location := fmt.Sprintf("{\"result\":\"1\",\"address\":\"%s\",\"latitude\":%s,\"longitude\":%s,\"altitude\":%s}", address, lat, lon, altitude)
	cookies := cookieData.ToCookies()
	r, err := svc.Rty.R().
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"enc":       enc,
			"activeId":  activeId,
			"uid":       cookieData.Uid,
			"location":  location,
			"appType":   "15",
			"fid":       cookieData.Fid,
			"name":      name,
			"clientip":  "",
			"latitude":  "-1",
			"longitude": "-1",
		}).
		Get(globals.PPT_SIGN_URL)

	if r.StatusCode() == 302 {
		log.Println("签到失败，可能是 Cookie 过期")
		return nil
	} else if err != nil {
		log.Printf("签到失败: %v\n", err)
		return err
	}

	log.Println("[Qrcode]: " + r.String())
	return nil
}

func LocationSign(ctx context.Context, signCfg models.SignConfigType, name, activeId, username string) error {
	cookieData, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return err
	}

	cookies := cookieData.ToCookies()
	for _, location := range signCfg.Locations {
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
			log.Println("签到失败，可能是 Cookie 过期")
			return nil
		} else if err != nil {
			log.Printf("签到失败: %v\n", err)
			return err
		}

		if r.String() != "success" {
			log.Printf("[Location]: 签到失败， %s\n", r.String())
			continue
		}
		log.Println("[Location]: " + r.String())
		break
	}

	return nil
}
