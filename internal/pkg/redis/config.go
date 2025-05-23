package redis

import (
	"chaoxing/internal/globals/config"
)

type RedisInfoConfig struct {
	Host     string
	Port     string
	DB       int
	Password string
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
	if config.Config.IsSet("redis.host") {
		Info.Host = config.Config.GetString("redis.host")
	}
	if config.Config.IsSet("redis.port") {
		Info.Port = config.Config.GetString("redis.port")
	}
	if config.Config.IsSet("redis.db") {
		Info.DB = config.Config.GetInt("redis.db")
	}
	if config.Config.IsSet("redis.pass") {
		Info.Password = config.Config.GetString("redis.pass")
	}
	return Info
}
