package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	// 使用bcrypt进行密码加密，成本因子为12
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hashedPassword string) bool {
	// 比较密码和哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) (bool, string) {
	if len(password) < 6 {
		return false, "密码长度至少6位"
	}
	
	if len(password) > 50 {
		return false, "密码长度不能超过50位"
	}
	
	return true, ""
} 