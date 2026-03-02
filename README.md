# 个人博客后端

一个使用 Go 编写的个人博客系统后端。

整体功能较为基础，主要用于个人学习与日常使用。

## 功能

当前已支持：

* 文章列表
* 文章详情
* 网站相关数据（大杂烩，k-v形式）
* 基础配置读取
* 日志输出

整体以常规 CRUD 为主，没有引入复杂架构设计。

## 技术栈

* Go
* Gin
* GORM
* PostgreSQL
* Redis
* Viper

## 项目结构

```
backend
├─cmd               # 启动入口
├─config            # 配置文件
└─internal
    ├─admin
    ├─article       # 文章模块
    ├─config        # 配置读取
    ├─global 
    ├─infra         # 数据库、Redis 初始化
    ├─logger        # 日志
    ├─middleware    # 中间件
    ├─response      # 响应封装
    ├─user          # 用户模块
    ├─utils 
    └─websiteData   # 网站相关数据
```

结构按模块进行拆分，保持相对清晰。

---

## 快速开始

1. 克隆项目

    ```
    git clone https://github.com/zygame03/z-blog-backend.git
    ```

2. 修改配置

    `config` 文件夹中提供示例配置文件，复制一份并重命名为 `config.json`，然后根据自己的环境修改相关参数。

3. 运行项目

    本地直接运行：

    ```
    go run .\cmd\main.go
    ```

    或编译为二进制后部署到服务器。

    可以直接运行 `build.bat` 脚本编译 Linux 可执行文件：

    ```
    .\build.bat
    ```

## Todo

* [ ] 完善用户模块与管理员接口
* [ ] 部分代码结构重构
* [ ] 持续修复问题

项目会根据实际使用情况持续调整结构
