package main

import (
	"chaoxing/internal/pkg/redis"
	"chaoxing/internal/pkg/resty"
	"chaoxing/internal/services"
	"context"
	"log"
)

var ctx = context.Background()

func main() {
	rty := resty.GetRty()
	rdb := redis.GetRdb()

	services.ServiceInit(rty, rdb)

	_, err := services.LoginByPass(ctx, "19033952880", "Zhz050108")
	if err != nil {
		log.Println(err)
	}

	// config := models.SignConfigType{
	// 	Locations: []models.LocationType{
	// 		{
	// 			Address:   "健行楼",
	// 			Longitude: "120.043059",
	// 			Latitude:  "30.230745",
	// 		},
	// 		{
	// 			Address:   "广知楼",
	// 			Longitude: "120.044254",
	// 			Latitude:  "30.230916",
	// 		},
	// 		{
	// 			Address:   "计算机楼",
	// 			Longitude: "120.054633",
	// 			Latitude:  "30.238795",
	// 		},
	// 	},
	// }
	// err = services.StoreSignConfig(ctx, "19033952880", config)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	config, err := services.GetSignConfig(ctx, "19033952880")
	if err != nil {
		log.Println(err)
	}
	log.Println(config)

	courses, err := services.GetCourses(ctx, "19033952880")

	if err != nil {
		log.Println(err)
	}
	log.Println(courses)

	// log.Println(courses[0])
	// act, err := services.GetActivity(ctx, courses[0], "19033952880")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// log.Println(act)

	// err = services.GetPPTActivityInfo(ctx, "19033952880", &act[0])

	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(act)
	// log.Println("otherID: ", act[0].OtherID)

	act, err := services.GetActivityLogic(ctx, courses[0], "19033952880")
	if err != nil {
		log.Println(err)
	}

	err = services.SignLogic(ctx, act[0], *config, "", "19033952880")
	if err != nil {
		log.Println(err)
	}

	// imparam, err := services.GetIMParams(ctx, "19033952880")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(imparam)

	// err = services.PreSign(ctx, act.ActivityID, courses[0].CourseID, courses[0].ClassID, "19033952880")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// name, err := services.GetUserName(ctx, "19033952880")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// log.Println(name)

	// // err = services.LocationSign(ctx, name, act.ActivityID, "", "120.043053", "30.230763", "19033952880")
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	// err = services.QrcodeSign(ctx, "C81E6EDFF728672D53B6C773058E4D05", name, act.ActivityID, "", "", "", "", "19033952880")
	// if err != nil {
	// 	fmt.Println(err)
	// }

}
