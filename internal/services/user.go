package services

import "chaoxing/internal/models"

func GetUserByID(ID int) (*models.User, error) {
	user, err := d.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByIDPass(ID int, pass string) (*models.User, error) {
	user, err := d.GetUserByIDPass(ctx, ID, pass)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserOnly(ID int) (*models.User, error) {
	user, err := d.GetUserOnlyByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserOnlyByIDPass(ID int, pass string) (*models.User, error) {
	user, err := d.GetUserOnlyByIDPass(ctx, ID, pass)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(user *models.User) error {
	err := d.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func DelUser(ID int, pass string) error {
	err := d.DelUser(ctx, ID, pass)
	if err != nil {
		return err
	}
	return nil
}