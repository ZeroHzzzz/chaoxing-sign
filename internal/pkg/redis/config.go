package redis

import (
	"chaoxing/internal/globals"
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
