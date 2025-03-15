# 后端 API 接口文档

## 基础信息

- **基础URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **字符编码**: UTF-8

## 统一响应格式

所有接口返回统一的响应格式：

```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": {}
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，0表示成功，其他值表示失败 |
| message | string | 响应消息 |
| data | any | 响应数据，根据接口不同而变化 |

### 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 500 | 失败 |
| 1001 | 请求参数错误 |
| 1002 | 数据库操作异常 |
| 2001 | 密码错误 |
| 2002 | 用户不存在 |

## 认证说明

部分接口需要JWT认证，需要在请求头中携带：

```
Authorization: Bearer {token}
```

Token通过登录接口获取，有效期为24小时。

---

## 一、文章相关接口

### 1.1 获取文章列表

**接口地址**: `/api/article`

**请求方法**: `GET`

**是否需要认证**: 否

**请求参数**:

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| pageSize | int | 否 | 10 | 每页数量 |

**请求示例**:
```
GET /api/article?page=1&pageSize=10
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": {
    "page": 1,
    "size": 10,
    "total": 100,
    "data": [
      {
        "id": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "title": "文章标题",
        "authorName": "作者名",
        "views": 100,
        "tags": "标签1,标签2",
        "cover": "封面图片URL"
      }
    ]
  }
}
```

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| page | int | 当前页码 |
| size | int | 每页数量 |
| total | int | 总记录数 |
| data | array | 文章列表 |
| data[].id | int | 文章ID |
| data[].created_at | string | 创建时间 |
| data[].updated_at | string | 更新时间 |
| data[].title | string | 文章标题 |
| data[].authorName | string | 作者名 |
| data[].views | int | 浏览量 |
| data[].tags | string | 标签（逗号分隔） |
| data[].cover | string | 封面图片URL |

---

### 1.2 获取热门文章

**接口地址**: `/api/article/hotArticles`

**请求方法**: `GET`

**是否需要认证**: 否

**请求参数**: 无

**请求示例**:
```
GET /api/article/hotArticles
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": [
    {
      "id": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "title": "文章标题",
      "authorName": "作者名",
      "views": 1000,
      "tags": "标签1,标签2",
      "cover": "封面图片URL"
    }
  ]
}
```

**响应字段说明**: 同文章列表接口的data字段

---

### 1.3 获取文章详情

**接口地址**: `/api/article/:id`

**请求方法**: `GET`

**是否需要认证**: 否

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文章ID |

**请求示例**:
```
GET /api/article/1
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": {
    "id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "title": "文章标题",
    "desc": "文章描述",
    "content": "文章正文内容",
    "authorName": "作者名",
    "views": 101,
    "tags": "标签1,标签2",
    "cover": "封面图片URL",
    "status": 0
  }
}
```

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 文章ID |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |
| title | string | 文章标题 |
| desc | string | 文章描述 |
| content | string | 文章正文内容 |
| authorName | string | 作者名 |
| views | int | 浏览量（访问后自动+1） |
| tags | string | 标签（逗号分隔） |
| cover | string | 封面图片URL |
| status | int | 状态（0:公开, 1:私有） |

---

## 二、用户相关接口

### 2.1 用户登录

**接口地址**: `/api/user/login`

**请求方法**: `POST`

**是否需要认证**: 否

**请求体**:

```json
{
  "username": "用户名",
  "password": "密码"
}
```

**请求示例**:
```json
POST /api/user/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| data | string | JWT Token，有效期24小时 |

---

### 2.2 用户注册

**接口地址**: `/api/user/register`

**请求方法**: `POST`

**是否需要认证**: 否

**请求体**:

```json
{
  "username": "用户名",
  "password": "密码"
}
```

**请求示例**:
```json
POST /api/user/register
Content-Type: application/json

{
  "username": "newuser",
  "password": "123456"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": {
    "id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "avatar": "头像URL",
    "username": "newuser",
    "signature": "个性签名",
    "github": "GitHub链接",
    "bilibili": "B站链接"
  }
}
```

---

### 2.3 获取用户资料

**接口地址**: `/api/user/profile`

**请求方法**: `GET`

**是否需要认证**: 否

**请求示例**:
```
GET /api/user/profile
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": {
    "id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "avatar": "头像URL",
    "username": "用户名",
    "signature": "个性签名",
    "github": "GitHub链接",
    "bilibili": "B站链接"
  }
}
```

---

### 2.4 获取个人介绍

**接口地址**: `/api/user/intro`

**请求方法**: `GET`

**是否需要认证**: 否

**请求示例**:
```
GET /api/user/intro
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": "个人介绍内容..."
}
```

---

### 2.5 获取技能列表

**接口地址**: `/api/user/skills`

**请求方法**: `GET`

**是否需要认证**: 否

**请求示例**:
```
GET /api/user/skills
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": "技能1,技能2,技能3"
}
```

---

### 2.6 获取兴趣爱好

**接口地址**: `/api/user/hobbies`

**请求方法**: `GET`

**是否需要认证**: 否

**请求示例**:
```
GET /api/user/hobbies
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": "爱好1,爱好2,爱好3"
}
```

---

### 2.7 获取时间线

**接口地址**: `/api/user/timeline`

**请求方法**: `GET`

**是否需要认证**: 否

**请求示例**:
```
GET /api/user/timeline
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": "时间线内容..."
}
```

---

### 2.8 获取未来目标

**接口地址**: `/api/user/futureGoals`

**请求方法**: `GET`

**是否需要认证**: 否

**请求示例**:
```
GET /api/user/futureGoals
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": "未来目标内容..."
}
```

---

## 三、管理员文章管理接口

> **注意**: 以下所有接口都需要JWT认证，请在请求头中携带 `Authorization: Bearer {token}`

### 3.1 获取文章列表（管理员）

**接口地址**: `/api/admin/article`

**请求方法**: `GET`

**是否需要认证**: 是

**请求参数**:

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| pagesize | int | 否 | 10 | 每页数量 |

**请求示例**:
```
GET /api/admin/article?page=1&pagesize=10
Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": {
    "page": 1,
    "size": 10,
    "total": 100,
    "data": [
      {
        "id": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "title": "文章标题",
        "desc": "文章描述",
        "content": "文章正文",
        "authorName": "作者名",
        "views": 100,
        "tags": "标签1,标签2",
        "cover": "封面图片URL",
        "status": 0
      }
    ]
  }
}
```

---

### 3.2 获取文章详情（管理员）

**接口地址**: `/api/admin/article/:id`

**请求方法**: `GET`

**是否需要认证**: 是

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文章ID |

**请求示例**:
```
GET /api/admin/article/1
Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**: 同3.1接口响应格式

---

### 3.3 创建或更新文章

**接口地址**: `/api/admin/article`

**请求方法**: `POST`

**是否需要认证**: 是

**请求体**:

```json
{
  "id": 0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "title": "文章标题",
  "desc": "文章描述",
  "content": "文章正文内容",
  "authorName": "作者名",
  "views": 0,
  "tags": "标签1,标签2",
  "cover": "封面图片URL",
  "status": 0
}
```

**请求参数说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 否 | 文章ID，大于0表示更新，等于0或不传表示新建 |
| title | string | 是 | 文章标题 |
| desc | string | 否 | 文章描述 |
| content | string | 是 | 文章正文 |
| authorName | string | 否 | 作者名 |
| views | int | 否 | 浏览量，默认0 |
| tags | string | 否 | 标签（逗号分隔） |
| cover | string | 否 | 封面图片URL |
| status | int | 否 | 状态（0:公开, 1:私有），默认0 |

**请求示例（新建）**:
```json
POST /api/admin/article
Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "title": "新文章标题",
  "desc": "文章描述",
  "content": "文章正文内容",
  "authorName": "作者名",
  "tags": "标签1,标签2",
  "cover": "封面图片URL",
  "status": 0
}
```

**请求示例（更新）**:
```json
POST /api/admin/article
Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "id": 1,
  "title": "更新后的标题",
  "desc": "更新后的描述",
  "content": "更新后的正文",
  "authorName": "作者名",
  "tags": "新标签1,新标签2",
  "cover": "新封面URL",
  "status": 0
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": {
    "id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "title": "文章标题",
    "desc": "文章描述",
    "content": "文章正文",
    "authorName": "作者名",
    "views": 0,
    "tags": "标签1,标签2",
    "cover": "封面图片URL",
    "status": 0
  }
}
```

---

### 3.4 删除文章

**接口地址**: `/api/admin/article/:id`

**请求方法**: `DELETE`

**是否需要认证**: 是

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 文章ID |

**请求示例**:
```
DELETE /api/admin/article/1
Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
  "code": 0,
  "message": "SUCCESS",
  "data": null
}
```

---

## 四、错误响应示例

### 参数错误
```json
{
  "code": 1001,
  "message": "请求参数错误",
  "data": null
}
```

### 数据库操作异常
```json
{
  "code": 1002,
  "message": "数据库操作异常",
  "data": null
}
```

### 认证失败
```json
{
  "error": "missing token"
}
```
或
```json
{
  "error": "invalid token"
}
```

---

## 五、注意事项

1. 所有时间字段使用ISO 8601格式（UTC时间）
2. JWT Token有效期为24小时，过期后需要重新登录
3. 文章状态：0表示公开，1表示私有
4. 文章删除为软删除，不会真正从数据库中删除
5. 获取文章详情时会自动增加浏览量（views + 1）
6. 分页查询默认返回公开且未删除的文章
7. 热门文章按浏览量（views）降序排列，默认返回前10条

---

## 六、更新日志

- 2024-01-01: 初始版本接口文档

