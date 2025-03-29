package services

import (
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Rty *resty.Client
	Rdb *redis.Client
}

var svc *ServiceContext

func ServiceInit(rty *resty.Client, rdb *redis.Client) {
	svc = &ServiceContext{
		Rty: rty,
		Rdb: rdb,
	}
}
