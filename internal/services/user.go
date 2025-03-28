package services

import (
	"context"
)

func LoginByPass(ctx context.Context, username string, password string) error {
	err := c.LoginByPass(ctx, username, password)
	if err != nil {
		return err
	}
	return nil
}
