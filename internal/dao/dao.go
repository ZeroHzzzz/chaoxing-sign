package dao

import (
	"chaoxing/internal/models"
	"context"
	"time"

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
	// 验证码相关
	StoreVerificationCode(ctx context.Context, email, code string, expiration time.Duration) error
	GetVerificationCode(ctx context.Context, email string) (string, error)

	// 超星cookie相关
	GetChaoxingCookies(ctx context.Context, phone string) (*models.ChaoxingCookieType, error)
	NewChaoxingCookies(ctx context.Context, phone string, cookie models.ChaoxingCookieType) error

	// 签到配置相关
	NewUserSignConfig(ctx context.Context, phone string, config models.SignConfigType) error
	GetUserSignConfig(ctx context.Context, phone string) (*models.SignConfigType, error)

	// 用户相关
	NewUser(ctx context.Context, user *models.User) error
	GetUsersByUsername(ctx context.Context, username string, page, pageSize int) ([]*models.User, int64, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByIDPass(ctx context.Context, ID int, password string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByEmailPass(ctx context.Context, email, password string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
	DeleteUserByPass(ctx context.Context, ID int, password string) error

	// 超星账号相关
	BindChaoxingAccount(ctx context.Context, account *models.ChaoxingAccount) error
	GetUserChaoxingAccountByPhone(ctx context.Context, phone string, userID int) (*models.ChaoxingAccount, error)
	GetUserChaoxingAccounts(ctx context.Context, userID int, page, pageSize int) ([]*models.ChaoxingAccount, int64, error)
	UpdateChaoxingAccount(ctx context.Context, account *models.ChaoxingAccount) error
	DeleteUserChaoxingAccounts(ctx context.Context, userID int) error
	UnbindChaoxingAccount(ctx context.Context, userID int, accountID string) error

	// Group
	NewGroup(ctx context.Context, group *models.Group) error
	GetGroupByGroupID(ctx context.Context, id int) (*models.Group, error)
	GetGroupByCaptainID(ctx context.Context, captainID int) ([]*models.Group, error)
	UpdateGroup(ctx context.Context, group *models.Group) error
	DeleteGroup(ctx context.Context, id int) error
	GetGroupsByUserID(ctx context.Context, userID int) ([]*models.Group, error)
	GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error)

	AddGroupMembership(ctx context.Context, member *models.GroupMembership) error
	RemoveGroupMembership(ctx context.Context, groupID, userID int) error
	UpdateGroupMembership(ctx context.Context, member *models.GroupMembership) error
	GetGroupMemberships(ctx context.Context, groupID int) ([]*models.GroupMembership, error)
	GetGroupMembership(ctx context.Context, groupID, userID int) (*models.GroupMembership, error)
	GetGroupMembersByGroupID(ctx context.Context, groupID int) ([]*models.User, error)
}
