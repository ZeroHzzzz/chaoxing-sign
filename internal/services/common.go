package services

import (
	"chaoxing/internal/chaoxing"
	"chaoxing/internal/dao"
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ctx = context.Background()
	d   dao.Daos
	c   chaoxing.Chaoxings
)

// Init 函数用于初始化服务。
func Init(db *gorm.DB, rdb *redis.Client, rty *resty.Client) {
	d = dao.NewDao(db, rdb)
	c = chaoxing.NewChaoxing(rty)
}
