package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"depression_go/internal/models"
)

// BaiduAIService 百度AI服务
type BaiduAIService struct {
	AppID     string
	APIKey    string
	SecretKey string
	client    *http.Client
}

// NewBaiduAIService 创建百度AI服务实例
func NewBaiduAIService() *BaiduAIService {
	appID := os.Getenv("BAIDU_APP_ID")
	apiKey := os.Getenv("BAIDU_API_KEY")
	secretKey := os.Getenv("BAIDU_SECRET_KEY")
	if appID == "" || apiKey == "" || secretKey == "" {
		panic("百度AI配置缺失，请检查环境变量 BAIDU_APP_ID、BAIDU_API_KEY、BAIDU_SECRET_KEY 是否设置")
	}
	return &BaiduAIService{
		AppID:     appID,
		APIKey:    apiKey,
		SecretKey: secretKey,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// GetAccessToken 获取百度AI访问令牌 AccessToken
func (s *BaiduAIService) GetAccessToken() (string, error) {
	requestURL := "https://aip.baidubce.com/oauth/2.0/token"
	params := url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", s.APIKey)
	params.Set("client_secret", s.SecretKey)
	resp, err := s.client.PostForm(requestURL, params)
	if err != nil {
		return "", fmt.Errorf("获取访问令牌失败: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	var tokenResp struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	if tokenResp.Error != "" {
		return "", fmt.Errorf("获取访问令牌错误: %s - %s", tokenResp.Error, tokenResp.ErrorDescription)
	}
	return tokenResp.AccessToken, nil
}

// DetectFace 人脸检测
func (s *BaiduAIService) DetectFace(imagePath string) (*models.BaiduAIResponse, error) {
	accessToken, err := s.GetAccessToken()
	if err != nil {
		return nil, err
	}
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("打开图片文件失败: %v", err)
	}
	defer imageFile.Close()
	imageBytes, err := io.ReadAll(imageFile)
	if err != nil {
		return nil, fmt.Errorf("读取图片内容失败: %v", err)
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
	requestURL := fmt.Sprintf("https://aip.baidubce.com/rest/2.0/face/v3/detect?access_token=%s", accessToken)
	params := url.Values{}
	params.Set("image", imageBase64)
	params.Set("image_type", "BASE64")
	params.Set("face_field", "age,beauty,expression,emotion,face_shape,landmark,landmark72,quality")
	params.Set("max_face_num", "1")
	resp, err := s.client.PostForm(requestURL, params)
	if err != nil {
		return nil, fmt.Errorf("人脸检测请求失败: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	var aiResp models.BaiduAIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	if aiResp.ErrorCode != 0 {
		return nil, fmt.Errorf("人脸检测失败: %s (错误码: %d)", aiResp.ErrorMsg, aiResp.ErrorCode)
	}

	return &aiResp, nil
}

// AnalyzeEmotion 分析情绪
func (s *BaiduAIService) AnalyzeEmotion(imagePath string) (*models.EmotionResult, error) {
	aiResp, err := s.DetectFace(imagePath)
	if err != nil {
		return nil, err
	}
	if aiResp.Result.FaceNum == 0 {
		return nil, fmt.Errorf("未检测到人脸")
	}
	face := aiResp.Result.FaceList[0]
	emotion := face.Emotion.Type
	confidence := face.Emotion.Probability
	score, level, description := s.calculateEmotionScore(emotion, confidence)
	return &models.EmotionResult{
		Emotion:     emotion,
		Confidence:  confidence,
		Score:       score,
		Level:       level,
		Description: description,
	}, nil
}

// calculateEmotionScore 计算情绪得分和等级
func (s *BaiduAIService) calculateEmotionScore(emotion string, confidence float64) (int, string, string) {
	var score int
	var level string
	var description string
	switch emotion {
	case "sad":
		score = int(confidence * 100)
		if score >= 80 {
			level = "severe"
			description = "检测到明显的悲伤情绪，建议寻求专业心理咨询"
		} else if score >= 60 {
			level = "moderate"
			description = "检测到中等程度的悲伤情绪，建议适当调节心情"
		} else if score >= 40 {
			level = "mild"
			description = "检测到轻微的悲伤情绪，属于正常范围"
		} else {
			level = "normal"
			description = "情绪状态正常"
		}
	case "angry":
		score = int(confidence * 90)
		if score >= 70 {
			level = "moderate"
			description = "检测到愤怒情绪，建议冷静处理"
		} else {
			level = "mild"
			description = "检测到轻微愤怒，属于正常情绪波动"
		}
	case "fear":
		score = int(confidence * 85)
		if score >= 70 {
			level = "moderate"
			description = "检测到恐惧情绪，建议寻求支持"
		} else {
			level = "mild"
			description = "检测到轻微恐惧，属于正常反应"
		}
	case "disgust":
		score = int(confidence * 80)
		level = "mild"
		description = "检测到厌恶情绪，属于正常反应"
	case "surprise":
		score = int(confidence * 60)
		level = "normal"
		description = "检测到惊讶情绪，属于正常反应"
	case "happy":
		score = int(confidence * 50)
		level = "normal"
		description = "检测到快乐情绪，情绪状态良好"
	default: // neutral
		score = int(confidence * 30)
		level = "normal"
		description = "情绪状态平静，属于正常范围"
	}
	return score, level, description
}

// SaveImage 保存上传的图片
func (s *BaiduAIService) SaveImage(file io.Reader, filename string) (string, error) {
	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("创建上传目录失败: %v", err)
	}
	filePath := fmt.Sprintf("%s/%s", uploadPath, filename)
	fileHandle, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer fileHandle.Close()
	_, err = io.Copy(fileHandle, file)
	if err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}
	return filePath, nil
}

// ValidateImage 验证图片格式和大小
func (s *BaiduAIService) ValidateImage(filename string, size int64) error {
	maxSize := int64(10485760) // 默认10MB
	if envMax := os.Getenv("MAX_FILE_SIZE"); envMax != "" {
		fmt.Sscanf(envMax, "%d", &maxSize)
	}
	if size > maxSize {
		return fmt.Errorf("文件大小超过限制，最大允许 %d 字节", maxSize)
	}

	// 如果文件名为空或者是从摄像头捕获的图像（如camera_capture.jpg），则跳过扩展名检查
	if filename == "" || strings.Contains(strings.ToLower(filename), "camera_capture") {
		return nil
	}

	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"}
	ext := filepath.Ext(strings.ToLower(filename))
	valid := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("不支持的文件格式，仅支持: %v", allowedExtensions)
	}
	return nil
}
