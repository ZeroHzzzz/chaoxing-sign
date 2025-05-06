package dao

import (
	"chaoxing/internal/models"
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

func (c *Dao) NewUserSignConfig(ctx context.Context, phone string, config models.SignConfigType) error {
	// 暂时考虑将用户cookie和签到配置分开存储

	configJSON, err := json.Marshal(config)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return err
	}

	err = c.Rdb.Set(ctx, "sign_config:"+phone, configJSON, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Dao) GetUserSignConfig(ctx context.Context, phone string) (*models.SignConfigType, error) {
	val, err := c.Rdb.Get(ctx, "sign_config:"+phone).Result()
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
