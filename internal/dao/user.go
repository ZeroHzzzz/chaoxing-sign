package dao

import (
	"chaoxing/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

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

func (d *Dao) BindChaoxingAccount(ctx context.Context, userID int, account *models.ChaoxingAccount) error {
	account.UserID = userID
	return d.DB.Create(account).Error
}

func (d *Dao) UnbindChaoxingAccount(ctx context.Context, userID int) error {
	return d.DB.Where("user_id = ?", userID).Delete(&models.ChaoxingAccount{}).Error
}

func (d *Dao) UpdateChaoxingAccount(ctx context.Context, account *models.ChaoxingAccount) error {
	return d.DB.Save(account).Error
}

func (d *Dao) GetUserChaoxingAccount(ctx context.Context, userID int) (*models.ChaoxingAccount, error) {
	var account models.ChaoxingAccount
	err := d.DB.Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}
