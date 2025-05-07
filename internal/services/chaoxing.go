package services

import "chaoxing/internal/models"

func NewChaoxingUser(user *models.ChaoxingAccount) error {
	err := d.NewChaoxingUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func GetChaoxingUserByPhone(phone string) (*models.ChaoxingAccount, error) {
	user, err := d.GetChaoxingUserByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	return user, nil
}
