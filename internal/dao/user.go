package dao

import (
	"chaoxing/internal/models"
	"context"
)

func (c *Dao) NewUser(ctx context.Context, user *models.User) error {
	if err := c.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (c *Dao) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	var user models.User
	if err := c.DB.First(&user, ID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Dao) GetUserByIDPass(ctx context.Context, ID int, pass string) (*models.User, error) {
	var user models.User
	if err := c.DB.Where("id = ? AND password = ?", ID, pass).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Dao) UpdateUser(ctx context.Context, user *models.User) error {
	if err := c.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (c *Dao) DelUser(ctx context.Context, ID int, pass string) error {
	var user models.User
	if err := c.DB.Where("id = ? AND password = ?", ID, pass).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (c *Dao) DelUserByID(ctx context.Context, ID int) error {
	var user models.User
	if err := c.DB.Delete(&user, ID).Error; err != nil {
		return err
	}
	return nil
}
