package inits

import (
	"depression_go/configs"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func InitConfig() {
	err := godotenv.Load("config.env") // 加载 env 文件
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	jwtExpire, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	maxFileSize, _ := strconv.ParseInt(os.Getenv("MAX_FILE_SIZE"), 10, 64)

	configs.GlobalConfig = &configs.Config{
		Database: configs.DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		JWT: configs.JWTConfig{
			Secret:      os.Getenv("JWT_SECRET"),
			ExpireHours: jwtExpire,
		},
		BaiduAI: configs.BaiduAIConfig{
			AppID:     os.Getenv("BAIDU_APP_ID"),
			APIKey:    os.Getenv("BAIDU_API_KEY"),
			SecretKey: os.Getenv("BAIDU_SECRET_KEY"),
		},
		Server: configs.ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
			Mode: os.Getenv("SERVER_MODE"),
		},
		Upload: configs.UploadConfig{
			Path:        os.Getenv("UPLOAD_PATH"),
			MaxFileSize: maxFileSize,
		},
	}
}
