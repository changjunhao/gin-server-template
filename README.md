# Gin Server Template

基于Gin框架的Go Web服务器模版，支持MySQL和MongoDB两种数据库实现。

> **注意**: 本仓库的`di-wire`分支是基于Google Wire框架实现依赖注入的模板版本。

## 项目概述

本项目是一个使用Go语言和Gin框架开发的Web服务器模板，采用了清晰的分层架构设计，支持MySQL和MongoDB两种数据库实现方式。项目通过配置文件可以轻松切换数据库类型，适合快速开发各类Web应用。

## 技术栈

- **Go**: 1.23.6+
- **Web框架**: Gin
- **数据库**: 支持MySQL和MongoDB
- **ORM**: GORM (MySQL)
- **认证**: JWT
- **配置管理**: Viper

## 项目结构

```
├── cmd/                # 应用程序入口
│   └── api/            # API服务入口
├── configs/            # 配置文件
├── internal/           # 内部包
│   ├── app/            # 应用程序初始化
│   ├── config/         # 配置结构定义
│   ├── controller/     # 控制器层
│   ├── database/       # 数据库连接管理
│   ├── entity/         # 实体模型定义
│   ├── middleware/     # 中间件
│   ├── repository/     # 数据访问层
│   │   ├── mongodb/    # MongoDB实现
│   │   └── mysql/      # MySQL实现
│   └── service/        # 业务逻辑层
└── pkg/                # 公共包
    └── response/       # 响应处理
```

## 如何切换数据库

在`configs/config.yaml`文件中修改以下配置：

```yaml
# 数据库配置
database:
  driver: mysql  # 可选值: mysql, mongodb
  # 其他配置...
```

## 安装和使用

### 前置条件

- Go 1.23.6+
- MySQL 或 MongoDB

### 安装步骤

1. 克隆仓库

```bash
git clone <repository-url>
cd gin-server
```

2. 安装依赖

```bash
go mod download
```

3. 配置数据库

编辑`configs/config.yaml`文件，设置数据库连接信息。

4. 运行服务

```bash
go run cmd/api/main.go
```

## API文档

启动服务后，可以通过以下端点访问API：

- 默认端口: 8080
- 健康检查: GET /health
- 用户API: 
  - 注册: POST /api/v1/users/register
  - 登录: POST /api/v1/users/login
  - 获取用户信息: GET /api/v1/users/:id
