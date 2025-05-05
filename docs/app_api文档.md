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
  "height_cm": 175.0,         // 身高(厘米)
  "birth_date": "2000-01-01", // 出生日期，格式YYYY-MM-DD
  "sex": "male",              // 性别: male/female/other
  "weight_kg": 70.5           // 当前体重(公斤)
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