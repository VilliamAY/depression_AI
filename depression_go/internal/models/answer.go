package models

import (
	"time"

	"gorm.io/gorm"
)

// Answer 答案模型
type Answer struct {
	gorm.Model
	UserID       uint   `json:"user_id" gorm:"not null"`
	QuestionID   uint   `json:"question_id" gorm:"not null"`
	AssessmentID uint   `json:"assessment_id" gorm:"not null"`
	Content      string `json:"content" gorm:"type:text;not null"` // 答案内容
	Score        int    `json:"score" gorm:"default:0"`            // 答案得分

	// 关联关系
	User       User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Question   Question   `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	Assessment Assessment `json:"assessment,omitempty" gorm:"foreignKey:AssessmentID"`
}

// TableName 指定表名
func (Answer) TableName() string {
	return "answers"
}

// AnswerCreateRequest 创建答案请求
type AnswerCreateRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Score      int    `json:"score"`
}

// AnswerUpdateRequest 更新答案请求
type AnswerUpdateRequest struct {
	Content string `json:"content"`
	Score   int    `json:"score"`
}

type AnswerRequest struct {
	QuestionID  uint `json:"question_id"`
	AnswerValue int  `json:"answer_value"`
}

// AnswerResponse 答案响应
type AnswerResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	QuestionID   uint      `json:"question_id"`
	AssessmentID uint      `json:"assessment_id"`
	Content      string    `json:"content"`
	Score        int       `json:"score"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AnswerWithQuestion 包含问题信息的答案
type AnswerWithQuestion struct {
	ID           uint             `json:"id"`
	UserID       uint             `json:"user_id"`
	QuestionID   uint             `json:"question_id"`
	AssessmentID uint             `json:"assessment_id"`
	Content      string           `json:"content"`
	Score        int              `json:"score"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Question     QuestionResponse `json:"question"`
}
