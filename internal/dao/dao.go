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
	NewChaoxingUser(ctx context.Context, user *models.ChaoxingAccount) error
	GetChaoxingUserByPhone(ctx context.Context, phone string) (*models.ChaoxingAccount, error)
	UpdateChaoxingUser(ctx context.Context, user *models.ChaoxingAccount) error
	DelChaoxingUser(ctx context.Context, phone string) error

	GetChaoxingCookies(ctx context.Context, phone string) (*models.ChaoxingCookieType, error)
	NewChaoxingCookies(ctx context.Context, phone string, cookie models.ChaoxingCookieType) error

	NewUserSignConfig(ctx context.Context, phone string, config models.SignConfigType) error
	GetUserSignConfig(ctx context.Context, phone string) (*models.SignConfigType, error)

	NewUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, ID int) (*models.User, error)
	GetUserOnlyByID(ctx context.Context, ID int) (*models.User, error)
	GetUserByIDPass(ctx context.Context, ID int, pass string) (*models.User, error)
	GetUserOnlyByIDPass(ctx context.Context, ID int, pass string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DelUser(ctx context.Context, ID int, pass string) error
	DelUserByID(ctx context.Context, ID int) error
}
