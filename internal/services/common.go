package services

import (
	"chaoxing/internal/chao"

	"github.com/go-resty/resty/v2"
)

var c *chao.Chao

func ServiceInit(rtyClient *resty.Client) {
	c = chao.NewChao(rtyClient)
}
