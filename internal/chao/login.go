package chao

import (
	"bytes"
	"context"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

var key = []byte("u2oh6Vu^HWe40fj")

type UserCookieType struct {
	Fid   string `json:"fid"`
	Pid   string `json:"pid"`
	Refer string `json:"refer"`
	Blank string `json:"_blank"`
	T     bool   `json:"t"`
	Vc3   string `json:"vc3"`
	Uid   string `json:"_uid"`
	D     string `json:"_d"`
	Uf    string `json:"uf"`
	Lv    string `json:"lv"`
}

var DefaultParams = UserCookieType{
	Fid:   "-1",
	Pid:   "-1",
	Refer: "http%3A%2F%2Fi.chaoxing.com",
	Blank: "1",
	T:     true,
	Vc3:   "",
	Uid:   "",
	D:     "",
	Uf:    "",
	Lv:    "",
}

func encoder(password string, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充密码以满足块大小要求
	padding := block.BlockSize() - len(password)%block.BlockSize()
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	passwordBytes := append([]byte(password), padText...)

	// 加密
	ciphertext := make([]byte, len(passwordBytes))
	mode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	mode.CryptBlocks(ciphertext, passwordBytes)

	// 返回 Base64 编码后的密文
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Chao) LoginByPass(ctx context.Context, username string, password string) error {
	encryptedPassword, err := encoder(password, key)
	if err != nil {
		return err
	}

	formData := map[string]string{
		"uname":            username,
		"password":         encryptedPassword,
		"fid":              "-1",
		"t":                "true",
		"refer":            "https%253A%252F%252Fi.chaoxing.com",
		"forbidotherlogin": "0",
		"validate":         "",
	}

	// 使用 resty 发送请求
	client := c.rty
	resp, err := client.R().
		SetFormData(formData).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Post("LOGIN.URL")

	if err != nil {
		return err
	}

	// 解析 JSON 响应
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	// 检查登录状态
	status, ok := result["status"].(bool)
	if !ok || !status {
		fmt.Println("登录失败")
		return err
	}

	// 获取 Set-Cookie
	cookies := resp.Header().Values("Set-Cookie")
	if len(cookies) == 0 {
		fmt.Println("网络异常，未获取到 Cookie")
		return err
	}

	// 解析 Cookie 并存入 Map
	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		parts := strings.Split(cookie, ";")
		for _, part := range parts {
			kv := strings.Split(strings.TrimSpace(part), "=")
			if len(kv) == 2 {
				cookieMap[kv[0]] = kv[1]
			}
		}
	}

	// 合并默认参数和 Cookie
	loginResult := DefaultParams
	for key, value := range cookieMap {
		switch key {
		case "_uid":
			loginResult.Uid = value
		case "_d":
			loginResult.D = value
		case "uf":
			loginResult.Uf = value
		case "lv":
			loginResult.Lv = value
		case "vc3":
			loginResult.Vc3 = value
		}
	}

	fmt.Println("登录成功")
	return nil
}
