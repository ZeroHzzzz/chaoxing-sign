package chaoxing

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"chaoxing/internal/pkg/xerr"
	"chaoxing/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"dario.cat/mergo"
	"github.com/redis/go-redis/v9"
)

func (c *Chaoxing) LoginByPass(ctx context.Context, username string, password string) (models.UserCookieType, error) {
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

	r, err := c.Rty.R().
		SetFormData(formData).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Post(globals.LOGIN_URL)

	if err != nil {
		return models.UserCookieType{}, err
	}

	// 解析 JSON 响应
	var result map[string]any
	err = json.Unmarshal(r.Body(), &result)
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
	cookies := r.Cookies()
	if len(cookies) == 0 {
		log.Println("网络异常，未获取到 Cookie")
		return models.UserCookieType{}, err
	}

	cookie := models.UserCookieType{
		Fid: "-1",
		// Pid:   "-1",
		// Refer: "https://i.chaoxing.com",
		// Blank: "1",
		// T:     true,
		Vc3: "",
		Uid: "",
		D:   "",
		Uf:  "",
		Lv:  "",
	}

	userCookie := utils.ParseCookies(cookies)

	err = mergo.Merge(&cookie, userCookie)
	if err != nil {
		log.Printf("合并失败: %v\n", err)
		return models.UserCookieType{}, err
	}

	err = c.StoreCookies(ctx, username, cookie)
	if err != nil {
		log.Printf("存储 Cookie 失败: %v\n", err)
		return models.UserCookieType{}, err
	}
	// log.Printf("登录成功: %v\n", cookie)
	return cookie, nil
}

func (c *Chaoxing) StoreCookies(ctx context.Context, key string, cookie models.UserCookieType) error {
	cookieJSON, err := json.Marshal(cookie)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return err
	}

	err = c.Rdb.Set(ctx, "cookie:"+key, cookieJSON, 2*time.Hour).Err() // 两小时过期
	if err != nil {
		return err
	}

	return nil
}

func (c *Chaoxing) GetCookies(ctx context.Context, key string) (*models.UserCookieType, error) {
	val, err := c.Rdb.Get(ctx, "cookie:"+key).Result()
	if err != nil {
		if err == redis.Nil {
			// todo: 处理 Cookie 过期的情况
			return nil, err
		}
		log.Printf("获取 Cookie 失败: %v\n", err)
		return nil, err
	}

	var cookie models.UserCookieType
	err = json.Unmarshal([]byte(val), &cookie)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return nil, err
	}

	return &cookie, nil
}

func (c *Chaoxing) GetPanToken(ctx context.Context, username string) (string, error) {
	cookie, err := c.GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return "", err
	}

	cookies := cookie.ToCookies()
	r, err := c.Rty.R().
		SetCookies(cookies).
		Get(globals.GET_PANTOKEN_URL)

	if err != nil {
		log.Printf("获取网盘 Token 失败: %v\n", err)
		return "", err
	}

	var result map[string]any
	err = json.Unmarshal(r.Body(), &result)
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

func (c *Chaoxing) GetCourses(ctx context.Context, username string) ([]models.CourseType, error) {
	cookieData, err := c.GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return nil, err
	}

	formData := map[string]string{
		"courseType":       "1",
		"courseFolderId":   "0",
		"courseFolderSize": "0",
	}

	r, err := c.Rty.R().
		SetHeaders(map[string]string{
			"Accept":          "text/html, */*; q=0.01",
			"Accept-Encoding": "gzip, deflate",
			"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
			"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
			"Cookie":          fmt.Sprintf("_uid=%s; _d=%s; vc3=%s", cookieData.Uid, cookieData.D, cookieData.Vc3),
		}). // 这里cookie格式特殊，因此使用了SetHeaders直接拼接
		SetFormData(formData).
		Post(globals.GET_COURSELIST_URL)

	if r.StatusCode() == 302 {
		log.Println("获取课程列表失败，可能是 Cookie 过期")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	data := r.String()

	courses := utils.ParseCourseData(data)

	if len(courses) == 0 {
		log.Println("无课程可查")
		return nil, nil
	}

	return courses, nil
}

func (c *Chaoxing) GetUserName(ctx context.Context, username string) (string, error) {
	cookieData, err := c.GetCookies(ctx, username)
	if err != nil {
		log.Println("获取 Cookie 失败: ", err)
		return "", xerr.NotLoginErr
	}

	cookies := cookieData.ToCookies()
	r, err := c.Rty.R().
		SetCookies(cookies).
		Get(globals.GET_USER_INFO_URL)
	if r.StatusCode() == 302 {
		log.Println("获取用户信息失败，可能是 Cookie 过期")
		return "", xerr.NotLoginErr
	} else if err != nil {
		log.Println("获取用户信息失败: ", err)
		return "", err
	}

	data := r.String()
	name := utils.ParseUserName(data)
	return name, nil
}

// 获取IM参数（登录用）
func (c *Chaoxing) GetIMParams(ctx context.Context, username string) (*models.IMParamsType, error) {
	cookieData, err := c.GetCookies(ctx, username)
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return nil, err
	}

	cookies := cookieData.ToCookies()
	r, err := c.Rty.R().
		SetCookies(cookies).
		Get(globals.GET_WEBIM_URL)

	if r.StatusCode() == 302 {
		log.Println("获取IM参数失败，可能是 Cookie 过期")
		return nil, xerr.NotLoginErr
	} else if err != nil {
		log.Printf("获取IM参数失败: %v\n", err)
		return nil, err
	}

	data := r.String()

	imParams := utils.ParseIMParams(data)
	// Puid为uid
	imParams.MyPuid = cookieData.Uid
	return &imParams, nil
}

func (c *Chaoxing) StoreSignConfig(ctx context.Context, username string, config models.SignConfigType) error {
	// 暂时考虑将用户cookie和签到配置分开存储

	configJSON, err := json.Marshal(config)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return err
	}

	err = c.Rdb.Set(ctx, "sign_config:"+username, configJSON, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Chaoxing) GetSignConfig(ctx context.Context, username string) (*models.SignConfigType, error) {
	val, err := c.Rdb.Get(ctx, "sign_config:"+username).Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("签到配置不存在")
			return nil, nil
		}
		log.Printf("获取签到配置失败: %v\n", err)
		return nil, err
	}

	var config models.SignConfigType
	err = json.Unmarshal([]byte(val), &config)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return nil, err
	}

	return &config, nil
}
