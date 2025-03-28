package chao

import (
	"bytes"
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"dario.cat/mergo"
)

// 使用 AES-CBC 模式加密消息
func encryptByAES(message, key string) (string, error) {
	keyBytes := []byte(key)
	iv := keyBytes // IV 和密钥相同

	// 创建 AES 加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 填充明文到块大小的整数倍
	blockSize := block.BlockSize()
	plainText := pkcs7Padding([]byte(message), blockSize)

	mode := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func parseCookie(cookies []*http.Cookie) models.UserCookieType {
	// Reflection
	var cookie models.UserCookieType
	val := reflect.ValueOf(&cookie).Elem()
	typ := val.Type()

	for i := range val.NumField() {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		for _, cookie := range cookies {
			if cookie.Name == jsonTag {
				fieldVal := val.Field(i)
				if fieldVal.CanSet() {
					switch field.Type.Kind() {
					case reflect.String:
						fieldVal.SetString(cookie.Value)
					case reflect.Bool:
						value := cookie.Value == "true" || cookie.Value == "1"
						fieldVal.SetBool(value)
					}
				}
				break
			}
		}
	}
	return cookie
}

func (c *Chao) LoginByPass(ctx context.Context, username string, password string) error {
	encryptedPassword, err := encryptByAES(password, globals.Secret)
	if err != nil {
		return err
	}

	encryptedUsername, err := encryptByAES(username, globals.Secret)
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
		fmt.Println("登录失败")
		return err
	}

	// 获取 Set-Cookie
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		fmt.Println("网络异常，未获取到 Cookie")
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

	userCookie := parseCookie(cookies)

	err = mergo.Merge(&cookie, userCookie)
	if err != nil {
		fmt.Printf("合并失败: %v\n", err)
		return err
	}

	fmt.Printf("登录成功: %v\n", cookie)
	return nil
}
