# z-blog-backend

一个基于 Go 构建的个人博客后端服务，面向博客站点与简单后台管理场景。

当前围绕文章、评论、站点信息、弹幕、用户鉴权、访问统计和动态配置展开，整体采用按业务模块拆分的方式组织代码，便于继续扩展。

## 项目简介

当前以实现以下能力：

- 文章列表、热门文章与文章详情查询
- 文章发布与评论审核
- 弹幕提交与首屏弹幕加载
- 站点简介、公告、弹幕展示与弹幕审核
- 用户注册、登录与基础资料查询
- 访问统计与查询
- 基于配置中心的动态模块配置读取与更新
- Swagger 接口文档

## 技术栈

- Go
- Gin
- GORM
- PostgreSQL
- Redis
- Viper
- Zap
- JWT
- Swagger

## 运行依赖

基础环境：

- Go 1.25+
- PostgreSQL
- Redis

## 配置说明

项目使用两类配置文件：

- `config/config.json`：静态配置，包含 HTTP 服务、数据库、Redis、JWT 等基础配置
- `config/dyconfig.json`：动态配置，包含文章、站点、统计等模块配置

初始化方式：

1. 复制 `config/config.example.json`
2. 重命名为 `config/config.json`
3. 根据本地或服务器环境填写数据库、Redis 和 `jwt_key` 等配置
4. 按需调整 `config/dyconfig.json`

示例中的关键配置包括：

- 服务端口与 CORS 白名单
- PostgreSQL 连接信息
- Redis 连接信息
- JWT 密钥
- 缓存 TTL 与定时同步间隔

## 快速开始

1. 克隆项目

```bash
git clone https://github.com/zygame03/z-blog-backend.git
cd z-blog-backend
```

2. 准备配置文件

```bash
cp config/config.example.json config/config.json
```

Windows PowerShell 也可以直接手动复制并重命名。

3. 安装依赖并启动

```bash
go run .\cmd\main.go
```

服务默认监听 `:8080`，接口基础前缀为 `/api`。

## 接口入口

常用接口入口如下：

- 文章模块：`/api/article`
- 用户模块：`/api/user`
- 站点模块：`/api/site`
- 管理端弹幕接口：`/api/admin/site`
- 统计模块：`/api/stats`
- 配置中心：`/api/home/config`
- Swagger 文档：`/api/home/swagger/index.html`

## 目录结构

当前项目主要目录如下：

```text
z-blog-backend
├─ cmd/                    # 程序启动入口
├─ config/                 # 静态配置与动态配置文件
├─ docs/                   # Swagger 生成产物
├─ internal/
│  ├─ article/             # 文章与评论相关逻辑
│  ├─ config/              # 配置加载与动态配置管理
│  ├─ global/              # 通用模型定义
│  ├─ home/                # Swagger 与配置中心接口
│  ├─ http/                # 统一响应结构与 schema 定义
│  ├─ infra/               # 数据库、Redis、Cron、HTTP Server 初始化
│  ├─ logger/              # 日志封装
│  ├─ middleware/          # JWT 与访问统计中间件
│  ├─ site/                # 站点信息与弹幕相关逻辑
│  ├─ stats/               # 访问统计逻辑
│  ├─ user/                # 用户注册、登录与资料
│  └─ zerrors/             # 业务错误定义
├─ build.bat               # 构建脚本
├─ go.mod
└─ README.md
```

## 模块说明

### article

- 提供文章分页、热门文章、文章详情查询
- 支持文章发布
- 支持评论提交与评论审核

### site

- 提供站点简介与公告展示
- 提供弹幕查询、发送和管理审核能力

### user

- 提供注册、登录和用户资料查询
- 登录成功后返回 JWT，可用于受保护接口访问

### stats

- 统计访问数据
- 支持按配置周期进行缓存与同步

### home

- 提供 Swagger 页面访问
- 提供动态配置 schema 查询与模块配置更新接口

## 响应约定

接口统一返回 JSON，基础结构如下：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

字段说明：

- `code`：业务状态码
- `message`：状态描述
- `data`：接口返回数据，类型随业务变化

## TODO

- 管理端接口边界与权限控制
- 更完整的错误码体系
- 自动化测试覆盖率
- 评论、文章与站点配置的后台化管理体验
