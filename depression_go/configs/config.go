package configs

// Config 全局配置结构体
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	BaiduAI  BaiduAIConfig
	Server   ServerConfig
	Upload   UploadConfig
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// BaiduAIConfig 百度AI配置
type BaiduAIConfig struct {
	AppID     string
	APIKey    string
	SecretKey string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Path        string
	MaxFileSize int64
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config
