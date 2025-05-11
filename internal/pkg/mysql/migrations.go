package mysql

import (
	"chaoxing/internal/models"

	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.ChaoxingAccount{},
	)
}
