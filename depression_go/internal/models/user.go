package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email    string `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Password string `json:"-" gorm:"size:255;not null"` // 密码不返回给前端
	Age      int    `json:"age" gorm:"default:0"`
	Gender   string `json:"gender" gorm:"size:10;default:'未知'"`
	Phone    string `json:"phone" gorm:"size:20"`
	Avatar   string `json:"avatar" gorm:"size:255"`
	Status   int    `json:"status" gorm:"default:1"` // 1:正常 0:禁用

	// 关联关系
	Assessments    []Assessment    `json:"assessments,omitempty" gorm:"foreignKey:UserID"`
	FaceDetections []FaceDetection `json:"face_detections,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子函数
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Status == 0 {
		u.Status = 1
	}
	return nil
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest 用户信息更新请求
type UserUpdateRequest struct {
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	Phone  string `json:"phone"`
	Avatar string `json:"avatar"`
}

// UserResponse 用户信息响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
