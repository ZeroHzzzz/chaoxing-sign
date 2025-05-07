package services

import "chaoxing/internal/models"

func GetUserByID(ID int) (*models.User, error) {
	user, err := d.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
