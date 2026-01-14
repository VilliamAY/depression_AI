package models

import (
	"time"

	"gorm.io/gorm"
)

// Assessment 评估模型
type Assessment struct {
	gorm.Model
	UserID     uint   `json:"user_id" gorm:"not null"`
	Title      string `json:"title" gorm:"size:200;not null"` // 评估标题
	Type       string `json:"type" gorm:"size:50;not null"` // 评估类型：questionnaire(问卷), face(人脸), combined(综合)
	TotalScore int    `json:"total_score" gorm:"default:0"` // 总分数
	MaxScore   int    `json:"max_score" gorm:"default:0"` // 最高可能分数
	Level      string `json:"level" gorm:"size:20"` // 评估等级：normal(正常), mild(轻度), moderate(中度), severe(重度)
	Result     string `json:"result" gorm:"type:text"` // 评估结果描述
	Status     int    `json:"status" gorm:"default:1"` // 1:完成 0:进行中

	// 关联关系
	User    User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Answers []Answer `json:"answers,omitempty" gorm:"foreignKey:AssessmentID"`
}

// TableName 指定表名
func (Assessment) TableName() string {
	return "assessments"
}

// BeforeCreate 创建前的钩子函数
func (a *Assessment) BeforeCreate(tx *gorm.DB) error {
	if a.Status == 0 {
		a.Status = 1
	}
	return nil
}

// AssessmentCreateRequest 创建评估请求
type AssessmentCreateRequest struct {
	Title string `json:"title" binding:"required"`
	Type  string `json:"type" binding:"required"`
}

// AssessmentUpdateRequest 更新评估请求
type AssessmentUpdateRequest struct {
	Title      string `json:"title"`
	TotalScore int    `json:"total_score"`
	MaxScore   int    `json:"max_score"`
	Level      string `json:"level"`
	Result     string `json:"result"`
	Status     int    `json:"status"`
}

// AssessmentResponse 评估响应
type AssessmentResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	TotalScore int       `json:"total_score"`
	MaxScore   int       `json:"max_score"`
	Level      string    `json:"level"`
	Result     string    `json:"result"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AssessmentWithAnswers 包含答案的评估
type AssessmentWithAnswers struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	TotalScore int       `json:"total_score"`
	MaxScore   int       `json:"max_score"`
	Level      string    `json:"level"`
	Result     string    `json:"result"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Answers    []AnswerWithQuestion `json:"answers"`
}

// AssessmentResult 评估结果
type AssessmentResult struct {
	Level       string  `json:"level"`
	Score       int     `json:"score"`
	MaxScore    int     `json:"max_score"`
	Percentage  float64 `json:"percentage"`
	Description string  `json:"description"`
	Suggestions string  `json:"suggestions"`
} 