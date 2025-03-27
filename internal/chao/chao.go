package chao

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type Chao struct {
	rty *resty.Client
}

func NewChao(rty *resty.Client) *Chao {
	return &Chao{rty: rty}
}

type Chaos interface {
	LoginByPass(ctx context.Context, username string, password string) error
	LoginByCode(ctx context.Context, username string, code string) error
}
