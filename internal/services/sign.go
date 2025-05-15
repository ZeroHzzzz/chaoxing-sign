package services

import (
	"context"
	"log"
)

func Test(ctx context.Context, phone string, pwd string, url string) bool {
	cookie, err := c.LoginByPass(ctx, phone, pwd)
	if err != nil {
		log.Println(err)
	}
	return c.ScanQrcodeSign(ctx, cookie, url)
}
