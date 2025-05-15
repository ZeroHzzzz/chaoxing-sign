package services

import (
	"chaoxing/internal/models"
	"context"
)

func GetUsersByUsername(ctx context.Context, username string, page, pageSize int) ([]*models.User, int64, error) {
	// 获取用户列表
	users, totalCount, err := d.GetUsersByUsername(ctx, username, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// 获取用户信息
	user, err := d.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmailPass(ctx context.Context, email, password string) (*models.User, error) {
	// 获取用户信息
	user, err := d.GetUserByEmailPass(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUserByPass(ctx context.Context, ID int, password string) error {
	return d.DeleteUserByPass(ctx, ID, password)
}

func BindChaoxingAccount(ctx context.Context, userID int, phone, password string) error {
	// 调用超星登录接口验证账号
	cookie, err := c.LoginByPass(ctx, phone, password)
	if err != nil {
		return err
	}

	// 获取用户名
	name, err := c.GetUserName(ctx, cookie)
	if err != nil {
		return err
	}

	// 创建超星账号记录
	account := &models.ChaoxingAccount{
		UserID: userID,
		Phone:  phone,
		Pass:   password,
		Name:   name,
	}

	// 绑定账号
	err = d.BindChaoxingAccount(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func UpdateChaoxingAccount(ctx context.Context, userID int, phone, password string) error {
	// 验证新账号信息
	cookie, err := c.LoginByPass(ctx, phone, password)
	if err != nil {
		return err
	}

	// 获取用户名
	name, err := c.GetUserName(ctx, cookie)
	if err != nil {
		return err
	}

	// 更新账号信息
	account := &models.ChaoxingAccount{
		UserID: userID,
		Phone:  phone,
		Pass:   password,
		Name:   name,
	}

	return d.UpdateChaoxingAccount(ctx, account)
}

func GetUserBindChaoxingAccounts(ctx context.Context, userID int, page, pageSize int) ([]*models.ChaoxingAccount, int64, error) {
	// 获取绑定的超星账号列表
	accounts, totalCount, err := d.GetUserChaoxingAccounts(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return accounts, totalCount, nil
}

func GetUserChaoxingAccountByPhone(ctx context.Context, phone string, userID int) (*models.ChaoxingAccount, error) {
	// 获取绑定的超星账号
	account, err := d.GetUserChaoxingAccountByPhone(ctx, phone, userID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func DeleteUserChaoxingAccounts(ctx context.Context, userID int) error {
	// 删除用户绑定的超星账号
	return d.DeleteUserChaoxingAccounts(ctx, userID)
}

func UnbindChaoxingAccount(ctx context.Context, userID int, accountID string) error {
	// 解除绑定
	return d.UnbindChaoxingAccount(ctx, userID, accountID)
}

// ValidateUser 验证用户token并获取用户信息
// func ValidateUser(ctx context.Context, token string) (*models.User, error) {
// 	// 解析token
// 	claims, err := utils.ParseToken(token)
// 	if err != nil {
// 		return nil, xerr.NotLoginErr
// 	}

// 	// 获取用户信息
// 	user, err := d.GetUserByID(ctx, claims.ID)
// 	if err != nil {

// 	// 获取用户信息
// 	user, err := d.GetUserByID(ctx, claims.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user == nil {
// 		return nil, errors.New("用户不存在")
// 	}

// 	return user, nil
// }
