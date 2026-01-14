# 抑郁倾向检测系统 API 文档

## 基础信息

- **基础URL**: `http://localhost:8088/api/v1`
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON
- **字符编码**: UTF-8

## 通用响应格式

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 422 | 验证错误 |
| 500 | 服务器内部错误 |

## 1. 认证相关接口

### 1.1 用户注册

**接口地址**: `POST /auth/register`

**请求参数**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "123456",
  "age": 25,
  "gender": "男",
  "phone": "13800138000"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "age": 25,
      "gender": "男",
      "phone": "13800138000",
      "avatar": "",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 1.2 用户登录

**接口地址**: `POST /auth/login`

**请求参数**:
```json
{
  "username": "testuser",
  "password": "123456"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "age": 25,
      "gender": "男",
      "phone": "13800138000",
      "avatar": "",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## 2. 用户相关接口（需要认证）

### 2.1 获取用户信息

**接口地址**: `GET /user/profile`

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "age": 25,
    "gender": "男",
    "phone": "13800138000",
    "avatar": "",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2.2 更新用户信息

**接口地址**: `PUT /user/profile`

**请求参数**:
```json
{
  "age": 26,
  "gender": "男",
  "phone": "13800138001",
  "avatar": "http://example.com/avatar.jpg"
}
```

### 2.3 修改密码

**接口地址**: `PUT /user/password`

**请求参数**:
```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

## 3. 问题相关接口

### 3.1 获取问题列表

**接口地址**: `GET /questions`

**查询参数**:
- `category`: 问题分类（depression, anxiety, stress）
- `status`: 状态（1:启用, 0:禁用）

**响应示例**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {
      "id": 1,
      "title": "您最近是否感到情绪低落或沮丧？",
      "description": "请根据最近两周的感受选择最符合的选项",
      "type": "single",
      "category": "depression",
      "options": "[\"从不\", \"偶尔\", \"经常\", \"总是\"]",
      "score": 10,
      "order_num": 1,
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 3.2 获取问题详情

**接口地址**: `GET /questions/{id}`

## 4. 人脸检测相关接口（需要认证）

### 4.1 上传图片进行人脸检测

**接口地址**: `POST /face/upload`

**请求类型**: `multipart/form-data`

**请求参数**:
- `image`: 图片文件

**响应示例**:
```json
{
  "code": 200,
  "message": "人脸检测完成",
  "data": {
    "id": 1,
    "user_id": 1,
    "image_path": "./uploads/20240101_120000_abc123.jpg",
    "image_url": "/uploads/20240101_120000_abc123.jpg",
    "emotion": "sad",
    "confidence": 0.85,
    "score": 85,
    "level": "moderate",
    "result": "检测到中等程度的悲伤情绪，建议适当调节心情",
    "status": 1,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### 4.2 获取检测历史

**接口地址**: `GET /face/history`

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**响应示例**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "image_path": "./uploads/20240101_120000_abc123.jpg",
        "image_url": "/uploads/20240101_120000_abc123.jpg",
        "emotion": "sad",
        "confidence": 0.85,
        "score": 85,
        "level": "moderate",
        "result": "检测到中等程度的悲伤情绪，建议适当调节心情",
        "status": 1,
        "created_at": "2024-01-01T12:00:00Z",
        "updated_at": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

### 4.3 获取检测统计

**接口地址**: `GET /face/stats`

**响应示例**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "total_detections": 10,
    "emotion_stats": {
      "sad": 5,
      "happy": 3,
      "neutral": 2
    },
    "level_stats": {
      "moderate": 5,
      "normal": 3,
      "mild": 2
    },
    "recent_detections": []
  }
}
```

## 5. 问卷相关接口（需要认证）

### 5.1 提交答案

**接口地址**: `POST /questionnaire/submit`

**请求参数**:
```json
{
  "assessment_id": 1,
  "answers": [
    {
      "question_id": 1,
      "content": "经常",
      "score": 3
    },
    {
      "question_id": 2,
      "content": "比较明显",
      "score": 3
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "提交成功",
  "data": {
    "level": "moderate",
    "score": 60,
    "max_score": 100,
    "percentage": 60.0,
    "description": "您存在中等程度的抑郁倾向，建议适当调节心情并考虑寻求专业帮助。",
    "suggestions": "1. 考虑寻求心理咨询师的帮助\n2. 增加户外活动和运动\n3. 培养兴趣爱好\n4. 保持社交活动"
  }
}
```

## 6. 评估结果相关接口（需要认证）

### 6.1 创建评估

**接口地址**: `POST /assessment`

**请求参数**:
```json
{
  "title": "抑郁倾向评估",
  "type": "questionnaire"
}
```

### 6.2 获取评估历史

**接口地址**: `GET /assessment/history`

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

### 6.3 获取评估详情

**接口地址**: `GET /assessment/history/{id}`

**响应示例**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "抑郁倾向评估",
    "type": "questionnaire",
    "total_score": 60,
    "max_score": 100,
    "level": "moderate",
    "result": "您存在中等程度的抑郁倾向，建议适当调节心情并考虑寻求专业帮助。",
    "status": 1,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z",
    "answers": [
      {
        "id": 1,
        "user_id": 1,
        "question_id": 1,
        "assessment_id": 1,
        "content": "经常",
        "score": 3,
        "created_at": "2024-01-01T12:00:00Z",
        "updated_at": "2024-01-01T12:00:00Z",
        "question": {
          "id": 1,
          "title": "您最近是否感到情绪低落或沮丧？",
          "description": "请根据最近两周的感受选择最符合的选项",
          "type": "single",
          "category": "depression",
          "options": "[\"从不\", \"偶尔\", \"经常\", \"总是\"]",
          "score": 10,
          "order_num": 1,
          "status": 1,
          "created_at": "2024-01-01T00:00:00Z",
          "updated_at": "2024-01-01T00:00:00Z"
        }
      }
    ]
  }
}
```

### 6.4 获取评估统计

**接口地址**: `GET /assessment/stats`

**响应示例**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "total_assessments": 5,
    "completed_assessments": 4,
    "level_stats": {
      "moderate": 2,
      "normal": 1,
      "mild": 1
    },
    "type_stats": {
      "questionnaire": 4,
      "face": 1
    },
    "recent_assessments": [],
    "average_score": 45.5
  }
}
```

### 6.5 获取综合评估结果

**接口地址**: `GET /assessment/combined`

**响应示例**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "combined_score": 65,
    "combined_level": "moderate",
    "description": "综合评估显示您存在中等程度的抑郁倾向，建议适当调节并考虑寻求专业帮助。",
    "suggestions": "1. 考虑寻求心理咨询师的帮助\n2. 增加户外活动和运动\n3. 培养兴趣爱好\n4. 保持社交活动\n5. 学习放松技巧",
    "questionnaire": {
      "score": 60,
      "level": "moderate"
    },
    "face_detection": {
      "score": 75,
      "level": "moderate",
      "emotion": "sad"
    },
    "assessment_date": "2024-01-01T12:00:00Z",
    "detection_date": "2024-01-01T12:00:00Z"
  }
}
```

## 7. 管理员接口（需要认证）

### 7.1 获取用户信息

**接口地址**: `GET /admin/users/{id}`

### 7.2 创建问题

**接口地址**: `POST /admin/questions`

**请求参数**:
```json
{
  "title": "新问题标题",
  "description": "问题描述",
  "type": "single",
  "category": "depression",
  "options": "[\"选项1\", \"选项2\", \"选项3\", \"选项4\"]",
  "score": 10,
  "order_num": 11
}
```

### 7.3 更新问题

**接口地址**: `PUT /admin/questions/{id}`

### 7.4 删除问题

**接口地址**: `DELETE /admin/questions/{id}`

## 8. 其他接口

### 8.1 健康检查

**接口地址**: `GET /health`

**响应示例**:
```json
{
  "status": "ok",
  "message": "服务运行正常"
}
```

## 使用说明

1. **认证流程**:
   - 用户注册或登录后获取JWT令牌
   - 在后续请求的Header中添加：`Authorization: Bearer <token>`

2. **文件上传**:
   - 支持格式：jpg, jpeg, png, bmp, gif
   - 最大文件大小：10MB

3. **分页查询**:
   - 默认页码：1
   - 默认每页数量：10
   - 最大每页数量：100

4. **错误处理**:
   - 所有接口都会返回统一的错误格式
   - 详细错误信息在message字段中

5. **数据验证**:
   - 用户名：3-50个字符
   - 密码：至少6个字符
   - 邮箱：标准邮箱格式
   - 年龄：0-150
   - 性别：男/女/未知 