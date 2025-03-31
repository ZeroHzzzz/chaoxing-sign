package services

import (
	"chaoxing/internal/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func StoreCookies(ctx context.Context, key string, cookie models.UserCookieType) error {
	cookieJSON, err := json.Marshal(cookie)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return err
	}

	err = svc.Rdb.Set(ctx, "cookie:"+key, cookieJSON, 2*time.Hour).Err() // 两小时过期
	if err != nil {
		return err
	}

	return nil
}

func GetCookies(ctx context.Context, key string) (*models.UserCookieType, error) {
	val, err := svc.Rdb.Get(ctx, "cookie:"+key).Result()
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
