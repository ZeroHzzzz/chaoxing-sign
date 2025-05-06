package mysql

import (
	"chaoxing/internal/globals"
	"chaoxing/internal/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		globals.Config.GetString("mysql.username"),
		globals.Config.GetString("mysql.password"),
		globals.Config.GetString("mysql.host"),
		globals.Config.GetInt("mysql.port"),
		globals.Config.GetString("mysql.database"),
		globals.Config.GetString("mysql.charset"),
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("MySQL连接失败: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(globals.Config.GetInt("mysql.max_idle_conns"))
	sqlDB.SetMaxOpenConns(globals.Config.GetInt("mysql.max_open_conns"))
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据表结构
	autoMigrate()

	log.Println("MySQL连接成功")
}

// autoMigrate 自动迁移数据表结构
func autoMigrate() {
	// 在这里添加需要自动迁移的模型
	err := DB.AutoMigrate(
		&models.User{},
		&models.ChaoxingUser{},
	)

	if err != nil {
		log.Fatalf("数据表迁移失败: %v", err)
	}

	log.Println("数据表迁移成功")
}
