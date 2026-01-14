package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"depression_go/inits"
	"depression_go/internal/models"
	"depression_go/middleware"
	"depression_go/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// QuestionnaireHandler 问卷处理器
type QuestionnaireHandler struct {
	db *gorm.DB
}

// NewQuestionnaireHandler 创建问卷处理器
func NewQuestionnaireHandler() *QuestionnaireHandler {
	return &QuestionnaireHandler{
		db: inits.DB,
	}
}

// GetQuestions 获取问题列表
func (h *QuestionnaireHandler) GetQuestions(c *gin.Context) {
	// 获取查询参数
	category := c.Query("category")
	status := c.Query("status")

	// 构建查询条件
	query := h.db.Model(&models.Question{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "" {
		if statusInt, err := strconv.Atoi(status); err == nil {
			query = query.Where("status = ?", statusInt)
		}
	}

	// 查询问题
	var questions []models.Question
	if err := query.Order("order_num ASC, id ASC").Find(&questions).Error; err != nil {
		response.InternalServerError(c, "查询失败")
		return
	}

	// 转换为响应格式
	var responses []models.QuestionResponse
	for _, question := range questions {
		responses = append(responses, models.QuestionResponse{
			ID:          question.ID,
			Title:       question.Title,
			Description: question.Description,
			Type:        question.Type,
			Category:    question.Category,
			Options:     question.Options,
			Score:       question.Score,
			OrderNum:    question.OrderNum,
			Status:      question.Status,
			CreatedAt:   question.CreatedAt,
			UpdatedAt:   question.UpdatedAt,
		})
	}

	response.Success(c, responses)
}

// GetQuestionByID 根据ID获取问题
func (h *QuestionnaireHandler) GetQuestionByID(c *gin.Context) {
	questionIDStr := c.Param("id")

	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的问题ID")
		return
	}

	var question models.Question
	if err := h.db.First(&question, questionID).Error; err != nil {
		response.NotFound(c, "问题不存在")
		return
	}

	response.Success(c, models.QuestionResponse{
		ID:          question.ID,
		Title:       question.Title,
		Description: question.Description,
		Type:        question.Type,
		Category:    question.Category,
		Options:     question.Options,
		Score:       question.Score,
		OrderNum:    question.OrderNum,
		Status:      question.Status,
		CreatedAt:   question.CreatedAt,
		UpdatedAt:   question.UpdatedAt,
	})
}

// SubmitAnswers 提交答案
func (h *QuestionnaireHandler) SubmitAnswers(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 1. 绑定请求参数
	var req struct {
		Answers []models.AnswerRequest `json:"answers" binding:"required"` // 包含question_id和answer_value
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 2. 开启事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. 创建评估记录
	assessment := models.Assessment{
		UserID: userID,
		Status: 1, // 直接标记为已完成
	}
	if err := tx.Create(&assessment).Error; err != nil {
		tx.Rollback()
		response.InternalServerError(c, "创建评估记录失败")
		return
	}

	// 4. 处理答案并计算总分
	var totalScore int
	for _, answerReq := range req.Answers {
		// 获取问题详情
		var question models.Question
		if err := tx.Where("id = ?", answerReq.QuestionID).First(&question).Error; err != nil {
			tx.Rollback()
			response.NotFound(c, fmt.Sprintf("问题%d不存在", answerReq.QuestionID))
			return
		}

		// 计算当前问题的得分
		score := calculateQuestionScore(question.Score, answerReq.AnswerValue)

		// 保存答案
		answer := models.Answer{
			UserID:       userID,
			QuestionID:   answerReq.QuestionID,
			AssessmentID: assessment.ID, // 使用新创建的评估ID
			Content:      getAnswerText(question, answerReq.AnswerValue),
			Score:        score,
		}
		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			response.InternalServerError(c, "保存答案失败")
			return
		}

		totalScore += score
	}

	// 5. 计算评估结果
	result := h.calculateAssessmentResult(totalScore)

	// 6. 更新评估记录（总分和结果）
	assessmentUpdates := map[string]interface{}{
		"total_score": totalScore,
		"level":       result.Level,
		"result":      result.Description,
	}
	if err := tx.Model(&assessment).Updates(assessmentUpdates).Error; err != nil {
		tx.Rollback()
		response.InternalServerError(c, "更新评估结果失败")
		return
	}

	// 7. 提交事务
	if err := tx.Commit().Error; err != nil {
		response.InternalServerError(c, "提交失败")
		return
	}

	// 8. 返回完整结果（包含评估ID和详情）
	response.Success(c, gin.H{
		"assessment_id": assessment.ID,
		"score":         totalScore,
		"level":         result.Level,
		"description":   result.Description,
		"suggestions":   result.Suggestions,
	})
}

// calculateAssessmentResult 计算评估结果
func (h *QuestionnaireHandler) calculateAssessmentResult(totalScore int) models.AssessmentResult {
	var level string
	var description string
	var suggestions string

	// 根据分数计算等级和建议
	if totalScore >= 80 {
		level = "severe"
		description = "您的抑郁倾向较为严重，建议立即寻求专业心理咨询师的帮助。"
		suggestions = "1. 尽快联系专业心理咨询师或精神科医生\n2. 保持规律的作息时间\n3. 多与家人朋友交流\n4. 避免独处时间过长"
	} else if totalScore >= 60 {
		level = "moderate"
		description = "您存在中等程度的抑郁倾向，建议适当调节心情并考虑寻求专业帮助。"
		suggestions = "1. 考虑寻求心理咨询师的帮助\n2. 增加户外活动和运动\n3. 培养兴趣爱好\n4. 保持社交活动"
	} else if totalScore >= 40 {
		level = "mild"
		description = "您存在轻微的抑郁倾向，属于正常范围，建议适当调节。"
		suggestions = "1. 多进行户外活动\n2. 保持规律作息\n3. 与朋友多交流\n4. 培养积极心态"
	} else {
		level = "normal"
		description = "您的心理状态良好，继续保持积极的生活态度。"
		suggestions = "1. 继续保持良好的生活习惯\n2. 定期进行心理健康检查\n3. 帮助身边的人保持心理健康"
	}

	return models.AssessmentResult{
		Level:       level,
		Score:       totalScore,
		MaxScore:    100, // 假设满分100
		Percentage:  float64(totalScore) / 100.0 * 100,
		Description: description,
		Suggestions: suggestions,
	}
}

// getAnswerText 获取选项对应的文本内容
func getAnswerText(question models.Question, answerValue int) string {
	// 解析问题的选项列表（假设options字段是JSON字符串）
	var options []string
	if err := json.Unmarshal([]byte(question.Options), &options); err != nil {
		return fmt.Sprintf("选项%d", answerValue) // 解析失败时返回默认值
	}

	// 检查选项值是否有效（数组索引从0开始）
	if answerValue > 0 && answerValue <= len(options) {
		return options[answerValue-1]
	}

	return fmt.Sprintf("无效选项%d", answerValue)
}

func calculateQuestionScore(questionWeight int, answerValue int) int {
	return answerValue * questionWeight
}
