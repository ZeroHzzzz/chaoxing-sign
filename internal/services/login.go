package services

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"chaoxing/internal/utils"
	"context"
	"encoding/json"
	"log"

	"dario.cat/mergo"
)

func LoginByPass(ctx context.Context, username string, password string) (models.UserCookieType, error) {
	encryptedPassword, err := utils.EncryptByAES(password, globals.Secret)
	if err != nil {
		return models.UserCookieType{}, err
	}

	encryptedUsername, err := utils.EncryptByAES(username, globals.Secret)
	if err != nil {
		return models.UserCookieType{}, err
	}

	formData := map[string]string{
		"uname":            encryptedUsername,
		"password":         encryptedPassword,
		"fid":              "-1",
		"t":                "true",
		"refer":            "https://i.chaoxing.com",
		"forbidotherlogin": "0",
		"validate":         "",
	}

	resp, err := svc.Rty.R().
		SetFormData(formData).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Post(globals.LOGIN_URL)

	if err != nil {
		return models.UserCookieType{}, err
	}

	// 解析 JSON 响应
	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return models.UserCookieType{}, err
	}

	// 检查登录状态
	status, ok := result["status"].(bool)
	if !ok || !status {
		log.Println("登陆失败")
		return models.UserCookieType{}, err
	}

	// 获取 Set-Cookie
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		log.Println("网络异常，未获取到 Cookie")
		return models.UserCookieType{}, err
	}

	cookie := models.UserCookieType{
		Fid:   "-1",
		Pid:   "-1",
		Refer: "https://i.chaoxing.com",
		Blank: "1",
		T:     true,
		Vc3:   "",
		Uid:   "",
		D:     "",
		Uf:    "",
		Lv:    "",
	}

	userCookie := utils.Cookie2Struct(cookies)

	err = mergo.Merge(&cookie, userCookie)
	if err != nil {
		log.Printf("合并失败: %v\n", err)
		return models.UserCookieType{}, err
	}

	err = StoreCookies(ctx, username, cookie)
	if err != nil {
		log.Printf("存储 Cookie 失败: %v\n", err)
		return models.UserCookieType{}, err
	}
	// log.Printf("登录成功: %v\n", cookie)
	return cookie, nil
}

func GetPanToken(ctx context.Context, username string) (string, error) {
	cookie, err := GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return "", err
	}

	cookies := cookie.ToCookies()
	resp, err := svc.Rty.R().
		SetCookies(cookies).
		Get(globals.GET_PANTOKEN_URL)

	if err != nil {
		log.Printf("获取网盘 Token 失败: %v\n", err)
		return "", err
	}

	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
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
