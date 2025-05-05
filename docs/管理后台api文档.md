# OME-SAAS 后端API接口文档

**注意 如果有变动需要同步到该文档中**

## 通用说明

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

### 错误码说明
- 0: 成功
- 10000: 服务内部错误
- 10001: 无效参数
- 10002: 找不到资源
- 10003: 未授权认证失败
- 10004: 未授权Token错误
- 10005: 未授权Token超时
- 10006: 请求过多
- 20001: 用户不存在
- 20002: 用户已存在
- 20003: 用户密码错误
- 20004: 创建用户失败
- 20005: 更新用户失败
- 20006: 删除用户失败
- 20007: 用户唯一编码无效
- 30001: 门店不存在
- 30002: 门店已存在
- 30003: 创建门店失败
- 30004: 更新门店失败
- 30005: 删除门店失败

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
  "username": "string",  // 必填，用户名，3-50字符
  "password": "string",  // 必填，密码，6-20字符
  "real_name": "string", // 选填，真实姓名
  "phone": "string",     // 选填，电话号码
  "email": "string"      // 选填，邮箱
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "username": "string",
    "real_name": "string",
    "phone": "string",
    "email": "string",
    "unique_code": "string",
    "role": 3,           // 角色：1-系统管理员，2-店家管理员，3-店员
    "parent_id": 0,
    "store_id": 0,
    "status": 1          // 状态：1-启用，0-禁用
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
  "username": "string",  // 必填，用户名
  "password": "string"   // 必填，密码
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "token": "string",   // JWT令牌
    "user": {
      "id": 1,
      "username": "string",
      "real_name": "string",
      "phone": "string",
      "email": "string",
      "unique_code": "string",
      "role": 3,
      "parent_id": 0,
      "store_id": 0,
      "status": 1
    }
  }
}
```

## 用户相关接口

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
    "username": "string",
    "real_name": "string",
    "phone": "string",
    "email": "string",
    "unique_code": "string",
    "role": 3,
    "parent_id": 0,
    "store_id": 0,
    "status": 1
  }
}
```

### 更新用户信息

**请求**
```
PUT /user/info
```

**请求参数**
```json
{
  "real_name": "string", // 选填，真实姓名
  "phone": "string",     // 选填，电话号码
  "email": "string",     // 选填，邮箱
  "status": 1            // 选填，状态：1-启用，0-禁用
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

### 修改密码

**请求**
```
PUT /user/password
```

**请求参数**
```json
{
  "old_password": "string",  // 必填，旧密码
  "new_password": "string"   // 必填，新密码，6-20字符
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

### 获取下级用户列表

**请求**
```
GET /user/subordinates
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "username": "string",
      "real_name": "string",
      "phone": "string",
      "email": "string",
      "unique_code": "string",
      "role": 3,
      "parent_id": 0,
      "store_id": 0,
      "status": 1
    }
  ]
}
```

### 通过唯一编码添加用户

**请求**
```
POST /user/add-by-code
```

**请求参数**
```json
{
  "unique_code": "string",  // 必填，用户唯一编码
  "store_id": 1             // 选填，门店ID，如果指定则将用户分配到该门店
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

### 获取用户权限

**请求**
```
GET /user/permissions
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "username": "string",
    "real_name": "string",
    "phone": "string",
    "email": "string",
    "unique_code": "string",
    "role": 3,
    "parent_id": 0,
    "store_id": 0,
    "status": 1,
    "permissions": [
      {
        "id": 1,
        "name": "string",
        "key": "string",
        "description": "string",
        "module": "string"
      }
    ]
  }
}
```

### 批量授予用户权限

**请求**
```
POST /user/grant-permissions
```

**请求参数**
```json
{
  "user_id": 1,              // 必填，用户ID
  "permission_ids": [1, 2]   // 必填，权限ID列表
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

## 门店相关接口

### 创建门店

**请求**
```
POST /store
```

**请求参数**
```json
{
  "name": "string",          // 必填，门店名称，1-100字符
  "address": "string",       // 选填，门店地址
  "phone": "string",         // 选填，联系电话
  "description": "string",   // 选填，描述
  "manager_id": 1            // 选填，店长用户ID
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "string",
    "address": "string",
    "phone": "string",
    "description": "string",
    "status": 1,              // 状态：1-启用，0-禁用
    "creator_id": 1,
    "manager_id": 1           // 店长ID
  }
}
```

### 更新门店

**请求**
```
PUT /store/:id
```

**路径参数**
- id: 门店ID

**请求参数**
```json
{
  "name": "string",          // 选填，门店名称
  "address": "string",       // 选填，门店地址
  "phone": "string",         // 选填，联系电话
  "description": "string",   // 选填，描述
  "status": 1,               // 选填，状态：1-启用，0-禁用
  "manager_id": 1            // 选填，店长用户ID
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

### 删除门店

**请求**
```
DELETE /store/:id
```

**路径参数**
- id: 门店ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取门店详情

**请求**
```
GET /store/:id
```

**路径参数**
- id: 门店ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "string",
    "address": "string",
    "phone": "string",
    "description": "string",
    "status": 1,
    "creator_id": 1,
    "manager_id": 1          // 店长ID
  }
}
```

### 获取门店列表

**请求**
```
GET /store
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "name": "string",
      "address": "string",
      "phone": "string",
      "description": "string",
      "status": 1,
      "creator_id": 1,
      "manager_id": 1        // 店长ID
    }
  ]
}
```

### 获取指定创建者的门店列表

**请求**
```
GET /store/by-creator/:creator_id
```

**路径参数**
- creator_id: 创建者用户ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "name": "string",
      "address": "string",
      "phone": "string",
      "description": "string",
      "status": 1,
      "creator_id": 1,
      "manager_id": 1        // 店长ID
    }
  ]
}
```

### 获取门店下的用户列表

**请求**
```
GET /store/:id/users
```

**路径参数**
- id: 门店ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "string",
    "address": "string",
    "phone": "string",
    "description": "string",
    "status": 1,
    "creator_id": 1,
    "manager_id": 1,          // 店长ID
    "users": [
      {
        "id": 1,
        "username": "string",
        "real_name": "string",
        "phone": "string",
        "email": "string",
        "unique_code": "string",
        "role": 3,
        "parent_id": 0,
        "store_id": 1,
        "status": 1
      }
    ]
  }
}
```

### 添加用户到门店

**请求**
```
POST /store/:id/users
```

**路径参数**
- id: 门店ID

**请求参数**
```json
{
  "user_id": 1  // 必填，用户ID
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

### 从门店移除用户

**请求**
```
DELETE /store/:id/users/:user_id
```

**路径参数**
- id: 门店ID
- user_id: 用户ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 设置门店店长

**请求**
```
PUT /store/:id/manager
```

**路径参数**
- id: 门店ID

**请求参数**
```json
{
  "manager_id": 1  // 必填，店长用户ID
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

## 权限相关接口

### 创建权限

**请求**
```
POST /permission
```

**请求参数**
```json
{
  "name": "string",          // 必填，权限名称，1-100字符
  "key": "string",           // 必填，权限键，1-50字符，唯一
  "description": "string",   // 选填，描述
  "module": "string"         // 必填，所属模块，1-50字符
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "string",
    "key": "string",
    "description": "string",
    "module": "string"
  }
}
```

### 更新权限

**请求**
```
PUT /permission/:id
```

**路径参数**
- id: 权限ID

**请求参数**
```json
{
  "name": "string",          // 选填，权限名称
  "description": "string",   // 选填，描述
  "module": "string"         // 选填，所属模块
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

### 删除权限

**请求**
```
DELETE /permission/:id
```

**路径参数**
- id: 权限ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取权限详情

**请求**
```
GET /permission/:id
```

**路径参数**
- id: 权限ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "string",
    "key": "string",
    "description": "string",
    "module": "string"
  }
}
```

### 获取权限列表

**请求**
```
GET /permission
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "name": "string",
      "key": "string",
      "description": "string",
      "module": "string"
    }
  ]
}
```

### 根据模块获取权限列表

**请求**
```
GET /permission/module/:module
```

**路径参数**
- module: 模块名称

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "name": "string",
      "key": "string",
      "description": "string",
      "module": "string"
    }
  ]
}
```

## 菜品分组相关接口

### 创建菜品分组

**请求**
```
POST /dish-group
```

**请求参数**
```json
{
  "name": "string",     // 必填，分组名称，1-64字符
  "sort_order": 0       // 选填，排序顺序
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "owner_id": 1,
    "name": "string",
    "sort_order": 0,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 更新菜品分组

**请求**
```
PUT /dish-group/:id
```

**路径参数**
- id: 分组ID

**请求参数**
```json
{
  "name": "string",     // 必填，分组名称
  "sort_order": 0       // 选填，排序顺序
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

### 删除菜品分组

**请求**
```
DELETE /dish-group/:id
```

**路径参数**
- id: 分组ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取菜品分组详情

**请求**
```
GET /dish-group/:id
```

**路径参数**
- id: 分组ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "owner_id": 1,
    "name": "string",
    "sort_order": 0,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取菜品分组列表

**请求**
```
GET /dish-group
```

**查询参数**
- owner_id: 选填，品牌管理员ID，筛选指定管理员下的菜品分组

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "owner_id": 1,
      "name": "string",
      "sort_order": 0,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

## 菜品相关接口

### 上传菜品图片

**请求**
```
POST /dish/upload-image
```

**请求参数**
使用表单数据(multipart/form-data)上传
- image: 图片文件，支持jpg、png、jpeg格式
  - 必须符合1:1比例，尺寸800*800px以上
  - 大小不超过3M
  - 每个菜品最少需上传1张，最多支持5张

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "url": "/uploads/dishes/image.jpg"
  }
}
```

### 创建菜品

**请求**
```
POST /dish
```

**请求参数**
```json
{
  "group_id": 1,                  // 必填，所属分组ID
  "name": "string",               // 必填，菜品名称，限制字数内
  "display_cat": "string",        // 选填，前端展示子分类
  "subtitle": "string",           // 选填，推荐文案、卖点，不超过200字
  "taste": "string",              // 必填，口味描述（前端以多选呈现：清淡、辣味、酸的、甜的）
  "temperature": "hot",           // 必填，温度，hot(热)或cold(凉)
  "protein_type": "meat",         // 选填，荤素类型，meat-荤 veg-素 mixed-荤素结合
  "packaging_fee": 2.0,           // 必填，打包费，默认0，支持一位小数
  "status": 1,                    // 选填，状态，1=启用 0=停用
  "images": [                     // 必填，菜品图片，至少1张，最多5张
    {
      "url": "string",            // 必填，图片URL
      "sort_order": 0             // 选填，排序顺序
    }
  ],
  "skus": [                       // 必填，至少一个规格
    {
      "name": "string",           // 必填，规格名称，10字以内
      "weight": 100,              // 必填，重量，仅可输入整数
      "weight_unit": "g",         // 必填，重量单位，仅支持"g"(克)或"个"
      "price": 15.5,              // 必填，价格，支持一位小数
      "is_default": 1             // 选填，是否默认规格，1=是 0=否
    }
  ],
  "ingredients": [                // 必填，食材列表
    {
      "name": "string",           // 必填，食材名称，10字以内
      "quantity": 100,            // 必填，用量，仅可输入整数
      "unit": "g",                // 必填，单位，仅支持"g"(克)或"个"
      "sort_order": 0             // 选填，排序顺序
    }
  ],
  "nutrition": {                  // 必填，营养信息
    "calories": 200,              // 必填，热量(千卡)，仅可输入整数
    "carbohydrate": 20.5,         // 必填，碳水化合物(克)，仅可输入整数
    "protein": 10.5,              // 必填，蛋白质(克)，仅可输入整数
    "fat": 5.5,                   // 必填，脂肪(克)，仅可输入整数
    "others": "string"            // 选填，其他营养信息，不超过50字
  },
  "option_groups": [              // 选填，选项组/属性列表
    {
      "name": "string",           // 必填，选项组名称，不超过10字
      "min_select": 0,            // 选填，最小选择数量
      "max_select": 1,            // 选填，最大选择数量
      "sort_order": 0,            // 选填，排序顺序
      "items": [                  // 必填，选项项列表
        {
          "name": "string",       // 必填，选项项名称，不超过8字
          "extra_price": 0,       // 选填，加价，支持一位小数
          "sort_order": 0         // 选填，排序顺序
        }
      ]
    }
  ]
}
```

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "owner_id": 1,
    "group_id": 1,
    "name": "string",
    "display_cat": "string",
    "subtitle": "string",
    "taste": "string",
    "temperature": "hot",
    "protein_type": "meat",
    "packaging_fee": 2.0,
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "images": [
      {
        "id": 1,
        "url": "string",
        "sort_order": 0
      }
    ],
    "skus": [
      {
        "id": 1,
        "name": "string",
        "weight": 100,
        "weight_unit": "g",
        "price": 15.5,
        "is_default": 1
      }
    ],
    "ingredients": [
      {
        "id": 1,
        "name": "string",
        "quantity": 100,
        "unit": "g",
        "sort_order": 0
      }
    ],
    "nutrition": {
      "calories": 200,
      "carbohydrate": 20.5,
      "protein": 10.5,
      "fat": 5.5,
      "others": "string"
    },
    "option_groups": [
      {
        "id": 1,
        "name": "string",
        "min_select": 0,
        "max_select": 1,
        "sort_order": 0,
        "items": [
          {
            "id": 1,
            "name": "string",
            "extra_price": 0,
            "sort_order": 0
          }
        ]
      }
    ]
  }
}
```

### 更新菜品

**请求**
```
PUT /dish/:id
```

**路径参数**
- id: 菜品ID

**请求参数**
```json
{
  "group_id": 1,                  // 选填，所属分组ID
  "name": "string",               // 选填，菜品名称，限制字数内
  "display_cat": "string",        // 选填，前端展示子分类
  "subtitle": "string",           // 选填，推荐文案、卖点，不超过200字
  "taste": "string",              // 选填，口味描述（前端以多选呈现：清淡、辣味、酸的、甜的）
  "temperature": "hot",           // 选填，温度，hot(热)或cold(凉)
  "protein_type": "meat",         // 选填，荤素类型，meat-荤 veg-素 mixed-荤素结合
  "packaging_fee": 2.0,           // 选填，打包费，支持一位小数
  "status": 1,                    // 选填，状态，1=启用 0=停用
  "images": [                     // 选填，菜品图片，至少1张，最多5张
    {
      "url": "string",            // 必填，图片URL
      "sort_order": 0             // 选填，排序顺序
    }
  ],
  "skus": [                       // 选填，规格列表
    {
      "name": "string",           // 必填，规格名称，10字以内
      "weight": 100,              // 必填，重量，仅可输入整数
      "weight_unit": "g",         // 必填，重量单位，仅支持"g"(克)或"个"
      "price": 15.5,              // 必填，价格，支持一位小数
      "is_default": 1             // 选填，是否默认规格，1=是 0=否
    }
  ],
  "ingredients": [                // 选填，食材列表
    {
      "name": "string",           // 必填，食材名称，10字以内
      "quantity": 100,            // 必填，用量，仅可输入整数
      "unit": "g",                // 必填，单位，仅支持"g"(克)或"个"
      "sort_order": 0             // 选填，排序顺序
    }
  ],
  "nutrition": {                  // 选填，营养信息
    "calories": 200,              // 必填，热量(千卡)，仅可输入整数
    "carbohydrate": 20.5,         // 必填，碳水化合物(克)，仅可输入整数
    "protein": 10.5,              // 必填，蛋白质(克)，仅可输入整数
    "fat": 5.5,                   // 必填，脂肪(克)，仅可输入整数
    "others": "string"            // 选填，其他营养信息，不超过50字
  },
  "option_groups": [              // 选填，选项组/属性列表
    {
      "name": "string",           // 必填，选项组名称，不超过10字
      "min_select": 0,            // 选填，最小选择数量
      "max_select": 1,            // 选填，最大选择数量
      "sort_order": 0,            // 选填，排序顺序
      "items": [                  // 必填，选项项列表
        {
          "name": "string",       // 必填，选项项名称，不超过8字
          "extra_price": 0,       // 选填，加价，支持一位小数
          "sort_order": 0         // 选填，排序顺序
        }
      ]
    }
  ]
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

### 删除菜品

**请求**
```
DELETE /dish/:id
```

**路径参数**
- id: 菜品ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": null
}
```

### 获取菜品详情

**请求**
```
GET /dish/:id
```

**路径参数**
- id: 菜品ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    // 与创建菜品的响应相同
  }
}
```

### 获取菜品列表

**请求**
```
GET /dish
```

**查询参数**
- group_id: 选填，分组ID，筛选指定分组下的菜品
- owner_id: 选填，品牌管理员ID，筛选指定管理员下的菜品

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      // 与获取菜品详情的响应相同
    }
  ]
}
```

## 门店菜品设置相关接口

### 设置门店菜品

**请求**
```
PUT /store/:id/dish/:dish_id/setting
```

**路径参数**
- id: 门店ID
- dish_id: 菜品ID

**请求参数**
```json
{
  "is_enabled": 1,               // 选填，是否启用，1=启用 0=停用
  "sale_start_time": "08:00:00", // 选填，每天销售开始时间
  "sale_end_time": "20:00:00",   // 选填，每天销售结束时间
  "sale_week_mask": "1111111"    // 选填，周一至周日是否销售，1=销售 0=不销售
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

### 获取门店菜品设置

**请求**
```
GET /store/:id/dish/:dish_id/setting
```

**路径参数**
- id: 门店ID
- dish_id: 菜品ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "store_id": 1,
    "dish_id": 1,
    "is_enabled": 1,
    "sale_start_time": "08:00:00", // 每天销售开始时间
    "sale_end_time": "20:00:00",   // 每天销售结束时间
    "sale_week_mask": "1111111",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取门店菜品设置列表

**请求**
```
GET /store/:id/dish/settings
```

**路径参数**
- id: 门店ID

**响应**
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "store_id": 1,
      "dish_id": 1,
      "is_enabled": 1,
      "sale_start_time": "08:00:00", // 每天销售开始时间
      "sale_end_time": "20:00:00",   // 每天销售结束时间
      "sale_week_mask": "1111111",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

### 批量设置门店菜品

**请求**
```
POST /store/:id/dish/batch-settings
```

**路径参数**
- id: 门店ID

**请求参数**
```json
{
  "dish_ids": [1, 2, 3],         // 必填，菜品ID列表
  "setting": {                    // 必填，设置内容
    "is_enabled": 1,
    "sale_start_time": "08:00:00", // 每天销售开始时间
    "sale_end_time": "20:00:00",   // 每天销售结束时间
    "sale_week_mask": "1111111"
  }
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