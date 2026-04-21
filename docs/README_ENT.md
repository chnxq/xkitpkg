# Ent 框架及工具集应用介绍

## 项目概述

Ent 是一个简单而强大的 Go 语言实体框架，它使得构建和维护具有大型数据模型的应用变得容易。本项目包含了 Ent 框架及其相关工具集。

## 主要应用 (Main Applications)

### 1. ent (entgo.io/ent/cmd/ent)

- **功能**: Ent 框架的主要命令行工具，用于代码生成和项目管理。
- **用途**: 创建新的 Ent 项目、生成代码、初始化项目结构等。
- **主要命令**:
  - `ent new [flags] [schemas]` - 创建新的 schema (已替代旧的 `ent init`)
  - `ent generate [flags] path` - 生成代码
  - `ent describe [flags] path` - 打印图 schema 的描述
  - `ent schema [flags] path` - 生成 DDL 语句用于数据库表创建
- **常用参数**:
  - `--target` - 指定 schema 目标目录 (默认: ./ent/schema)
  - `--storage` - 指定存储驱动 (默认: sql)
  - `--feature` - 启用额外功能 (例如: privacy, entql 等)
    - `privacy` - 隐私层功能，通过 schema 配置提供隐私保护
    - `entql` - 运行时通用过滤功能
    - `intercept` - 拦截器功能，简化拦截器的使用
    - `namedges` - 为边加载提供动态命名API
    - `bidiedges` - 为边加载设置双向引用
    - `sql/schemaconfig` - 允许为每个模型配置不同的 schema 名称
    - `sql/lock` - SQL 行级锁定功能
    - `sql/modifier` - 为查询添加自定义修饰符
    - `sql/execquery` - 暴露底层 SQL 驱动的 ExecContext/QueryContext 方法
    - `sql/upsert` - INSERT 语句的 upsert (ON CONFLICT) 功能
    - `sql/versioned-migration` - 版本化迁移功能
    - `sql/globalid` - 确保所有节点都有唯一的全局标识符
    - `schema/snapshot` - 存储 ent/schema 的快照并自动解决合并冲突
  - `--header` - 代码生成头部注释
  - `--dialect` - 指定数据库方言 (用于 schema 命令)
- **参数使用示例**:
  - 使用自定义目标目录: `ent generate --target ./mypass/schema ./ent/schema`
  - 使用不同存储驱动: `ent generate --storage gremlin ./ent/schema`
  - 启用隐私功能: `ent generate --feature privacy ./ent/schema`
  - 自定义头部注释: `ent generate --header "// Custom header" ./ent/schema`
  - 指定数据库方言: `ent schema --dialect mysql ./ent/schema`
- **综合示例**:
  ```bash
  go install entgo.io/ent/cmd/ent@latest
  ent new Todo
  ent generate ./ent/schema
  ent describe ./ent/schema
  # 带参数的生成示例
  ent generate --feature privacy,entql ./ent/schema
  ```

### 2. entc (entgo.io/ent/cmd/entc)

- **功能**: Ent 代码生成器的命令行界面，是 ent 工具的另一种实现方式。
- **用途**: 提供更高级的代码生成功能，支持更多自定义选项。
- **主要命令**:
  - `entc new [flags] [schemas]` - 创建新的 schema
  - `entc generate [flags] path` - 生成代码
  - `entc describe [flags] path` - 描述 schema 结构
- **特点**:
  - 更多自定义选项
  - 扩展性强
  - 适合集成到自定义代码生成流程中
- **常用参数**:
  - `--idtype` - ID 字段类型 (int, int64, uint, uint64, string)
  - `--template` - 指定外部模板
  - `--build-tags` - 指定 Go 构建标签

### 3. entfix (entgo.io/ent/cmd/entfix)

- **功能**: Ent 数据修复和迁移工具。
- **用途**: 用于修复和迁移 Ent 数据库中的数据，如全局 ID 迁移等。
- **主要功能**:
  - 全局 ID 迁移 (`entfix globalid`)
  - 数据库连接支持多种方言 (MySQL, PostgreSQL, SQLite3)
- **命令**:
  - `entfix globalid` - 迁移唯一全局 ID 类型到 Ent 全局功能
- **参数**:
  - `--dialect` - 数据库方言 (mysql, postgres, sqlite3)
  - `--dsn` - 数据源名称
  - `--path` - 生成的 Ent 代码路径
- **示例**:
  ```bash
  entfix globalid --dialect mysql --dsn "user:password@tcp(localhost:3306)/dbname" --path ./ent
  ```

### 4. entproto (entgo.io/contrib/entproto/cmd/entproto)

- **功能**: 将 Ent 模式转换为 Protocol Buffer 文件。
- **用途**: 用于生成 gRPC 服务定义和消息类型，实现与 Protocol Buffer 的集成。
- **参数**:
  - `-path` - 指定 schema 目录路径
- **示例**:
  ```bash
  entproto -path ./ent/schema
  ```
- **使用案例**:
  1. 在 schema 上添加注解:
     ```go
     func (User) Annotations() []schema.Annotation {
         return []schema.Annotation{
             entproto.Message(),
             entproto.Service(), // 生成 gRPC 服务定义
         }
     }
     ```
  2. 为字段添加 proto 字段号:
     ```go
     func (User) Fields() []ent.Field {
         return []ent.Field{
             field.String("user_name").
                 Annotations(entproto.Field(2)), // proto 字段号为 2
         }
     }
     ```
  3. 运行命令生成 proto 文件: `entproto -path ./ent/schema`
  4. 生成的 proto 文件位于 `./ent/proto/entpb/entpb.proto`

### 5. protoc-gen-ent (entgo.io/contrib/entproto/cmd/protoc-gen-ent)

- **功能**: Protocol Buffers 编译器插件，用于从 .proto 文件生成 Ent 模式。
- **用途**: 与 protoc 配合使用，实现从 Protocol Buffer 到 Ent 模式的双向转换。
- **特点**:
  - 作为 protoc 插件运行
  - 支持模式反向生成
  - 可通过 `-schemadir` 参数指定 schema 目录
- **使用案例**:
  1. 安装插件: `go get entgo.io/contrib/entproto/cmd/protoc-gen-ent`
  2. 在 proto 文件中导入并使用注解:
     ```protobuf
     syntax = "proto3";
     package entpb;
     
     import "options/opts.proto";  // 导入 entproto 选项
     option go_package = "example.com/project/ent/proto/entpb";
     
     message User {
       option (ent.schema).gen = true;  // 告诉 protoc-gen-ent 从此消息生成 schema
       string name = 1;
       string email_address = 2;
     }
     ```
  3. 运行 protoc 命令: `protoc -I=proto/ --ent_out=. --ent_opt=schemadir=./schema proto/entpb/user.proto`
  4. 生成的 schema 文件位于 `./schema/user.go`

### 6. entoas (entgo.io/contrib/entoas)

- **功能**: 生成完全兼容的可扩展 OpenAPI 规范文档，用于使用 Swagger 工具生成 RESTful 服务器存根和客户端。
- **用途**: 将 Ent 模式转换为 OpenAPI 规范文档，便于生成 RESTful API 接口。
- **特点**:
  - 生成符合 OpenAPI 3.0.3 规范的文档
  - 与 Swagger 工具链集成
  - 支持生成服务器存根和客户端代码
  - 可与 elk 项目配合使用生成完整的服务器实现
  - 使用 ogen 的 OAS 结构定义创建 OAS 文档
- **使用方法**:
  1. 安装模块: `go get -u entgo.io/contrib/entoas`
  2. 在 entc.go 中注册 entoas 扩展
  3. 运行代码生成器后会生成 openapi.json 文件
- **使用案例**:
  1. 创建 entc.go 文件配置扩展:
     ```go
     ex, err := entoas.NewExtension(
         entoas.SimpleModels(),
         entoas.Mutations(func(_ *gen.Graph, spec *ogen.Spec) error {
             spec.Info.SetTitle("My Simple API").
                 SetDescription("API to demonstrate **simple model** generation.").
                 SetVersion("0.0.1")
             return nil
         }),
     )
     ```
  2. 在 ent/generate.go 中引用 entc.go: `//go:generate go run -mod=mod entc.go`
  3. 运行 `go generate ./...` 后会在项目根目录生成 openapi.json 文件

### 7. protoc-gen-entgrpc (entgo.io/contrib/entproto/cmd/protoc-gen-entgrpc)

- **功能**: Protocol Buffers 编译器插件，用于从 .proto 文件生成 gRPC 服务代码。
- **用途**: 生成与 Ent 模式对应的 gRPC 服务端和客户端代码。
- **参数**:
  - `-schema_path` - 指定 Ent schema 路径
- **特点**:
  - 生成完整的 gRPC 服务实现
  - 与 Ent 模式紧密集成
- **使用案例**:
  1. 安装插件: `go get entgo.io/contrib/entproto/cmd/protoc-gen-entgrpc`
  2. 使用 entproto 生成 proto 文件后，系统会自动创建 generate.go 文件
  3. 运行 `go generate ./ent/proto/...` 来生成 gRPC 服务代码
  4. 或手动运行 protoc 命令:
     ```bash
     protoc -I=.. --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --entgrpc_out=.. --entgrpc_opt=paths=source_relative,schema_path=../../schema --go-grpc_opt=paths=source_relative entpb/entpb.proto
     ```
  5. 生成的 gRPC 服务实现在 `ent/proto/entpb/` 目录下

## 使用方法

### 安装主工具

```bash
# 安装 ent 主工具
go install entgo.io/ent/cmd/ent@latest

# 安装 entc 工具
go install entgo.io/ent/cmd/entc@latest

# 安装 entfix 工具
go install entgo.io/ent/cmd/entfix@latest
```

### 安装扩展工具

```bash
# 安装 entoas
go get -u entgo.io/contrib/entoas

# 安装 entproto
go install entgo.io/contrib/entproto/cmd/entproto@latest

# 安装 protoc-gen-ent 插件
go install entgo.io/contrib/entproto/cmd/protoc-gen-ent@latest

# 安装 protoc-gen-entgrpc 插件
go install entgo.io/contrib/entproto/cmd/protoc-gen-entgrpc@latest
```

## 适用场景

- **ent**: 适用于所有基于 Ent 框架的项目开发，特别是新项目快速启动
- **entc**: 适用于需要更高级代码生成功能的项目，或需要深度定制代码生成过程的场景
- **entfix**: 适用于需要数据迁移或修复的生产环境，特别是处理全局 ID 等复杂迁移任务
- **entoas**: 适用于需要生成 RESTful API 和 OpenAPI 文档的项目，便于前端开发和第三方集成
- **entproto 及相关工具**: 适用于需要 gRPC 和 Protocol Buffer 集成的微服务架构

## 特性与功能

Ent 框架提供以下核心功能：
- **Schema As Code** - 将任何数据库 schema 建模为 Go 对象
- **轻松遍历任意图** - 轻松运行查询、聚合和遍历任何图结构
- **静态类型和显式 API** - 100% 静态类型和显式 API，使用代码生成
- **多存储驱动** - 支持 MySQL、MariaDB、TiDB、PostgreSQL、CockroachDB、SQLite 和 Gremlin
- **可扩展** - 使用 Go 模板轻松扩展和自定义

## 许可证

Ent 框架及其工具遵循 Apache 2.0 许可证。