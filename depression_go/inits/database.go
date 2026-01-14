package inits

import (
	"fmt"
	"log"
	"os"

	"depression_go/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	// 从环境变量获取数据库配置
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// 检查必要的环境变量
	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatalf("数据库配置不完整 - Host: %s, Port: %s, User: %s, Database: %s", dbHost, dbPort, dbUser, dbName)
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// 配置GORM日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // 只显示警告和错误，不显示每条SQL
	}

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 获取底层sql.DB对象
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100) // 最大打开连接数

	// 自动迁移数据库表结构
	AutoMigrate()
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() {
	// 迁移所有模型
	err := DB.AutoMigrate(
		&models.User{},
		&models.Question{},
		&models.Assessment{},
		&models.Answer{},
		&models.FaceDetection{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
