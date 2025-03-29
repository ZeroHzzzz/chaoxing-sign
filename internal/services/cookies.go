package services

import (
	"chaoxing/internal/models"
	"context"
	"encoding/json"
	"log"
)

func StoreCookies(ctx context.Context, key string, cookie models.UserCookieType) error {
	cookieJSON, err := json.Marshal(cookie)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return err
	}

	err = svc.Rdb.Set(ctx, "cookie:"+key, cookieJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetCookies(ctx context.Context, key string) (models.UserCookieType, error) {
	cookieJSON, err := svc.Rdb.Get(ctx, "cookie:"+key).Result()
	if err != nil {
		log.Printf("获取 Cookie 失败: %v\n", err)
		return models.UserCookieType{}, err
	}

	var cookie models.UserCookieType
	err = json.Unmarshal([]byte(cookieJSON), &cookie)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return models.UserCookieType{}, err
	}

	return cookie, nil
}
