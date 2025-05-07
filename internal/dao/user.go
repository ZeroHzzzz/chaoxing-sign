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
	var user models.User
	if err := d.DB.Where("id = ? AND password = ?", ID, pass).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (d *Dao) DelUserByID(ctx context.Context, ID int) error {
	var user models.User
	if err := d.DB.Delete(&user, ID).Error; err != nil {
		return err
	}
	return nil
}
