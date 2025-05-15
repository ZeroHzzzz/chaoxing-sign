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

func BatchScanQrCodeSign(ctx context.Context, phone []string, url string) ([]bool, error) {
	var results []bool
	for _, p := range phone {
		cookie, err := d.GetChaoxingCookies(ctx, p)
		if err != nil {
			log.Println(err)
			results = append(results, false)
			continue
		}
		result := c.ScanQrcodeSign(ctx, *cookie, url)
		results = append(results, result)
	}
	return results, nil
}
