package services

func AddChaoxingAccount(user *models.ChaoxingAccount) error {
	err := d.AddChaoxingAccount(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func GetChaoxingAccountsByUserID(userID int) ([]models.ChaoxingAccount, error) {
	accounts, err := d.GetChaoxingAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func DelChaoxingAccount(phone string) error {
	err := d.DelChaoxingAccount(ctx, phone)
	if err != nil {
		return err
	}
	return nil
}

func DelChaoxingAccountByID(id int) error {
	err := d.DelChaoxingAccountByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChaoxingAccount(account *models.ChaoxingAccount) error {
	err := d.UpdateChaoxingAccount(ctx, account)
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
