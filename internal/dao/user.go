package dao

import (
	"chaoxing/internal/models"
	"context"
)

func (d *Dao) NewUser(ctx context.Context, user *models.User) error {
	if err := d.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *Dao) GetUserOnlyByID(ctx context.Context, ID int) (*models.User, error) {
	var user models.User
	if err := d.DB.First(&user, ID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Dao) GetUserOnlyByIDPass(ctx context.Context, ID int, pass string) (*models.User, error) {
	var user models.User
	if err := d.DB.Where("id = ? AND password = ?", ID, pass).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Dao) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	var user models.User
	if err := d.DB.Preload("ChaoxingAccounts").First(&user, ID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Dao) GetUserByIDPass(ctx context.Context, ID int, pass string) (*models.User, error) {
	var user models.User
	if err := d.DB.Where("id = ? AND password = ?", ID, pass).Preload("ChaoxingAccounts").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Dao) UpdateUser(ctx context.Context, user *models.User) error {
	if err := d.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *Dao) DelUser(ctx context.Context, ID int, pass string) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		// 先查找用户
		var user models.User
		if err := tx.Where("id = ? AND password = ?", ID, pass).First(&user).Error; err != nil {
			return err
		}
		// 删除关联的超星账号
		if err := tx.Where("user_id = ?", ID).Delete(&models.ChaoxingAccount{}).Error; err != nil {
			return err
		}
		// 删除用户
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *Dao) DelUserByID(ctx context.Context, ID int) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		// 先查找用户
		var user models.User
		if err := tx.First(&user, ID).Error; err != nil {
			return err
		}
		// 删除关联的超星账号
		if err := tx.Where("user_id = ?", ID).Delete(&models.ChaoxingAccount{}).Error; err != nil {
			return err
		}
		// 删除用户
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		return nil
	})
}
