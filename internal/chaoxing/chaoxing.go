package chaoxing

import (
	"chaoxing/internal/models"
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
)

type Chaoxing struct {
	Rty *resty.Client
	Rdb *redis.Client
}

func NewChaoxing(ctx context.Context, rty *resty.Client, rdb *redis.Client) Chaoxings {
	chao := &Chaoxing{
		Rty: rty,
		Rdb: rdb,
	}
	return chao
}

type Chaoxings interface {
	// user
	GetCookies(ctx context.Context, key string) (*models.ChaoxingCookieType, error)
	StoreCookies(ctx context.Context, key string, cookie models.ChaoxingCookieType) error
	LoginByPass(ctx context.Context, username string, password string) (models.ChaoxingCookieType, error)
	GetPanToken(ctx context.Context, username string) (string, error)
	GetCourses(ctx context.Context, username string) ([]models.CourseType, error)
	GetUserName(ctx context.Context, username string) (string, error)
	GetIMParams(ctx context.Context, username string) (*models.IMParamsType, error)

	StoreSignConfig(ctx context.Context, username string, config models.SignConfigType) error
	GetSignConfig(ctx context.Context, username string) (*models.SignConfigType, error)

	// sign
	GetPPTActivityInfo(ctx context.Context, username string, activity *models.ActivityType) error
	GetActivity(ctx context.Context, course models.CourseType, username string) ([]models.ActivityType, error)
	GetActivityLogic(ctx context.Context, course models.CourseType, username string) ([]models.ActivityType, error)
	SignLogic(ctx context.Context, act models.ActivityType, signCfg models.SignConfigType, enc, signCode, username string) error
	PreSign(ctx context.Context, act models.ActivityType, username string) bool
	GeneralSign(ctx context.Context, act models.ActivityType, username string) bool
	CodeSign(ctx context.Context, act models.ActivityType, signCode, username string) bool
	QrcodeSign(ctx context.Context, location models.LocationType, enc, name, activeId, username string) bool
	LocationSign(ctx context.Context, location models.LocationType, name, activeId, username string) bool
}
