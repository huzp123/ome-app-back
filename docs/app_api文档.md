# OME-APP 前端API接口文档

## 通用说明

*有任何相关变化都需要同步到该api文档中*

### 基础URL
```
/api/v1
```

### 认证方式
除了登录、注册、微信登录和健康检查API外，所有请求都需要在Header中加入认证信息：
```
Authorization: Bearer {token}
```

### 响应格式
所有API响应都遵循以下格式：
```json
{
  "code": 0,       // 0表示成功，非0表示错误
  "msg": "成功",    // 错误信息
  "data": {},      // 响应数据，错误时可能为null
  "details": []    // 详细错误信息，仅在错误时有值
}
```

## 公开接口

### 健康检查

**请求**
```
GET /ping
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "message": "pong"
  }
}
```

### 用户注册

**请求**
```
POST /register
```

**请求参数**
```json
{
  "phone": "string",     // 手机号
  "email": "string",     // 邮箱
  "user_name": "string", // 用户名
  "password": "string"   // 密码
}
```

**说明**
- 用户注册时不需要提供身高等个人信息，这些信息将在个人资料更新时填写
- 手机号和邮箱至少需要提供一个，两者都可以为空，但不能同时为空
- 手机号和邮箱字段都有唯一性约束，不能使用已注册的手机号或邮箱

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "user_id": 1,
    "token": "string"
  }
}
```

### 用户登录

**请求**
```
POST /login
```

**请求参数**
```json
{
  "account": "string",  // 账号，可以是手机号或邮箱
  "password": "string"  // 密码
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "user_id": 1,
    "user_name": "string",
    "token": "string",
    "is_profile_complete": true  // 用户档案是否已完善
  }
}
```

### 微信登录

**请求**
```
POST /wechat/login
```

**请求参数**
```json
{
  "openid": "string",        // 微信OpenID（必填）
  "user_name": "string",     // 用户昵称（可选）
  "avatar_url": "string"     // 头像URL（可选）
}
```

**说明**
- 微信登录支持新用户自动注册和已有用户登录
- 如果用户不存在，系统会自动创建新用户
- 如果用户已存在，系统会更新用户信息
- 微信登录的用户password_hash字段为空，只能通过微信登录

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "user_id": 123,
    "user_name": "用户昵称",
    "token": "JWT令牌",
    "is_new_user": true,           // 是否为新用户
    "is_profile_complete": false   // 用户档案是否已完善
  }
}
```

## 用户相关接口（需要认证）

### 获取用户信息

**请求**
```
GET /user/info
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "user_name": "张三",
    "phone": "13800138000",
    "email": "zhangsan@example.com",
    "wechat_openid": "微信OpenID",    // 微信登录用户的OpenID
    "avatar_url": "头像URL",          // 用户头像URL
    "birth_date": "1990-01-01",
    "sex": "male",
    "created_at": "2023-04-01T12:00:00Z",
    "updated_at": "2023-04-15T10:30:00Z",
    "is_profile_complete": true
  }
}
```

**说明**
- `wechat_openid` 和 `avatar_url` 字段仅对微信登录用户有值
- 普通注册用户这些字段为 `null`

### 更新用户档案

**请求**
```
PUT /user/profile
```

**请求参数**
```json
{
  "phone": "string",            // 手机号
  "email": "string",            // 邮箱
  "birth_date": "2000-01-01",   // 出生日期，格式YYYY-MM-DD
  "sex": "male",                // 性别: male/female/other
  "weight_kg": 70.5             // 当前体重(公斤)
}
```

**说明**
- 所有字段均为可选，只需填写需要更新的信息
- 更新手机号或邮箱时会检查唯一性，不能使用已被其他用户注册的联系方式

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 更新用户健康目标

**请求**
```
PUT /user/goal
```

**请求参数**
```json
{
  "goal_type": "lose_fat",              // 目标类型: lose_fat/keep_fit/gain_muscle，必填
  "target_weight_kg": 65.0,             // 目标体重(公斤)，必填
  "weekly_change_kg": 0.5,              // 每周计划变化的体重(公斤)，必填
  "target_date": "2023-12-31",          // 目标日期，格式YYYY-MM-DD，必填
  "diet_type": "normal",                // 饮食类型: normal/vegetarian/meat_lover，必填
  "taste_preferences": ["清淡", "酸的"], // 口味偏好，必填，至少选择1个
  "food_intolerances": ["海鲜"]         // 食物不耐受/禁忌，必填，至少选择1个
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取用户健康目标

**请求**
```
GET /user/goal
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "goal_type": "lose_fat",              // 目标类型: lose_fat/keep_fit/gain_muscle
    "target_weight_kg": 65.0,             // 目标体重(公斤)
    "weekly_change_kg": 0.5,              // 每周计划变化的体重(公斤)
    "target_date": "2023-12-31",          // 目标日期，格式YYYY-MM-DD
    "diet_type": "normal",                // 饮食类型: normal/vegetarian/meat_lover
    "taste_preferences": ["清淡", "酸的"], // 口味偏好
    "food_intolerances": ["海鲜"],         // 食物不耐受/禁忌
    "created_at": "2023-04-01T12:00:00Z"  // 创建时间
  }
}
```

**说明**
- 如果用户还没有设置健康目标，`data`字段将为`null`
- 首次使用的用户需要先通过更新健康目标接口设置目标后才能获取到数据

## 身高管理相关接口（需要认证）

### 记录身高

**请求**
```
POST /user/height
```

**请求参数**
```json
{
  "height_cm": 175.0  // 身高(厘米)，范围50-300cm
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取身高历史记录

**请求**
```
GET /user/height/history?limit=30
```

**查询参数**
- limit: 可选，限制返回的记录数量，默认为30，最多获取一年内的数据

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "height_cm": 175.0,
      "record_date": "2023-05-01T00:00:00Z",
      "created_at": "2023-05-01T10:30:00Z"
    },
    {
      "id": 2,
      "height_cm": 174.8,
      "record_date": "2023-04-30T00:00:00Z",
      "created_at": "2023-04-30T09:15:00Z"
    }
  ]
}
```

### 获取当前身高信息

**请求**
```
GET /user/height/current
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "height_cm": 175.0,
    "record_date": "2023-05-01T00:00:00Z",
    "days_ago": 2  // 距离现在多少天前记录的
  }
}
```

### 删除身高记录

**请求**
```
DELETE /user/height/{id}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

**说明**
- 只能删除自己的身高记录
- 删除不存在的记录会返回错误

### 获取身高统计分析

**请求**
```
GET /user/height/statistics?days=30
```

**查询参数**
- days: 可选，统计天数，默认为30天

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "current_height": 175.0,    // 当前身高
    "min_height": 174.5,        // 最低身高
    "max_height": 175.2,        // 最高身高
    "avg_height": 174.8,        // 平均身高
    "height_change": 0.5,       // 身高变化(正数为增加，负数为减少)
    "trend_data": [             // 趋势数据点，用于绘制图表
      {
        "date": "2023-04-01T00:00:00Z",
        "height_cm": 174.5
      },
      {
        "date": "2023-04-15T00:00:00Z",
        "height_cm": 174.8
      },
      {
        "date": "2023-05-01T00:00:00Z",
        "height_cm": 175.0
      }
    ]
  }
}
```

**说明**
- 统计数据基于指定天数内的身高记录
- trend_data 按时间顺序排列，可用于绘制身高变化曲线
- 如果没有身高记录会返回错误

## 健康分析相关接口（需要认证）

### 生成健康分析报告

**请求**
```
GET /health/analysis
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "bmi": 23.5,                     // BMI指数
    "bmi_category": "正常",           // BMI分类
    "bmr": 1550.0,                   // 基础代谢率(千卡)
    "tdee": 2130.0,                  // 每日总能量消耗(千卡)
    "recommended_calories": 1880.0,   // 推荐每日摄入热量(千卡)
    "protein_need_g": 140.0,          // 蛋白质需求(克)
    "carb_need_g": 210.0,             // 碳水需求(克)
    "fat_need_g": 60.0,               // 脂肪需求(克)
    "analysis_content": "string",     // 分析结果文本内容
    "current_weight_kg": 70.5,        // 当前体重(公斤)
    "target_weight_kg": 65.0,         // 目标体重(公斤)
    "weekly_change_kg": 0.5,          // 每周计划变化的体重(公斤)
    "target_date": "2023-12-31",      // 目标日期
    "days_to_target": 120             // 距离目标日期天数
  }
}
```

### 获取健康分析历史记录

**请求**
```
GET /health/history?limit=10
```

**查询参数**
- limit: 可选，限制返回的记录数量，默认为10

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "bmi": 23.5,
      "bmr": 1550.0,
      "tdee": 2130.0,
      "protein_need_g": 140.0,
      "carb_need_g": 210.0,
      "fat_need_g": 60.0,
      "recommended_calories": 1880.0,
      "analysis_content": "string",
      "created_at": "2023-04-01T12:00:00Z"
    }
  ]
}
```

## 体重管理相关接口（需要认证）

### 手动记录体重

**请求**
```
POST /user/weight
```

**请求参数**
```json
{
  "weight_kg": 85.5  // 体重(公斤)，范围20-500kg
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取体重历史记录

**请求**
```
GET /user/weight/history?limit=30
```

**查询参数**
- limit: 可选，限制返回的记录数量，默认为30，最多获取一年内的数据

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "weight_kg": 85.5,
      "record_date": "2023-05-01T00:00:00Z",
      "created_at": "2023-05-01T10:30:00Z"
    },
    {
      "id": 2,
      "weight_kg": 85.2,
      "record_date": "2023-04-30T00:00:00Z",
      "created_at": "2023-04-30T09:15:00Z"
    }
  ]
}
```

### 获取当前体重信息

**请求**
```
GET /user/weight/current
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "weight_kg": 85.5,
    "record_date": "2023-05-01T00:00:00Z",
    "days_ago": 2  // 距离现在多少天前记录的
  }
}
```

### 删除体重记录

**请求**
```
DELETE /user/weight/{id}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

**说明**
- 只能删除自己的体重记录
- 删除不存在的记录会返回错误

### 获取体重统计分析

**请求**
```
GET /user/weight/statistics?days=30
```

**查询参数**
- days: 可选，统计天数，默认为30天

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "current_weight": 85.5,    // 当前体重
    "min_weight": 84.8,        // 最低体重
    "max_weight": 86.2,        // 最高体重
    "avg_weight": 85.3,        // 平均体重
    "weight_change": -0.7,     // 体重变化(正数为增加，负数为减少)
    "trend_data": [            // 趋势数据点，用于绘制图表
      {
        "date": "2023-04-01T00:00:00Z",
        "weight_kg": 86.2
      },
      {
        "date": "2023-04-15T00:00:00Z",
        "weight_kg": 85.8
      },
      {
        "date": "2023-05-01T00:00:00Z",
        "weight_kg": 85.5
      }
    ]
  }
}
```

**说明**
- 统计数据基于指定天数内的体重记录
- trend_data 按时间顺序排列，可用于绘制体重变化曲线
- 如果没有体重记录会返回错误

## 每日营养相关接口（需要认证）

### 获取今日营养数据

**请求**
```
GET /api/v1/nutrition/today
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "date": "2023-05-01",
    "calories_intake": 1200.5,
    "protein_intake_g": 65.2,
    "carb_intake_g": 150.3,
    "fat_intake_g": 40.1,
    "target_calories": 1800.0,
    "target_protein_g": 90.0,
    "target_carb_g": 220.0,
    "target_fat_g": 60.0,
    "calories_completion_rate": 66.69,
    "created_at": "2023-05-01T08:30:00Z",
    "updated_at": "2023-05-01T18:45:00Z"
  }
}
```

### 更新今日营养摄入数据

**请求**
```
PUT /api/v1/nutrition/today
```

**请求参数**
```json
{
  "calories_intake": 1500.0,
  "protein_intake_g": 75.5,
  "carb_intake_g": 180.2,
  "fat_intake_g": 45.8
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "date": "2023-05-01",
    "calories_intake": 1500.0,
    "protein_intake_g": 75.5,
    "carb_intake_g": 180.2,
    "fat_intake_g": 45.8,
    "target_calories": 1800.0,
    "target_protein_g": 90.0,
    "target_carb_g": 220.0,
    "target_fat_g": 60.0,
    "calories_completion_rate": 83.33,
    "created_at": "2023-05-01T08:30:00Z",
    "updated_at": "2023-05-01T19:15:00Z"
  }
}
```

### 获取营养历史记录

**请求**
```
GET /api/v1/nutrition/history?start_date=2023-04-25&end_date=2023-05-01
```

**查询参数**
- start_date: 开始日期，格式YYYY-MM-DD
- end_date: 结束日期，格式YYYY-MM-DD

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "date": "2023-05-01",
      "calories_intake": 1500.0,
      "protein_intake_g": 75.5,
      "carb_intake_g": 180.2,
      "fat_intake_g": 45.8,
      "target_calories": 1800.0,
      "target_protein_g": 90.0,
      "target_carb_g": 220.0,
      "target_fat_g": 60.0,
      "calories_completion_rate": 83.33,
      "created_at": "2023-05-01T08:30:00Z",
      "updated_at": "2023-05-01T19:15:00Z"
    },
    // ... 其他日期的记录
  ]
}
```

### 获取周营养摄入统计

**请求**
```
GET /api/v1/nutrition/weekly-summary
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "avg_calories": 1450.5,
    "avg_protein": 72.3,
    "avg_carb": 175.8,
    "avg_fat": 48.2,
    "avg_completion_rate": 80.58
  }
}
```

## AI对话相关接口（需要认证）

### 创建聊天会话

**请求**
```
POST /api/v1/chat/sessions
```

**请求参数**
```json
{
  "title": "午餐咨询" // 可选，会话标题
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": "sess_1234567890abcdef",
    "user_id": 1,
    "title": "午餐咨询",
    "created_at": "2023-05-01T10:30:00Z",
    "updated_at": "2023-05-01T10:30:00Z"
  }
}
```

### 获取会话列表

**请求**
```
GET /api/v1/chat/sessions
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": "sess_1234567890abcdef",
      "user_id": 1,
      "title": "午餐咨询",
      "created_at": "2023-05-01T10:30:00Z",
      "updated_at": "2023-05-01T11:45:00Z"
    },
    {
      "id": "sess_abcdef1234567890",
      "user_id": 1,
      "title": "健康建议",
      "created_at": "2023-04-28T14:20:00Z",
      "updated_at": "2023-04-28T14:55:00Z"
    }
  ]
}
```

### 更新会话标题

**请求**
```
PUT /api/v1/chat/sessions/{session_id}
```

**请求参数**
```json
{
  "title": "减脂午餐咨询"
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 删除会话

**请求**
```
DELETE /api/v1/chat/sessions/{session_id}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取会话消息历史

**请求**
```
GET /api/v1/chat/sessions/{session_id}/messages
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "session_id": "sess_1234567890abcdef",
      "user_id": 1,
      "role": "user",
      "content": "我想知道中午吃什么比较健康?",
      "created_at": "2023-05-01T10:30:30Z"
    },
    {
      "id": 2,
      "session_id": "sess_1234567890abcdef",
      "user_id": 1,
      "role": "assistant",
      "content": "午餐建议摄入均衡的蛋白质、碳水和蔬菜...",
      "created_at": "2023-05-01T10:30:32Z"
    }
  ]
}
```

### 发送消息 (流式响应)

**请求**
```
POST /api/v1/chat/sessions/{session_id}/messages
```

**请求参数**
```json
{
  "content": "我应该如何搭配一顿减脂午餐?"
}
```

**响应**
- 接口会返回一个 `Content-Type: text/event-stream` 的流式响应。
- 客户端应监听 `message` 事件来接收AI回复的文本块。
- 流结束时，连接将自动关闭。

**响应示例**
```
event: message
data: 减脂午餐

event: message
data: 可以

event: message
data: 考虑以下

event: message
data: 搭配：

event: message
data: \n1. 主食：

event: message
data: 选择全谷物...
```

**说明**
- **重要变更**: 此接口已从返回单个JSON对象改为返回 Server-Sent Events (SSE) 流。前端需要相应地调整来处理流式数据。
- 响应不再包含完整的 `user_message` 和 `assistant_message` 对象。客户端发送消息后，会立即开始接收AI的流式回复。
- 用户发送的消息会由后端保存，但不会在本次响应中返回。
- 完整的AI回复需要客户端将所有 `message` 事件的 `data` 拼接起来。

## 食物识别相关接口（需要认证）

### 识别食物图片

**请求**
```
POST /api/v1/food/recognize
```

**说明**
- 使用multipart/form-data格式上传
- 文件大小限制10MB

**表单参数**
- food_image: 食物图片文件
- session_id: 可选，关联的聊天会话ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "image_url": "uploads/user_1/1683806400_abcdef.jpg",
    "recognized_foods": [
      {
        "name": "烤鸡胸肉",
        "quantity": "约100克",
        "calories": 165
      },
      {
        "name": "糙米饭",
        "quantity": "约150克",
        "calories": 180
      },
      {
        "name": "西兰花",
        "quantity": "约80克",
        "calories": 27
      }
    ],
    "nutrition_summary": {
      "calories_intake": 372,
      "protein_intake_g": 35.6,
      "carb_intake_g": 42.8,
      "fat_intake_g": 7.2
    },
    "ai_analysis": "这是一顿均衡的健康餐，蛋白质来源充足，含有复合碳水和蔬菜，总热量适中。"
  }
}
```

### 获取识别记录详情

**请求**
```
GET /api/v1/food/recognition/{id}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 123,
    "image_url": "uploads/user_1/1683806400_abcdef.jpg",
    "recognized_foods": [
      {
        "name": "烤鸡胸肉",
        "quantity": "约100克",
        "calories": 165
      },
      {
        "name": "糙米饭",
        "quantity": "约150克",
        "calories": 180
      },
      {
        "name": "西兰花",
        "quantity": "约80克",
        "calories": 27
      }
    ],
    "nutrition_summary": {
      "calories_intake": 372,
      "protein_intake_g": 35.6,
      "carb_intake_g": 42.8,
      "fat_intake_g": 7.2
    },
    "ai_analysis": "这是一顿均衡的健康餐，蛋白质来源充足，含有复合碳水和蔬菜，总热量适中。",
    "is_adopted": false,
    "record_date": "2023-05-01"
  }
}
```

### 获取今日识别记录

**请求**
```
GET /api/v1/food/recognition/today
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 123,
      "image_url": "uploads/user_1/1683806400_abcdef.jpg",
      "recognized_foods": [
        {
          "name": "烤鸡胸肉",
          "quantity": "约100克",
          "calories": 165
        },
        {
          "name": "糙米饭",
          "quantity": "约150克",
          "calories": 180
        },
        {
          "name": "西兰花",
          "quantity": "约80克",
          "calories": 27
        }
      ],
      "nutrition_summary": {
        "calories_intake": 372,
        "protein_intake_g": 35.6,
        "carb_intake_g": 42.8,
        "fat_intake_g": 7.2
      },
      "ai_analysis": "这是一顿均衡的健康餐，蛋白质来源充足，含有复合碳水和蔬菜，总热量适中。",
      "is_adopted": true,
      "record_date": "2023-05-01"
    },
    // ... 其他今日记录
  ]
}
```

### 保存食物识别结果到营养摄入

**请求**
```
POST /api/v1/food/recognition/{id}/save
```

**说明**
- 使用该接口将食物识别的营养数据保存到用户当日的营养摄入记录中
- 用户需要先查看识别结果后决定是否保存，而不是自动保存
- 保存后识别记录的`is_adopted`字段会被更新为`true`，表示该记录已被采用

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取用户已采用的食物识别记录

**请求**
```
GET /api/v1/food/recognition/adopted?page=1&page_size=20&start_date=2023-01-01&end_date=2023-01-31
```

**查询参数**
- `page`: 可选，页码，默认为1
- `page_size`: 可选，每页记录数，默认为20，最大100
- `start_date`: 可选，开始日期，格式YYYY-MM-DD，默认为30天前
- `end_date`: 可选，结束日期，格式YYYY-MM-DD，默认为当天

**说明**
- 使用该接口获取用户已保存到营养摄入的食物识别记录
- 返回结果包含按日期分组的记录，便于前端展示历史饮食记录

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "total": 15,
    "records": [
      {
        "id": 123,
        "image_url": "uploads/user_1/1683806400_abcdef.jpg",
        "recognized_foods": [
          {
            "name": "烤鸡胸肉",
            "quantity": "约100克",
            "calories": 165
          },
          {
            "name": "糙米饭",
            "quantity": "约150克",
            "calories": 180
          }
        ],
        "nutrition_summary": {
          "calories_intake": 345,
          "protein_intake_g": 30.5,
          "carb_intake_g": 40.8,
          "fat_intake_g": 6.2
        },
        "ai_analysis": "这是一顿均衡的健康餐，蛋白质来源充足，含有复合碳水，总热量适中。",
        "is_adopted": true,
        "record_date": "2023-05-01"
      },
      // ... 其他记录
    ],
    "date_groups": {
      "2023-05-01": [
        {
          "id": 123,
          "image_url": "uploads/user_1/1683806400_abcdef.jpg",
          "recognized_foods": [
            {
              "name": "烤鸡胸肉",
              "quantity": "约100克",
              "calories": 165
            },
            {
              "name": "糙米饭",
              "quantity": "约150克",
              "calories": 180
            }
          ],
          "nutrition_summary": {
            "calories_intake": 345,
            "protein_intake_g": 30.5,
            "carb_intake_g": 40.8,
            "fat_intake_g": 6.2
          },
          "ai_analysis": "这是一顿均衡的健康餐，蛋白质来源充足，含有复合碳水，总热量适中。",
          "is_adopted": true,
          "record_date": "2023-05-01"
        },
        // ... 同一天的其他记录
      ],
      "2023-04-30": [
        // ... 另一天的记录
      ]
      // ... 其他日期分组
    }
  }
}
```

## 运动记录相关接口（需要认证）

### 创建运动记录

**请求**
```
POST /api/v1/exercise
```

**请求参数**
```json
{
  "exercise_type": "跑步",           // 运动类型，必填
  "duration_min": 30.5,             // 持续时间（分钟），必填，大于0
  "calories_burned": 250.0,         // 消耗热量（千卡），必填，大于等于0
  "distance_km": 5.2,               // 距离（公里），可选
  "start_time": "2023-12-01T10:30:00Z"  // 运动开始时间，RFC3339格式，必填
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "exercise_type": "跑步",
    "duration_min": 30.5,
    "calories_burned": 250.0,
    "distance_km": 5.2,
    "start_time": "2023-12-01T10:30:00Z",
    "created_at": "2023-12-01T10:30:15Z",
    "updated_at": "2023-12-01T10:30:15Z"
  }
}
```

### 获取单个运动记录

**请求**
```
GET /api/v1/exercise/{id}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "exercise_type": "跑步",
    "duration_min": 30.5,
    "calories_burned": 250.0,
    "distance_km": 5.2,
    "start_time": "2023-12-01T10:30:00Z",
    "created_at": "2023-12-01T10:30:15Z",
    "updated_at": "2023-12-01T10:30:15Z"
  }
}
```

### 获取运动历史记录

**请求**
```
GET /api/v1/exercise/history?start_date=2023-12-01&end_date=2023-12-07&limit=20
```

**查询参数**
- start_date: 开始日期，YYYY-MM-DD格式，必填
- end_date: 结束日期，YYYY-MM-DD格式，必填
- limit: 限制数量，可选

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "exercise_type": "跑步",
      "duration_min": 30.5,
      "calories_burned": 250.0,
      "distance_km": 5.2,
      "start_time": "2023-12-01T10:30:00Z",
      "created_at": "2023-12-01T10:30:15Z",
      "updated_at": "2023-12-01T10:30:15Z"
    }
    // ... 其他记录
  ]
}
```

### 获取今日运动记录

**请求**
```
GET /api/v1/exercise/today
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "exercise_type": "跑步",
      "duration_min": 30.5,
      "calories_burned": 250.0,
      "distance_km": 5.2,
      "start_time": "2023-12-01T10:30:00Z",
      "created_at": "2023-12-01T10:30:15Z",
      "updated_at": "2023-12-01T10:30:15Z"
    }
    // ... 其他今日记录
  ]
}
```

### 更新运动记录

**请求**
```
PUT /api/v1/exercise/{id}
```

**请求参数**
```json
{
  "exercise_type": "慢跑",           // 可选
  "duration_min": 35.0,             // 可选
  "calories_burned": 280.0,         // 可选
  "distance_km": 6.0,               // 可选
  "start_time": "2023-12-01T10:30:00Z"  // 可选
}
```

**说明**
- 所有字段均为可选，只更新提供的字段

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "exercise_type": "慢跑",
    "duration_min": 35.0,
    "calories_burned": 280.0,
    "distance_km": 6.0,
    "start_time": "2023-12-01T10:30:00Z",
    "created_at": "2023-12-01T10:30:15Z",
    "updated_at": "2023-12-01T11:15:30Z"
  }
}
```

### 删除运动记录

**请求**
```
DELETE /api/v1/exercise/{id}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取运动统计数据

**请求**
```
GET /api/v1/exercise/statistics?start_date=2023-12-01&end_date=2023-12-07
```

**查询参数**
- start_date: 开始日期，YYYY-MM-DD格式，必填
- end_date: 结束日期，YYYY-MM-DD格式，必填

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "total_exercises": 5,          // 总运动次数
    "total_duration": 150.5,       // 总持续时间（分钟）
    "total_calories": 1250.0,      // 总消耗热量（千卡）
    "total_distance": 25.8,        // 总距离（公里）
    "avg_duration": 30.1,          // 平均持续时间（分钟）
    "avg_calories": 250.0          // 平均消耗热量（千卡）
  }
}
```

### 获取运动选项配置

**请求**
```
GET /api/v1/exercise/options
```

**说明**
- 获取前端展示用的运动记录选项配置
- 无需参数，返回所有可选的配置项

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "exercise_types": [
      "跑步", "走路", "骑行", "游泳", "瑜伽",
      "健身", "篮球", "足球", "网球", "羽毛球",
      "乒乓球", "爬山", "跳舞", "其他"
    ]
  }
}
```

## 心情记录相关接口（需要认证）

### 创建心情记录

**请求**
```
POST /api/v1/mood
```

**请求参数**
```json
{
  "time_context": "now",           // 时间上下文："now"表示当下，"today"表示当天，必填
  "mood_level": 3,                 // 情绪等级：1-7级（1=非常愉快，4=不悲不喜，7=非常不愉快），必填
  "mood_tags": ["平静", "满足"],    // 情绪标签数组，可选
  "influences": ["工作", "家人"]    // 影响因素数组，可选
}
```

**最简请求示例**
```json
{
  "time_context": "now",
  "mood_level": 4
}
```

**说明**
- time_context: 记录的是什么时候的情绪，必填
- mood_level: 情绪等级，必填。1=非常愉快, 2=愉快, 3=有点愉快, 4=不悲不喜, 5=有点不愉快, 6=不愉快, 7=非常不愉快
- mood_tags: 描述具体的情绪感受，可选，支持多选
- influences: 对情绪产生影响的因素，可选，支持多选

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "time_context": "now",
    "mood_level": 3,
    "mood_tags": ["平静", "满足"],
    "influences": ["工作", "家人"],
    "record_time": "2023-12-01T15:30:00Z",
    "created_at": "2023-12-01T15:30:05Z"
  }
}
```

### 获取单个心情记录

**请求**
```
GET /api/v1/mood/{id}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "time_context": "now",
    "mood_level": 3,
    "mood_tags": ["平静", "满足"],
    "influences": ["工作", "家人"],
    "record_time": "2023-12-01T15:30:00Z",
    "created_at": "2023-12-01T15:30:05Z"
  }
}
```

### 获取心情历史记录

**请求**
```
GET /api/v1/mood/history?start_date=2023-12-01&end_date=2023-12-07&limit=20
```

**查询参数**
- start_date: 开始日期，YYYY-MM-DD格式，必填
- end_date: 结束日期，YYYY-MM-DD格式，必填
- limit: 限制数量，可选

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "time_context": "now",
      "mood_level": 3,
      "mood_tags": ["平静", "满足"],
      "influences": ["工作", "家人"],
      "record_time": "2023-12-01T15:30:00Z",
      "created_at": "2023-12-01T15:30:05Z"
    }
    // ... 其他记录
  ]
}
```

### 获取今日心情记录

**请求**
```
GET /api/v1/mood/today
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "time_context": "now",
      "mood_level": 3,
      "mood_tags": ["平静", "满足"],
      "influences": ["工作", "家人"],
      "record_time": "2023-12-01T15:30:00Z",
      "created_at": "2023-12-01T15:30:05Z"
    }
    // ... 其他今日记录
  ]
}
```

### 删除心情记录

**请求**
```
DELETE /api/v1/mood/{id}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取心情统计数据

**请求**
```
GET /api/v1/mood/statistics?start_date=2023-12-01&end_date=2023-12-07
```

**查询参数**
- start_date: 开始日期，YYYY-MM-DD格式，必填
- end_date: 结束日期，YYYY-MM-DD格式，必填

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "total_records": 8,            // 总记录数
    "avg_mood_level": 3.5          // 平均情绪等级
  }
}
```

### 获取心情选项配置

**请求**
```
GET /api/v1/mood/options
```

**说明**
- 获取前端展示用的心情记录选项配置
- 无需参数，返回所有可选的配置项

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "time_contexts": ["now", "today"],
    "mood_levels": {
      "1": "非常愉快",
      "2": "愉快",
      "3": "有点愉快",
      "4": "不悲不喜",
      "5": "有点不愉快",
      "6": "不愉快",
      "7": "非常不愉快"
    },
    "influences": [
      "健康", "健身", "自我照顾", "爱好",
      "身份", "心灵", 
      "社群", "家人", "朋友", "伴侣", "约会",
      "家务", "工作", "教育", "旅行",
      "天气", "时事", "金钱"
    ],
    "common_mood_tags": [
      "平静", "开心", "兴奋", "满足", "放松",
      "烦躁", "焦虑", "沮丧", "愤怒", "疲惫",
      "无聊", "紧张", "困惑", "失望", "孤独"
    ]
  }
}
```

### 提交订单

```
GET /api/v1/order/create
```

**请求示例**
```json
{
  "store_id": 1,                      // 必填，门店ID
  "customer_name": "张三",             // 选填，客户姓名
  "customer_phone": "13800138000",    // 选填，客户电话
  "order_type": "dine_in",            // 必填，订单类型：dine_in-堂食，takeaway-自取，delivery-外送
  "payment_method": "wechat",         // 选填，支付方式
  "voucher_amount": 5.0,              // 选填，代金券抵扣金额
  "remark": "少盐少油",                // 选填，订单备注
  "items": [                          // 必填，订单项列表，至少1项
    {
      "dish_id": 1,                   // 必填，菜品ID
      "sku_info": {                   // 必填，SKU信息
        "id": 1,
        "name": "大份",
        "price": 25.0,
        "weight": 300,
        "weight_unit": "g"
      },
      "options_info": [               // 选填，选项信息
        {
          "group_id": 1,
          "group_name": "辣度",
          "items": [
            {
              "id": 1,
              "name": "微辣",
              "extra_price": 0
            }
          ]
        }
      ],
      "quantity": 2,                  // 必填，数量
      "remark": "不要香菜"             // 选填，备注
    }
  ]
}
```

**响应**
```json
{
    "code": 0,
    "data": {
        "OrderID": 1,
        "OrderNo": "20250804238C4CDA",
        "StoreID": 1,
        "CustomerName": "张三",
        "CustomerPhone": "13800138000",
        "TotalAmount": 50,
        "ActualAmount": 49,
        "VoucherAmount": 5,
        "PackagingFee": 4,
        "OrderType": "dine_in",
        "PaymentMethod": "wechat",
        "Remark": "少盐少油",
        "CreatedAt": "2025-08-04 18:55:07.607 +0800 CST",
        "UpdatedAt": "2025-08-04 18:55:07.607 +0800 CST"
    },
    "msg": "成功"
}
```


## 文件访问相关接口

### 获取文件（公共）

**请求**
```
GET /api/v1/files/{filepath}
```

**说明**
- 使用此接口获取系统中可公开访问的文件
- 接口会返回文件内容而非JSON响应
- 仅允许访问uploads目录下的文件，出于安全考虑有路径限制

**响应**
文件内容（非JSON格式），同时设置适当的Content-Type头部

### 获取用户文件（需要认证）

**请求**
```
GET /api/v1/user/files/{filepath}
```

**说明**
- 使用此接口获取当前认证用户有权限访问的文件
- 接口会返回文件内容而非JSON响应
- 仅允许用户访问自己uploads/user_{user_id}目录下的文件

**响应**
文件内容（非JSON格式），同时设置适当的Content-Type头部 
