package models

import (
	"time"

	"gorm.io/gorm"
)

// Question 问题模型
type Question struct {
	gorm.Model
	Title       string `json:"title" gorm:"size:500;not null"`
	Description string `json:"description" gorm:"size:1000"`
	Type        string `json:"type" gorm:"size:20;not null"`     // 问题类型：single(单选), multiple(多选), text(文本)
	Category    string `json:"category" gorm:"size:50;not null"` // 问题分类：depression(抑郁), anxiety(焦虑), stress(压力)
	Options     string `json:"options" gorm:"type:text"`         // JSON格式的选项
	Score       int    `json:"score" gorm:"default:0"`           // 问题权重分数
	OrderNum    int    `json:"order_num" gorm:"default:0"`       // 排序号
	Status      int    `json:"status" gorm:"default:1"`          // 1:启用 0:禁用

	// 关联关系
	Answers []Answer `json:"answers,omitempty" gorm:"foreignKey:QuestionID"`
}

// TableName 指定表名
func (Question) TableName() string {
	return "questions"
}

// BeforeCreate 创建前的钩子函数
func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.Status == 0 {
		q.Status = 1
	}
	return nil
}

// QuestionCreateRequest 创建问题请求
type QuestionCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Options     string `json:"options"`
	Score       int    `json:"score"`
	OrderNum    int    `json:"order_num"`
}

// QuestionUpdateRequest 更新问题请求
type QuestionUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Options     string `json:"options"`
	Score       int    `json:"score"`
	OrderNum    int    `json:"order_num"`
	Status      int    `json:"status"`
}

// QuestionResponse 问题响应
type QuestionResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Options     string    `json:"options"`
	Score       int       `json:"score"`
	OrderNum    int       `json:"order_num"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
