package dao

import (
	"chaoxing/internal/models"
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dao struct {
	DB  *gorm.DB
	Rdb *redis.Client
}

func NewDao(db *gorm.DB, rdb *redis.Client) Daos {
	return &Dao{
		DB:  db,
		Rdb: rdb,
	}
}

type Daos interface {
	GetChaoxingCookies(ctx context.Context, key string) (*models.ChaoxingCookieType, error)
	NewChaoxingCookies(ctx context.Context, key string, cookie models.ChaoxingCookieType) error

	NewUserSignConfig(ctx context.Context, username string, config models.SignConfigType) error
	GetUserSignConfig(ctx context.Context, username string) (*models.SignConfigType, error)
}
