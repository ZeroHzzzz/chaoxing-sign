package main

import (
	"chaoxing/internal/pkg/resty"
	"chaoxing/internal/services"
	"context"
	"fmt"
)

var ctx = context.Background()

func main() {
	rtyClient := resty.GetClient()
	services.ServiceInit(rtyClient)

	err := services.LoginByPass(ctx, "15918991630", "Zhz050108")
	if err != nil {
		fmt.Println(err)
	}
}
