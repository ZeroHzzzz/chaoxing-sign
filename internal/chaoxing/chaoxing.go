package chaoxing

import (
	"chaoxing/internal/models"
	"context"

	"github.com/go-resty/resty/v2"
)

type Chaoxing struct {
	Rty *resty.Client
}

func NewChaoxing(rty *resty.Client, cookie *models.ChaoxingCookieType) Chaoxings {
	return &Chaoxing{
		Rty: rty,
	}
}

type Chaoxings interface {
	// user
	// GetCookies(ctx context.Context, key string) (*models.ChaoxingCookieType, error)
	// StoreCookies(ctx context.Context, key string, cookie models.ChaoxingCookieType) error
	LoginByPass(ctx context.Context, phone string, password string) (models.ChaoxingCookieType, error)
	GetPanToken(ctx context.Context, cookie models.ChaoxingCookieType) (string, error)
	GetCourses(ctx context.Context, cookie models.ChaoxingCookieType) ([]models.CourseType, error)
	GetUserName(ctx context.Context, cookie models.ChaoxingCookieType) (string, error)
	GetIMParams(ctx context.Context, cookie models.ChaoxingCookieType) (*models.IMParamsType, error)

	// StoreSignConfig(ctx context.Context, phone string, config models.SignConfigType) error
	// GetSignConfig(ctx context.Context, phone string) (*models.SignConfigType, error)

	// sign
	GetPPTActivityInfo(ctx context.Context, cookie models.ChaoxingCookieType, activity *models.ActivityType) error
	GetActivity(ctx context.Context, cookie models.ChaoxingCookieType, course models.CourseType) ([]models.ActivityType, error)
	GetActivityLogic(ctx context.Context, cookie models.ChaoxingCookieType, course models.CourseType) ([]models.ActivityType, error)

	SignLogic(ctx context.Context, cookie models.ChaoxingCookieType, act models.ActivityType, signCfg models.SignConfigType, enc, signCode string) error
	PreSign(ctx context.Context, cookie models.ChaoxingCookieType, act models.ActivityType) bool
	GeneralSign(ctx context.Context, cookie models.ChaoxingCookieType, act models.ActivityType, phone string) bool
	CodeSign(ctx context.Context, cookie models.ChaoxingCookieType, act models.ActivityType, signCode, phone string) bool
	QrcodeSign(ctx context.Context, cookie models.ChaoxingCookieType, location models.LocationType, enc, name, activeId string) bool
	LocationSign(ctx context.Context, cookie models.ChaoxingCookieType, location models.LocationType, name, activeId string) bool
}
