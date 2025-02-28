# 依赖注入指南：增加新功能时如何自动生成依赖注入

## 概述

本项目使用 Google 的 [Wire](https://github.com/google/wire) 框架来实现依赖注入。Wire 是一个代码生成工具，它使用「依赖注入」设计模式，无需大量手动编写初始化代码。当你需要增加新功能时，只需按照以下步骤操作，就能自动生成依赖注入代码。

## 步骤

### 1. 创建新组件

首先，创建你的新组件（如新的 Repository、Service 或 Controller）。确保它们有明确的构造函数，并且构造函数接收所有必要的依赖作为参数。

**示例：创建新的 Repository 接口**

```go
// internal/repository/product.go
package repository

import "gin-server/internal/entity"

// ProductRepository 产品仓库接口
type ProductRepository interface {
	Create(product *entity.Product) error
	GetByID(id uint) (*entity.Product, error)
	// 其他方法...
}
```

**示例：创建 MySQL 实现**

```go
// internal/repository/mysql/product.go
package mysql

import "gin-server/internal/entity"

// ProductRepository MySQL实现的产品仓库
type ProductRepository struct {
	// 可能的依赖...
}

// NewProductRepository 创建MySQL产品仓库实例
func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// Create 创建产品
func (r *ProductRepository) Create(product *entity.Product) error {
	// 实现...
	return nil
}

// GetByID 根据ID获取产品
func (r *ProductRepository) GetByID(id uint) (*entity.Product, error) {
	// 实现...
	return nil, nil
}
```

**示例：创建 MongoDB 实现**

```go
// internal/repository/mongodb/product.go
package mongodb

import "gin-server/internal/entity"

// ProductRepository MongoDB实现的产品仓库
type ProductRepository struct {
	// MongoDB 相关字段...
}

// NewProductRepository 创建MongoDB产品仓库实例
func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// 实现接口方法...
```

**示例：创建 Service**

```go
// internal/service/product.go
package service

import (
	"gin-server/internal/entity"
	"gin-server/internal/repository"
)

// ProductService 产品服务
type ProductService struct {
	productRepo repository.ProductRepository
}

// NewProductService 创建产品服务
func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// CreateProduct 创建产品
func (s *ProductService) CreateProduct(product *entity.Product) error {
	return s.productRepo.Create(product)
}

// GetProduct 获取产品
func (s *ProductService) GetProduct(id uint) (*entity.Product, error) {
	return s.productRepo.GetByID(id)
}
```

**示例：创建 Controller**

```go
// internal/controller/product.go
package controller

import (
	"gin-server/internal/config"
	"gin-server/internal/service"
	"github.com/gin-gonic/gin"
)

// ProductController 产品控制器
type ProductController struct {
	productService *service.ProductService
	jwtConfig     *config.JWTConfig
}

// NewProductController 创建产品控制器
func NewProductController(productService *service.ProductService, jwtConfig *config.JWTConfig) *ProductController {
	return &ProductController{
		productService: productService,
		jwtConfig:     jwtConfig,
	}
}

// CreateProduct 创建产品
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	// 实现...
}

// GetProduct 获取产品
func (c *ProductController) GetProduct(ctx *gin.Context) {
	// 实现...
}
```

### 2. 更新 Container 结构体

在 `internal/di/container.go` 文件中，更新 `Container` 结构体，添加新组件：

```go
// Container 依赖注入容器
type Container struct {
	Config            *config.Config
	DB                *gorm.DB
	UserRepository    repository.UserRepository
	UserService       *service.UserService
	UserController    *controller.UserController
	// 添加新组件
	ProductRepository repository.ProductRepository
	ProductService    *service.ProductService
	ProductController *controller.ProductController
}
```

### 3. 添加 Provider 函数

在 `internal/di/wire.go` 文件中，添加新组件的 provider 函数：

```go
// provideProductRepository 根据配置提供产品仓库实现
func provideProductRepository(cfg *config.Config, db *gorm.DB) (repository.ProductRepository, error) {
	// 根据配置选择仓库实现
	switch cfg.Database.Driver {
	case "mysql":
		return mysql.NewProductRepository(), nil
	case "mongodb":
		return mongodb.NewProductRepository(), nil
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Driver)
	}
}

// provideProductController 提供产品控制器
func provideProductController(productService *service.ProductService, cfg *config.Config) *controller.ProductController {
	return controller.NewProductController(productService, &cfg.JWT)
}
```

### 4. 更新 Wire.Build

在 `InitializeContainer` 函数中，更新 `wire.Build` 调用，添加新的 provider 函数：

```go
// InitializeContainer 使用Wire注入依赖
func InitializeContainer(cfg *config.Config) (*Container, error) {
	wire.Build(
		wire.Struct(new(Container), "*"),
		provideDB,
		provideUserRepository,
		service.NewUserService,
		provideUserController,
		// 添加新组件
		provideProductRepository,
		service.NewProductService,
		provideProductController,
	)
	return nil, nil
}
```

### 5. 生成 Wire 代码

运行以下命令生成依赖注入代码：

```bash
go run github.com/google/wire/cmd/wire
```

或者如果你在项目根目录：

```bash
cd internal/di
go run -mod=mod github.com/google/wire/cmd/wire
```

这将生成更新后的 `wire_gen.go` 文件，其中包含所有依赖注入的实现代码。

### 6. 更新路由

最后，在 `Container` 的 `SetupRouter` 方法中添加新的路由：

```go
// SetupRouter 设置路由
func (c *Container) SetupRouter() *gin.Engine {
	// 现有代码...

	// 公共路由组
	public := router.Group("/api/v1")
	{
		// 用户相关路由...

		// 产品相关路由
		productGroup := public.Group("/products")
		{
			productGroup.POST("", c.ProductController.CreateProduct)
			productGroup.GET("/:id", c.ProductController.GetProduct)
		}
	}

	// 需要认证的路由组
	authorized := router.Group("/api/v1")
	// authorized.Use(middleware.JWTAuth())
	{
		// 用户相关路由...

		// 产品相关路由
		productGroup := authorized.Group("/products")
		{
			// 需要认证的产品路由...
		}
	}

	return router
}
```

## 最佳实践

1. **保持构造函数简单**：构造函数应该只接收必要的依赖，并返回接口而不是具体实现。

2. **使用接口**：尽可能使用接口而不是具体类型，这样可以更容易地进行单元测试和替换实现。

3. **分层依赖**：确保依赖关系是单向的，例如 Controller 依赖 Service，Service 依赖 Repository。

4. **错误处理**：在 provider 函数中正确处理错误，并将其传播到调用者。

5. **避免循环依赖**：确保你的依赖图是无环的，否则 Wire 将无法生成代码。

## 示例：完整流程

以下是添加新功能（产品管理）的完整流程示例：

1. 创建实体：`internal/entity/product.go`
2. 创建仓库接口：`internal/repository/product.go`
3. 创建仓库实现：`internal/repository/mysql/product.go` 和 `internal/repository/mongodb/product.go`
4. 创建服务：`internal/service/product.go`
5. 创建控制器：`internal/controller/product.go`
6. 更新容器：`internal/di/container.go`
7. 添加 provider 函数：`internal/di/wire.go`
8. 生成依赖注入代码：运行 Wire 命令
9. 更新路由：`internal/di/container.go` 中的 `SetupRouter` 方法

按照这些步骤，你可以轻松地向项目添加新功能，并自动生成所有必要的依赖注入代码。