package redis

import "github.com/redis/go-redis/v9"

func GetRdb() *redis.Client {
	info := getConfig()

	return initRdb(info)
}

func initRdb(info RedisInfoConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     info.Host + ":" + info.Port,
		Password: info.Password,
		DB:       info.DB,
	})
	return redisClient
}
