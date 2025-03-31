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

	_, err := services.LoginByPass(ctx, "19033952880", "Zhz050108")
	if err != nil {
		fmt.Println(err)
	}

	courses, err := services.GetCourses(ctx, "19033952880")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(courses[0])
	services.GetActivity(ctx, courses[0], "19033952880")
}
