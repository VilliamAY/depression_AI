package handlers

import (
	"depression_go/inits"
	"depression_go/internal/models"
	"depression_go/middleware"
	"depression_go/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ResultHandler 评估结果处理器
type ResultHandler struct {
	db *gorm.DB
}

// NewResultHandler 创建评估结果处理器
func NewResultHandler() *ResultHandler {
	return &ResultHandler{
		db: inits.DB,
	}
}

// CreateAssessment 创建评估
func (h *ResultHandler) CreateAssessment(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.AssessmentCreateRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 创建评估
	assessment := models.Assessment{
		UserID: userID,
		Title:  req.Title,
		Type:   req.Type,
		Status: 0, // 进行中
	}

	if err := h.db.Create(&assessment).Error; err != nil {
		response.InternalServerError(c, "创建评估失败")
		return
	}

	response.SuccessWithMessage(c, "评估创建成功", models.AssessmentResponse{
		ID:         assessment.ID,
		UserID:     assessment.UserID,
		Title:      assessment.Title,
		Type:       assessment.Type,
		TotalScore: assessment.TotalScore,
		MaxScore:   assessment.MaxScore,
		Level:      assessment.Level,
		Result:     assessment.Result,
		Status:     assessment.Status,
		CreatedAt:  assessment.CreatedAt,
		UpdatedAt:  assessment.UpdatedAt,
	})
}

// GetCombinedResult 获取综合评估结果（结合问卷和人脸检测）
func (h *ResultHandler) GetCombinedResult(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取最近的问卷评估
	var questionnaireAssessment models.Assessment
	if err := h.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&questionnaireAssessment).Error; err != nil {
		response.NotFound(c, "未找到评估记录记录")
		return
	}

	// 获取最近的人脸检测
	var faceDetection models.FaceDetection
	if err := h.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&faceDetection).Error; err != nil {
		response.NotFound(c, "未找到人脸检测记录")
		return
	}

	// 计算综合得分
	questionnaireScore := questionnaireAssessment.TotalScore
	faceScore := faceDetection.Score

	// 权重分配：问卷70%，人脸检测30%
	combinedScore := int(float64(questionnaireScore)*0.7 + float64(faceScore)*0.3)

	// 计算综合等级
	var combinedLevel string
	var description string
	var suggestions string

	if combinedScore >= 80 {
		combinedLevel = "severe"
		description = "综合评估显示您的抑郁倾向较为严重，建议立即寻求专业帮助。"
		suggestions = "1. 立即联系专业心理咨询师或精神科医生\n2. 保持规律的作息时间\n3. 多与家人朋友交流\n4. 避免独处时间过长\n5. 考虑药物治疗"
	} else if combinedScore >= 60 {
		combinedLevel = "moderate"
		description = "综合评估显示您存在中等程度的抑郁倾向，建议适当调节并考虑寻求专业帮助。"
		suggestions = "1. 考虑寻求心理咨询师的帮助\n2. 增加户外活动和运动\n3. 培养兴趣爱好\n4. 保持社交活动\n5. 学习放松技巧"
	} else if combinedScore >= 40 {
		combinedLevel = "mild"
		description = "综合评估显示您存在轻微的抑郁倾向，属于正常范围，建议适当调节。"
		suggestions = "1. 多进行户外活动\n2. 保持规律作息\n3. 与朋友多交流\n4. 培养积极心态\n5. 学习压力管理"
	} else {
		combinedLevel = "normal"
		description = "综合评估显示您的心理状态良好，继续保持积极的生活态度。"
		suggestions = "1. 继续保持良好的生活习惯\n2. 定期进行心理健康检查\n3. 帮助身边的人保持心理健康\n4. 培养兴趣爱好"
	}

	result := gin.H{
		"combined_score": combinedScore,
		"combined_level": combinedLevel,
		"description":    description,
		"suggestions":    suggestions,
		"questionnaire": gin.H{
			"score": questionnaireAssessment.TotalScore,
			"level": questionnaireAssessment.Level,
		},
		"face_detection": gin.H{
			"score":   faceDetection.Score,
			"level":   faceDetection.Level,
			"emotion": faceDetection.Emotion,
		},
		"assessment_date": questionnaireAssessment.CreatedAt,
		"detection_date":  faceDetection.CreatedAt,
	}

	response.Success(c, result)
}
