package handlers

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"depression_go/inits"
	"depression_go/internal/models"
	"depression_go/middleware"
	"depression_go/pkg/response"
	"depression_go/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FaceDetectionHandler 人脸检测处理器
type FaceDetectionHandler struct {
	db           *gorm.DB
	baiduService *services.BaiduAIService
}

// NewFaceDetectionHandler 创建人脸检测处理器
func NewFaceDetectionHandler() *FaceDetectionHandler {
	return &FaceDetectionHandler{
		db:           inits.DB,
		baiduService: services.NewBaiduAIService(),
	}
}

// UploadImage 上传图片进行人脸检测
func (h *FaceDetectionHandler) UploadImage(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取上传的文件
	file, err := c.FormFile("image")
	fmt.Println(err)
	if err != nil {
		response.BadRequest(c, "请选择要上传的图片")
		return
	}

	// 验证文件
	if err := h.baiduService.ValidateImage(file.Filename, file.Size); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// 如果没有扩展名，默认使用.jpg
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%s_%s%s",
		time.Now().Format("20060102_150405"),
		uuid.New().String()[:8],
		ext,
	)

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.InternalServerError(c, "打开文件失败")
		return
	}
	defer src.Close()

	// 保存文件
	filePath, err := h.baiduService.SaveImage(src, filename)
	if err != nil {
		response.InternalServerError(c, "保存文件失败: "+err.Error())
		return
	}

	// 进行人脸检测和情绪分析
	emotionResult, err := h.baiduService.AnalyzeEmotion(filePath)
	if err != nil {
		response.InternalServerError(c, "人脸检测失败: "+err.Error())
		return
	}

	// 获取原始API响应数据
	aiResp, _ := h.baiduService.DetectFace(filePath)
	rawData := ""
	if aiResp != nil {
		// 这里可以序列化原始响应数据
		rawData = fmt.Sprintf("检测到%d个人脸", aiResp.Result.FaceNum)
	}

	// 创建人脸检测记录
	faceDetection := models.FaceDetection{
		UserID:     userID,
		ImagePath:  filePath,
		ImageURL:   fmt.Sprintf("/uploads/%s", filename),
		Emotion:    emotionResult.Emotion,
		Confidence: emotionResult.Confidence,
		Score:      emotionResult.Score,
		Level:      emotionResult.Level,
		Result:     emotionResult.Description,
		RawData:    rawData,
		Status:     1,
	}

	if err := h.db.Create(&faceDetection).Error; err != nil {
		response.InternalServerError(c, "保存检测记录失败")
		return
	}

	// 返回检测结果
	response.SuccessWithMessage(c, "人脸检测完成", models.FaceDetectionResponse{
		ID:         faceDetection.ID,
		UserID:     faceDetection.UserID,
		ImagePath:  faceDetection.ImagePath,
		ImageURL:   faceDetection.ImageURL,
		Emotion:    faceDetection.Emotion,
		Confidence: faceDetection.Confidence,
		Score:      faceDetection.Score,
		Level:      faceDetection.Level,
		Result:     faceDetection.Result,
		Status:     faceDetection.Status,
		CreatedAt:  faceDetection.CreatedAt,
		UpdatedAt:  faceDetection.UpdatedAt,
	})
}

// GetDetectionHistory 获取检测历史
func (h *FaceDetectionHandler) GetDetectionHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 查询检测记录
	var detections []models.FaceDetection
	var total int64

	offset := (page - 1) * pageSize

	// 统计总数
	if err := h.db.Model(&models.FaceDetection{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		response.InternalServerError(c, "查询失败")
		return
	}

	// 查询数据
	if err := h.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&detections).Error; err != nil {
		response.InternalServerError(c, "查询失败")
		return
	}

	// 转换为响应格式
	var responses []models.FaceDetectionResponse
	for _, detection := range detections {
		responses = append(responses, models.FaceDetectionResponse{
			ID:         detection.ID,
			UserID:     detection.UserID,
			ImagePath:  detection.ImagePath,
			ImageURL:   detection.ImageURL,
			Emotion:    detection.Emotion,
			Confidence: detection.Confidence,
			Score:      detection.Score,
			Level:      detection.Level,
			Result:     detection.Result,
			Status:     detection.Status,
			CreatedAt:  detection.CreatedAt,
			UpdatedAt:  detection.UpdatedAt,
		})
	}

	response.SuccessWithPage(c, responses, total, page, pageSize)
}
