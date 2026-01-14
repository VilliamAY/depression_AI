package routes

import (
	"depression_go/handlers"
	"depression_go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 创建处理器实例
	authHandler := handlers.NewAuthHandler()
	faceDetectionHandler := handlers.NewFaceDetectionHandler()
	questionnaireHandler := handlers.NewQuestionnaireHandler()
	resultHandler := handlers.NewResultHandler()

	// API版本组
	api := r.Group("/api/v1")

	// 公开路由（无需认证）
	public := api.Group("")
	{
		// 认证相关
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 问题相关（公开访问）
		questions := public.Group("/questions")
		{
			questions.GET("", questionnaireHandler.GetQuestions)
			questions.GET("/:id", questionnaireHandler.GetQuestionByID)
		}
	}

	// 需要认证的路由
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		user := protected.Group("/user")
		{
			//获取用户信息
			user.GET("/profile", authHandler.GetProfile)
		}

		// 人脸检测相关
		face := protected.Group("/face")
		{
			//上传图片进行人脸检测，返回结果
			face.POST("/upload", faceDetectionHandler.UploadImage)
			//获取检测历史
			face.GET("/history", faceDetectionHandler.GetDetectionHistory)
		}

		// 问卷相关
		questionnaire := protected.Group("/questionnaire")
		{
			//提交答案
			questionnaire.POST("/submit", questionnaireHandler.SubmitAnswers)
		}

		// 评估结果相关
		assessment := protected.Group("/assessment")
		{
			//创建评估
			assessment.POST("", resultHandler.CreateAssessment)
			assessment.GET("/total", resultHandler.GetCombinedResult)
		}
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "接口不存在",
		})
	})
}
