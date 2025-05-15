package dao

import (
	"chaoxing/internal/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func (d *Dao) NewChaoxingCookies(ctx context.Context, key string, cookie models.ChaoxingCookieType) error {
	cookieJSON, err := json.Marshal(cookie)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return err
	}

	err = d.Rdb.Set(ctx, "cookie:"+key, cookieJSON, 2*time.Hour).Err() // 两小时过期
	if err != nil {
		return err
	}

	return nil
}

func (d *Dao) GetChaoxingCookies(ctx context.Context, key string) (*models.ChaoxingCookieType, error) {
	val, err := d.Rdb.Get(ctx, "cookie:"+key).Result()
	if err != nil {
		if err == redis.Nil {
			// todo: 处理 Cookie 过期的情况
			return nil, err
		}
		log.Printf("获取 Cookie 失败: %v\n", err)
		return nil, err
	}

	var cookie models.ChaoxingCookieType
	err = json.Unmarshal([]byte(val), &cookie)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return nil, err
	}

	return &cookie, nil
}

func (d *Dao) NewUserSignConfig(ctx context.Context, phone string, config models.SignConfigType) error {
	// 暂时考虑将用户cookie和签到配置分开存储

	configJSON, err := json.Marshal(config)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return err
	}

	err = d.Rdb.Set(ctx, "sign_config:"+phone, configJSON, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) GetUserSignConfig(ctx context.Context, phone string) (*models.SignConfigType, error) {
	val, err := d.Rdb.Get(ctx, "sign_config:"+phone).Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("签到配置不存在")
			return nil, nil
		}
		log.Printf("获取签到配置失败: %v\n", err)
		return nil, err
	}

	var config models.SignConfigType
	err = json.Unmarshal([]byte(val), &config)
	if err != nil {
		log.Printf("数据转换失败: %v\n", err)
		return nil, err
	}

	return &config, nil
}

// GetChaoxingAccountByPhone 通过手机号查询超星账号
func (d *Dao) GetUserChaoxingAccountByPhone(ctx context.Context, phone string, userID int) (*models.ChaoxingAccount, error) {
	var account models.ChaoxingAccount
	err := d.DB.Where("phone = ? AND user_id = ?", phone, userID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (d *Dao) GetUserChaoxingAccounts(ctx context.Context, userID int, page, pageSize int) ([]*models.ChaoxingAccount, int64, error) {
	var (
		accounts   []*models.ChaoxingAccount
		totalCount int64
	)

	// 构建基础查询条件
	db := d.DB.Model(&models.ChaoxingAccount{}).Where("user_id = ?", userID)

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
		Find(&accounts).Error

	if err != nil {
		return nil, 0, err
	}

	return accounts, totalCount, nil
}

// DeleteUserChaoxingAccounts 删除用户的所有超星账号
func (d *Dao) DeleteUserChaoxingAccounts(ctx context.Context, userID int) error {
	return d.DB.Where("user_id = ?", userID).Delete(&models.ChaoxingAccount{}).Error
}

// UpdateChaoxingAccountByPhone 通过手机号更新超星账号信息
func (d *Dao) UpdateChaoxingAccount(ctx context.Context, account *models.ChaoxingAccount) error {
	return d.DB.Model(account).Updates(account).Error
}

func (d *Dao) UnbindChaoxingAccount(ctx context.Context, userID int, accountID string) error {
	return d.DB.Where("user_id = ? AND id = ?", userID, accountID).Delete(&models.ChaoxingAccount{}).Error
}

func (d *Dao) BindChaoxingAccount(ctx context.Context, account *models.ChaoxingAccount) error {
	return d.DB.Create(account).Error
}
