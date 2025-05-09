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
	GetChaoxingCookies(ctx context.Context, phone string) (*models.ChaoxingCookieType, error)
	NewChaoxingCookies(ctx context.Context, phone string, cookie models.ChaoxingCookieType) error

	NewUserSignConfig(ctx context.Context, phone string, config models.SignConfigType) error
	GetUserSignConfig(ctx context.Context, phone string) (*models.SignConfigType, error)

	NewUser(ctx context.Context, user *models.User) error
	GetUsersByUsername(ctx context.Context, username string, page, pageSize int) ([]*models.User, int64, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByIDPass(ctx context.Context, ID int, password string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByEmailPass(ctx context.Context, email, password string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
	DeleteUserByPass(ctx context.Context, ID int, password string) error

	BindChaoxingAccount(ctx context.Context, userID int, account *models.ChaoxingAccount) error
	GetUserChaoxingAccount(ctx context.Context, userID int) (*models.ChaoxingAccount, error)
	UpdateChaoxingAccount(ctx context.Context, account *models.ChaoxingAccount) error
	UnbindChaoxingAccount(ctx context.Context, userID int) error
}
