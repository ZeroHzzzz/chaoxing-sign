package chao

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"chaoxing/internal/utils"
	"context"
	"encoding/json"
	"log"

	"dario.cat/mergo"
)

func (c *Chao) LoginByPass(ctx context.Context, username string, password string) error {
	encryptedPassword, err := utils.EncryptByAES(password, globals.Secret)
	if err != nil {
		return err
	}

	encryptedUsername, err := utils.EncryptByAES(username, globals.Secret)
	if err != nil {
		return err
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

	client := c.rtyClient
	resp, err := client.R().
		SetFormData(formData).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Post(globals.LOGIN_URL)

	if err != nil {
		return err
	}

	// 解析 JSON 响应
	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	// 检查登录状态
	status, ok := result["status"].(bool)
	if !ok || !status {
		log.Println("登陆失败")
		return err
	}

	// 获取 Set-Cookie
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		log.Println("网络异常，未获取到 Cookie")
		return err
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

	userCookie := utils.ParseCookie(cookies)

	err = mergo.Merge(&cookie, userCookie)
	if err != nil {
		log.Printf("合并失败: %v\n", err)
		return err
	}

	// log.Printf("登录成功: %v\n", cookie)
	return nil
}
