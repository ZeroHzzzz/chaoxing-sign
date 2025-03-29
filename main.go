package main

import (
	"chaoxing/internal/pkg/redis"
	"chaoxing/internal/pkg/resty"
	"chaoxing/internal/services"
	"context"
	"fmt"
)

var ctx = context.Background()

func main() {
	rty := resty.GetRty()
	rdb := redis.GetRdb()

	services.ServiceInit(rty, rdb)

	err := services.LoginByPass(ctx, "15918991630", "Zhz050108")
	if err != nil {
		fmt.Println(err)
	}

	_, err = services.GetPanToken(ctx, "15918991630")

	if err != nil {
		fmt.Println(err)
	}
}
