package redis

import (
	"chaoxing/internal/globals"

	"github.com/go-redis/redis/v8"
)

type RedisInfoConfig struct {
	Host     string
	Port     string
	DB       int
	Password string
}

func getRedisClient(info RedisInfoConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     info.Host + ":" + info.Port,
		Password: info.Password,
		DB:       info.DB,
	})
	return redisClient
}

// func defaultRedisConfig() RedisInfoConfig {
// 	return RedisInfoConfig{
// 		Host:     "localhost",
// 		Port:     "6379",
// 		DB:       0,
// 		Password: "",
// 	}
// }

func getConfig() RedisInfoConfig {
	Info := RedisInfoConfig{
		Host:     "localhost",
		Port:     "6379",
		DB:       0,
		Password: "",
	}
	if globals.Config.IsSet("redis.host") {
		Info.Host = globals.Config.GetString("redis.host")
	}
	if globals.Config.IsSet("redis.port") {
		Info.Port = globals.Config.GetString("redis.port")
	}
	if globals.Config.IsSet("redis.db") {
		Info.DB = globals.Config.GetInt("redis.db")
	}
	if globals.Config.IsSet("redis.pass") {
		Info.Password = globals.Config.GetString("redis.pass")
	}
	return Info
}
