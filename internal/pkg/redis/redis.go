package redis

import "github.com/go-redis/redis/v8"

var RedisClient *redis.Client
var RedisInfo RedisInfoConfig

func init() {
	info := getConfig()

	RedisClient = getRedisClient(info)
	RedisInfo = info

}
