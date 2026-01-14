package handlers

import (
	"depression_go/inits"
	"depression_go/internal/models"
	"depression_go/middleware"
	"depression_go/pkg/response"
	"depression_go/pkg/token"
	"depression_go/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		db: inits.DB,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegisterRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证密码强度
	if valid, msg := utils.ValidatePassword(req.Password); !valid {
		response.BadRequest(c, msg)
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		response.BadRequest(c, "用户名已存在")
		return
	}

	// 检查邮箱是否已存在
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		response.BadRequest(c, "邮箱已存在")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		response.InternalServerError(c, "密码加密失败")
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Age:      req.Age,
		Gender:   req.Gender,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := h.db.Create(&user).Error; err != nil {
		response.InternalServerError(c, "用户创建失败")
		return
	}

	// 生成JWT令牌
	tokenString, err := token.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.InternalServerError(c, "令牌生成失败")
		return
	}

	// 返回用户信息（不包含密码）
	userResponse := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		Gender:    user.Gender,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}

	response.SuccessWithMessage(c, "注册成功", gin.H{
		"user":  userResponse,
		"token": tokenString,
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 查找用户
	var user models.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		response.BadRequest(c, "用户名或密码错误")
		return
	}

	// 检查用户状态
	if user.Status == 0 {
		response.Forbidden(c, "账户已被禁用")
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		response.BadRequest(c, "用户名或密码错误")
		return
	}

	// 生成JWT令牌
	tokenString, err := token.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.InternalServerError(c, "令牌生成失败")
		return
	}

	// 返回用户信息（不包含密码）
	userResponse := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		Gender:    user.Gender,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}

	response.SuccessWithMessage(c, "登录成功", gin.H{
		"user":  userResponse,
		"token": tokenString,
	})
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 返回用户信息（不包含密码）
	userResponse := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		Gender:    user.Gender,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}

	response.Success(c, userResponse)
}
