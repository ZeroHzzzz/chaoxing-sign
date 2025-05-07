package mysql

import "gorm.io/gorm"

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
	// TODO
	)
}
