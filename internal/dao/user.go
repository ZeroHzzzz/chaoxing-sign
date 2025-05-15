package dao

import (
	"chaoxing/internal/models"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// StoreVerificationCode 存储邮箱验证码到Redis
func (d *Dao) StoreVerificationCode(ctx context.Context, email, code string, expiration time.Duration) error {
	return d.Rdb.Set(ctx, "verification_code:"+email, code, expiration).Err()
}

// GetVerificationCode 从Redis获取邮箱验证码
func (d *Dao) GetVerificationCode(ctx context.Context, email string) (string, error) {
	code, err := d.Rdb.Get(ctx, "verification_code:"+email).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("验证码不存在或已过期")
		}
		return "", err
	}
	return code, nil
}

func (d *Dao) NewUser(ctx context.Context, user *models.User) error {
	return d.DB.Create(user).Error
}

func (d *Dao) GetUsersByUsername(ctx context.Context, username string, page, pageSize int) ([]*models.User, int64, error) {
	var (
		users      []*models.User
		totalCount int64
	)

	// 构建基础查询条件
	db := d.DB.Model(&models.User{}).Where("username LIKE ?", "%"+username+"%")

	// 查询总数（用于分页）
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 防止过大分页
	}

	err := db.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (d *Dao) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := d.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *Dao) GetUserByIDPass(ctx context.Context, ID int, pass string) (*models.User, error) {
	var user models.User
	err := d.DB.Where("id = ? AND password = ?", ID, pass).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *Dao) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := d.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *Dao) GetUserByEmailPass(ctx context.Context, email, pass string) (*models.User, error) {
	var user models.User
	err := d.DB.Where("email = ? AND password = ?", email, pass).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *Dao) UpdateUser(ctx context.Context, user *models.User) error {
	return d.DB.Model(user).Updates(user).Error
}

func (d *Dao) DeleteUser(ctx context.Context, id int) error {
	return d.DB.Delete(&models.User{}, id).Error
}

func (d *Dao) DeleteUserByPass(ctx context.Context, ID int, password string) error {
	return d.DB.Where("id = ? AND password = ?", ID, password).Delete(&models.User{}).Error
}
