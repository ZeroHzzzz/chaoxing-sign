package services

import (
	"chaoxing/internal/models"
	"chaoxing/internal/pkg/xerr"
	"chaoxing/internal/utils"
	"errors"
)

func Register(username, password string) error {
	// 创建新用户
	user := &models.User{
		Username: username,
		Password: password, // 实际应用中应该对密码进行加密
	}

	return d.NewUser(ctx, user)
}

func Login(ID int, password string) (string, error) {
	// 获取用户信息
	user, err := d.GetUserByIDPass(ctx, ID, password)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", xerr.OtherError("用户不存在")
	}

	// 验证密码
	if user.Password != password { // 实际应用中应该对密码进行验证
		return "", xerr.OtherError("密码错误")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserByID(ID int) (*models.User, error) {
	// 获取用户信息
	user, err := d.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerr.OtherError("用户不存在")
	}

	return user, nil
}

func GetUsersByUsername(username string, page, pageSize int) ([]*models.User, int64, error) {
	// 获取用户列表
	users, totalCount, err := d.GetUsersByUsername(ctx, username, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func DeleteUserByPass(ID int, password string) error {
	// 获取用户信息
	user, err := d.GetUserByIDPass(ctx, ID, password)
	if err != nil {
		return err
	}
	if user == nil {
		return xerr.OtherError("用户不存在")
	}

	// 删除用户
	return d.DeleteUserByPass(ctx, ID, password)
}

func BindChaoxingAccount(userID int, phone, password string) error {
	// 调用超星登录接口验证账号
	cookie, err := c.LoginByPass(ctx, phone, password)
	if err != nil {
		return xerr.OtherError("超星账号验证失败")
	}

	// 获取用户名
	name, err := c.GetUserName(ctx, cookie)
	if err != nil {
		return err
	}

	// 创建超星账号记录
	account := &models.ChaoxingAccount{
		Phone: phone,
		Pass:  password,
		Name:  name,
	}

	// 绑定账号
	err = d.BindChaoxingAccount(ctx, userID, account)
	if err != nil {
		return err
	}

	return nil
}

func GetUserChaoxingAccount(userID int) (*models.ChaoxingAccount, error) {
	return d.GetUserChaoxingAccount(ctx, userID)
}

func UpdateChaoxingAccount(userID int, phone, password string) error {
	// 获取现有账号信息
	account, err := d.GetUserChaoxingAccount(ctx, userID)
	if err != nil {
		return err
	}
	if account == nil {
		return xerr.OtherError("未绑定超星账号")
	}

	// 验证新账号信息
	cookie, err := c.LoginByPass(ctx, phone, password)
	if err != nil {
		return xerr.OtherError("超星账号验证失败")
	}

	// 获取用户名
	name, err := c.GetUserName(ctx, cookie)
	if err != nil {
		return err
	}

	// 更新账号信息
	account.Phone = phone
	account.Pass = password
	account.Name = name

	return d.UpdateChaoxingAccount(ctx, account)
}

func UnbindChaoxingAccount(userID int) error {
	return d.UnbindChaoxingAccount(ctx, userID)
}

// ValidateUser 验证用户token并获取用户信息
func ValidateUser(token string) (*models.User, error) {
	// 解析token
	claims, err := utils.ParseToken(token)
	if err != nil {
		return nil, xerr.NotLoginErr
	}

	// 获取用户信息
	user, err := d.GetUserByID(ctx, claims.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	return user, nil
}
