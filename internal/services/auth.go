package services

import (
	"chaoxing/internal/models"
	"chaoxing/internal/pkg/email"
	"chaoxing/internal/pkg/utils"
	"chaoxing/internal/pkg/verification"
	"chaoxing/internal/pkg/xerr"
	"context"
)

func SendVerificationCode(ctx context.Context, emailAddr string) error {
	code := verification.GenerateCode()

	err := d.StoreVerificationCode(ctx, emailAddr, code, verification.CodeTTL)
	if err != nil {
		return err
	}

	return email.SendVerificationCode(emailAddr, code)
}

func RegisterByEmail(ctx context.Context, username, email, password, code string) error {
	// 验证邮箱验证码
	storedCode, err := d.GetVerificationCode(ctx, email)
	if err != nil {
		return err
	}

	if !verification.VerifyCode(storedCode, code) {
		return xerr.EmailVerifyErr
	}

	// 创建新用户
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	return d.NewUser(ctx, user)
}

func RegisterTest(ctx context.Context, username, email, password string) error {
	// 创建新用户
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	return d.NewUser(ctx, user)
}

func LoginByID(ctx context.Context, ID int, password string) (string, *models.User, error) {
	user, err := d.GetUserByIDPass(ctx, ID, password)
	if err != nil {
		return "", nil, err
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func LoginByEmail(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := d.GetUserByEmailPass(ctx, email, password)
	if err != nil {
		return "", nil, err
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	// 获取用户信息
	user, err := d.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
