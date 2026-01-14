# 抑郁倾向检测系统 - Go后端

基于Gin和Gorm框架的抑郁倾向检测系统后端API服务。

## 项目特性

- **Gorm.Model集成**: 所有模型都使用`gorm.Model`，自动包含ID、CreatedAt、UpdatedAt和DeletedAt字段
- **软删除支持**: 通过`gorm.Model`的`DeletedAt`字段实现软删除
- **自动时间戳**: 自动管理创建和更新时间
- **RESTful API**: 完整的RESTful接口设计
- **JWT认证**: 基于JWT的用户认证系统
- **人脸识别**: 集成百度云人脸识别API
- **问卷评估**: 完整的抑郁倾向问卷评估系统
- **CORS支持**: 跨域资源共享支持

## 技术栈

- **框架**: Gin (Web框架)
- **ORM**: Gorm (数据库ORM)
- **数据库**: MySQL
- **认证**: JWT
- **配置**: 环境变量 + godotenv
- **API**: 百度云人脸识别API

## 项目结构

```
depression_go/
├── config.env.example      # 环境变量示例
├── go.mod                  # Go模块文件
├── go.sum                  # 依赖校验文件
├── main.go                 # 主程序入口
├── database/
│   └── init.sql           # 数据库初始化脚本
├── inits/
│   └── database.go        # 数据库初始化
├── internal/
│   ├── models/            # 数据模型 (使用gorm.Model)
│   │   ├── user.go
│   │   ├── question.go
│   │   ├── answer.go
│   │   ├── assessment.go
│   │   └── face_detection.go
│   ├── controllers/       # 控制器层
│   │   ├── auth_controller.go
│   │   ├── questionnaire_controller.go
│   │   └── face_detection_controller.go
│   └── services/          # 服务层
│       ├── auth_service.go
│       ├── questionnaire_service.go
│       ├── face_detection_service.go
│       └── baidu_ai_service.go
├── pkg/
│   ├── middleware/        # 中间件
│   │   ├── auth.go
│   │   └── cors.go
│   ├── response/          # 响应工具
│   │   └── response.go
│   ├── token/             # JWT工具
│   │   └── jwt.go
│   └── utils/             # 工具函数
│       └── password.go
└── routers/
    └── router.go          # 路由配置
```

## Gorm.Model 使用说明

### 模型定义

所有模型都使用`gorm.Model`，它包含以下字段：

```go
type Model struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### 示例模型

```go
type User struct {
    gorm.Model              // 包含ID、CreatedAt、UpdatedAt、DeletedAt
    Username string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
    Email    string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
    Password string         `json:"-" gorm:"size:255;not null"`
    // ... 其他字段
}
```

### 软删除

使用`gorm.Model`后，删除操作默认是软删除：

```go
// 软删除 - 只设置DeletedAt字段
db.Delete(&user)

// 硬删除 - 真正从数据库删除
db.Unscoped().Delete(&user)

// 查询时自动过滤已软删除的记录
db.Find(&users) // 只查询未删除的记录

// 查询包含已删除的记录
db.Unscoped().Find(&users)
```

### 自动时间戳

- `CreatedAt`: 创建记录时自动设置
- `UpdatedAt`: 更新记录时自动更新
- `DeletedAt`: 软删除时自动设置

## 环境配置

复制`config.env.example`为`config.env`并配置：

```env
# 服务器配置
PORT=8088

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=depression_ai

# JWT配置
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE_HOURS=24

# 百度云AI配置
BAIDU_APP_ID=your_app_id
BAIDU_API_KEY=your_api_key
BAIDU_SECRET_KEY=your_secret_key
```

## 安装和运行

1. **安装依赖**
   ```bash
   go mod tidy
   ```

2. **配置环境变量**
   ```bash
   cp config.env.example config.env
   # 编辑config.env文件
   ```

3. **初始化数据库**
   ```bash
   mysql -u root -p < database/init.sql
   ```

4. **运行服务**
   ```bash
   go run main.go
   ```

## API接口

### 用户认证
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/profile` - 获取用户信息
- `PUT /api/auth/profile` - 更新用户信息

### 问卷评估
- `GET /api/questions` - 获取问题列表
- `POST /api/assessments` - 创建评估
- `POST /api/assessments/:id/answers` - 提交答案
- `GET /api/assessments/:id/result` - 获取评估结果
- `GET /api/assessments` - 获取用户评估历史

### 人脸检测
- `POST /api/face/detect` - 人脸情绪检测
- `GET /api/face/history` - 获取检测历史

## 数据库设计

### 表结构 (使用gorm.Model)

所有表都包含以下基础字段：
- `id` (BIGINT UNSIGNED, 主键)
- `created_at` (TIMESTAMP, 创建时间)
- `updated_at` (TIMESTAMP, 更新时间)
- `deleted_at` (TIMESTAMP NULL, 软删除时间)

### 主要表

1. **users** - 用户表
2. **questions** - 问题表
3. **assessments** - 评估表
4. **answers** - 答案表
5. **face_detections** - 人脸检测表

## 开发说明

### 添加新模型

1. 在`internal/models/`目录创建新模型文件
2. 使用`gorm.Model`作为基础结构
3. 在`inits/database.go`中添加模型到AutoMigrate
4. 创建对应的控制器和服务

### 数据库操作示例

```go
// 查询
var users []models.User
db.Find(&users)

// 创建
user := models.User{Username: "test", Email: "test@example.com"}
db.Create(&user)

// 更新
db.Model(&user).Update("status", 0)

// 软删除
db.Delete(&user)

// 关联查询
var assessment models.Assessment
db.Preload("Answers").Preload("User").First(&assessment, id)
```

## 许可证

MIT License

## 联系方式

如有问题或建议，请通过以下方式联系：

- 项目Issues: [GitHub Issues](https://github.com/your-repo/depression_ai/issues)
- 邮箱: your-email@example.com

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 支持用户认证
- 支持问卷评估
- 支持人脸检测
- 支持综合评估 