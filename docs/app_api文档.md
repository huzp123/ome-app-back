# OME-APP 前端API接口文档

## 通用说明

*有任何相关变化都需要同步到该api文档中*

### 基础URL
```
/api/v1
```

### 认证方式
除了登录、注册和健康检查API外，所有请求都需要在Header中加入认证信息：
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

## 用户相关接口（需要认证）

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
  "height_cm": 175.0,           // 身高(厘米)，必须在50-300cm范围内
  "birth_date": "2000-01-01",   // 出生日期，格式YYYY-MM-DD
  "sex": "male",                // 性别: male/female/other
  "weight_kg": 70.5             // 当前体重(公斤)
}
```

**说明**
- 所有字段均为可选，只需填写需要更新的信息
- 身高字段有限制，必须在50-300cm范围内
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
  "goal_type": "lose_fat",              // 目标类型: lose_fat/keep_fit/gain_muscle
  "target_weight_kg": 65.0,             // 目标体重(公斤)
  "weekly_change_kg": 0.5,              // 每周计划变化的体重(公斤)
  "target_date": "2023-12-31",          // 目标日期，格式YYYY-MM-DD
  "diet_type": "normal",                // 饮食类型: normal/vegetarian/low_carb等
  "taste_preferences": ["清淡", "酸的"], // 口味偏好
  "food_intolerances": ["海鲜"]         // 食物不耐受/禁忌
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
    "diet_type": "normal",                // 饮食类型: normal/vegetarian/low_carb等
    "taste_preferences": ["清淡", "酸的"], // 口味偏好
    "food_intolerances": ["海鲜"],         // 食物不耐受/禁忌
    "created_at": "2023-04-01T12:00:00Z"  // 创建时间
  }
}
```

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

### 发送消息

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
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "user_message": {
      "id": 3,
      "session_id": "sess_1234567890abcdef",
      "user_id": 1,
      "role": "user",
      "content": "我应该如何搭配一顿减脂午餐?",
      "created_at": "2023-05-01T10:35:00Z"
    },
    "assistant_message": {
      "id": 4,
      "session_id": "sess_1234567890abcdef",
      "user_id": 1,
      "role": "assistant",
      "content": "减脂午餐可以考虑以下搭配：\n1. 主食：选择全谷物...",
      "created_at": "2023-05-01T10:35:02Z"
    }
  }
}
```

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
    },
    // ... 其他今日记录
  ]
}
``` 