package main

import (
	"chaoxing/internal/pkg/mysql"
	"chaoxing/internal/pkg/redis"
	"chaoxing/internal/pkg/resty"
	"chaoxing/internal/services"
	"context"
	"log"
)

// var ctx = context.Background()
var uname = "19033952880"

func main() {
	// rty := resty.GetRty()
	// rdb := redis.GetRdb()

	// services.ServiceInit(rty, rdb)

	// chao := chaoxing.NewChaoxing(rty, rdb)

	// data, err := chao.LoginByPass(ctx, uname, "Zhz050108")
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(data)
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

	// config, err := services.GetSignConfig(ctx, uname)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(config)

	// courses, err := services.GetCourses(ctx, uname)

	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(courses)

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

	// act, err := services.GetActivityLogic(ctx, courses[0], uname)
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println(act)
	// err = services.SignLogic(ctx, act[0], *config, "", uname)
	// if err != nil {
	// 	log.Println(err)
	// }

	// imparam, err := services.GetIMParams(ctx, uname)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(imparam)

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
	db := mysql.Init()
	rdb := redis.GetRdb()
	rty := resty.GetRty()
	services.Init(db, rdb, rty)
	url := "https://mobilelearn.chaoxing.com/widget/sign/e?id=2000124711806&c=2000124711806&enc=9BBD13396668E1F461BA1A7064847F3B&DB_STRATEGY=PRIMARY_KEY&STRATEGY_PARA=id"
	status := services.Test(context.Background(), uname, "Zhz050108", url)
	log.Println(status)
	// r := router.InitRouter()
	// r.Run(":8080")
}
