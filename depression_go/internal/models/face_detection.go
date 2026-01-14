package models

import (
	"time"

	"gorm.io/gorm"
)

// FaceDetection 人脸检测模型
type FaceDetection struct {
	gorm.Model
	UserID     uint    `json:"user_id" gorm:"not null"`
	ImagePath  string  `json:"image_path" gorm:"size:500;not null"` // 图片路径
	ImageURL   string  `json:"image_url" gorm:"size:500"` // 图片URL
	Emotion    string  `json:"emotion" gorm:"size:50"` // 检测到的情绪：happy, sad, angry, fear, surprise, disgust, neutral
	Confidence float64 `json:"confidence" gorm:"type:decimal(5,4)"` // 置信度
	Score      int     `json:"score" gorm:"default:0"` // 情绪得分
	Level      string  `json:"level" gorm:"size:20"` // 情绪等级：normal, mild, moderate, severe
	Result     string  `json:"result" gorm:"type:text"` // 检测结果描述
	RawData    string  `json:"raw_data" gorm:"type:text"` // 原始API返回数据
	Status     int     `json:"status" gorm:"default:1"` // 1:成功 0:失败

	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (FaceDetection) TableName() string {
	return "face_detections"
}

// BeforeCreate 创建前的钩子函数
func (f *FaceDetection) BeforeCreate(tx *gorm.DB) error {
	if f.Status == 0 {
		f.Status = 1
	}
	return nil
}

// FaceDetectionCreateRequest 创建人脸检测请求
type FaceDetectionCreateRequest struct {
	ImagePath string `json:"image_path" binding:"required"`
	ImageURL  string `json:"image_url"`
}

// FaceDetectionUpdateRequest 更新人脸检测请求
type FaceDetectionUpdateRequest struct {
	Emotion    string  `json:"emotion"`
	Confidence float64 `json:"confidence"`
	Score      int     `json:"score"`
	Level      string  `json:"level"`
	Result     string  `json:"result"`
	RawData    string  `json:"raw_data"`
	Status     int     `json:"status"`
}

// FaceDetectionResponse 人脸检测响应
type FaceDetectionResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	ImagePath  string    `json:"image_path"`
	ImageURL   string    `json:"image_url"`
	Emotion    string    `json:"emotion"`
	Confidence float64   `json:"confidence"`
	Score      int       `json:"score"`
	Level      string    `json:"level"`
	Result     string    `json:"result"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// EmotionResult 情绪检测结果
type EmotionResult struct {
	Emotion     string  `json:"emotion"`
	Confidence  float64 `json:"confidence"`
	Score       int     `json:"score"`
	Level       string  `json:"level"`
	Description string  `json:"description"`
}

// BaiduAIResponse 百度AI接口响应
type BaiduAIResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	LogID     int64  `json:"log_id"`
	Timestamp int    `json:"timestamp"`
	Cached    int    `json:"cached"`
	Result    struct {
		FaceNum int `json:"face_num"`
		FaceList []struct {
			FaceToken string `json:"face_token"`
			Location  struct {
				Left   float64 `json:"left"`
				Top    float64 `json:"top"`
				Width  float64 `json:"width"`
				Height float64 `json:"height"`
				Rotation int    `json:"rotation"`
			} `json:"location"`
			FaceProbability float64 `json:"face_probability"`
			Angle           struct {
				Yaw   float64 `json:"yaw"`
				Pitch float64 `json:"pitch"`
				Roll  float64 `json:"roll"`
			} `json:"angle"`
			Age       int     `json:"age"`
			Beauty    float64 `json:"beauty"`
			Expression struct {
				Type        string  `json:"type"`
				Probability float64 `json:"probability"`
			} `json:"expression"`
			Emotion struct {
				Type        string  `json:"type"`
				Probability float64 `json:"probability"`
			} `json:"emotion"`
			FaceShape struct {
				Type        string  `json:"type"`
				Probability float64 `json:"probability"`
			} `json:"face_shape"`
			Landmark []struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"landmark"`
			Landmark72 []struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"landmark72"`
			Quality struct {
				Occlusion struct {
					LeftEye   float64 `json:"left_eye"`
					RightEye  float64 `json:"right_eye"`
					Nose      float64 `json:"nose"`
					Mouth     float64 `json:"mouth"`
					LeftCheek float64 `json:"left_cheek"`
					RightCheek float64 `json:"right_cheek"`
					Chin      float64 `json:"chin"`
				} `json:"occlusion"`
				Blur       float64 `json:"blur"`
				Illumination float64 `json:"illumination"`
				Completeness float64 `json:"completeness"`
			} `json:"quality"`
		} `json:"face_list"`
	} `json:"result"`
} 