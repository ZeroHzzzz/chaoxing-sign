package chaoxing

import (
	"chaoxing/internal/models"
	"context"

	"github.com/go-resty/resty/v2"
)

type Chaoxing struct {
	Rty    *resty.Client
	Cookie *models.ChaoxingCookieType
}

func NewChaoxing(rty *resty.Client, cookie *models.ChaoxingCookieType) Chaoxings {
	return &Chaoxing{
		Rty:    rty,
		Cookie: cookie,
	}
}

type Chaoxings interface {
	// user
	// GetCookies(ctx context.Context, key string) (*models.ChaoxingCookieType, error)
	// StoreCookies(ctx context.Context, key string, cookie models.ChaoxingCookieType) error
	UpdateCookie(cookie models.ChaoxingCookieType)
	LoginByPass(ctx context.Context, phone string, password string) (models.ChaoxingCookieType, error)
	GetPanToken(ctx context.Context) (string, error)
	GetCourses(ctx context.Context) ([]models.CourseType, error)
	GetUserName(ctx context.Context) (string, error)
	GetIMParams(ctx context.Context) (*models.IMParamsType, error)

	// StoreSignConfig(ctx context.Context, phone string, config models.SignConfigType) error
	// GetSignConfig(ctx context.Context, phone string) (*models.SignConfigType, error)

	// sign
	GetPPTActivityInfo(ctx context.Context, activity *models.ActivityType) error
	GetActivity(ctx context.Context, course models.CourseType) ([]models.ActivityType, error)
	GetActivityLogic(ctx context.Context, course models.CourseType) ([]models.ActivityType, error)

	SignLogic(ctx context.Context, act models.ActivityType, signCfg models.SignConfigType, enc, signCode string) error
	PreSign(ctx context.Context, act models.ActivityType) bool
	GeneralSign(ctx context.Context, act models.ActivityType, phone string) bool
	CodeSign(ctx context.Context, act models.ActivityType, signCode, phone string) bool
	QrcodeSign(ctx context.Context, location models.LocationType, enc, name, activeId string) bool
	LocationSign(ctx context.Context, location models.LocationType, name, activeId string) bool
}
