# CurrencyExchangeApp

一个基于 Go + Gin 的汇率与文章管理后端服务。项目提供用户注册登录、JWT 鉴权、汇率记录管理、文章发布查询与点赞接口，使用 MySQL 存储数据，Redis 用于文章详情等缓存场景。

## 功能特性

- 用户注册、登录与 JWT 鉴权
- 汇率数据创建与查询
- 文章创建、分页查询、详情查询
- 文章点赞与点赞数查询
- MySQL 数据持久化，GORM 自动迁移表结构
- Redis 缓存支持
- Docker Compose 一键启动应用、MySQL、Redis

## 技术栈

- Go 1.25.4
- Gin
- GORM
- MySQL 8
- Redis 7
- Viper
- JWT
- Docker / Docker Compose

## 项目结构

```text
.
├── config/             # 配置加载、数据库和 Redis 初始化
├── controllers/        # HTTP 控制器
├── global/             # 全局数据库、Redis 等对象
├── middlewares/        # Gin 中间件
├── models/             # GORM 数据模型
├── router/             # 路由注册
├── services/           # 业务逻辑
├── utils/              # 响应、日志、JWT、错误处理等工具
├── docker-compose.yml  # Docker Compose 编排
├── dockerfile          # 应用镜像构建文件
├── Makefile            # 常用开发命令
├── go.mod
└── main.go
```

## 快速开始

### 使用 Docker Compose 启动

推荐直接使用 Docker Compose，它会同时启动应用、MySQL 和 Redis。

```bash
docker-compose up -d --build
```

服务默认访问地址：

```text
http://localhost:3000
```

停止服务：

```bash
docker-compose down
```

### 本地启动

本地启动前需要先准备 MySQL 和 Redis，并根据实际连接地址修改 `config/config.yml`。

默认配置如下：

```yaml
app:
  port: :3000

database:
  dsn: "root:123456@tcp(mysql:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: redis:6379
```

如果在宿主机直接运行程序，通常需要把数据库和 Redis 地址改为：

```yaml
database:
  dsn: "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: 127.0.0.1:6379
```

安装依赖并启动：

```bash
go mod tidy
go run main.go
```

构建二进制文件：

```bash
go build -o bin/exchangeapp main.go
```

## 配置说明

主要配置文件位于 `config/config.yml`。

| 配置项 | 说明 |
| --- | --- |
| `app.name` | 应用名称 |
| `app.port` | HTTP 服务监听端口 |
| `app.mode` | 运行模式 |
| `app.log_level` | 日志级别 |
| `app.log_file` | 日志文件路径 |
| `database.dsn` | MySQL 连接字符串 |
| `database.max_idle_conns` | 最大空闲连接数 |
| `database.max_open_conns` | 最大打开连接数 |
| `database.conn_max_lifetime` | 连接最大生命周期，单位秒 |
| `redis.addr` | Redis 地址 |
| `redis.password` | Redis 密码 |
| `redis.db` | Redis 数据库编号 |
| `jwt.secret` | JWT 签名密钥 |
| `jwt.expiration_hours` | Token 过期时间，单位小时 |

生产环境请务必修改 `jwt.secret`、MySQL 密码等敏感配置。

## API 接口

接口统一返回结构：

```json
{
  "code": 0,
  "message": "Success",
  "data": {}
}
```

### 认证

#### 注册

```http
POST /api/auth/register
Content-Type: application/json
```

请求示例：

```json
{
  "username": "test",
  "password": "123456",
  "email": "test@example.com"
}
```

#### 登录

```http
POST /api/auth/login
Content-Type: application/json
```

请求示例：

```json
{
  "username": "test",
  "password": "123456"
}
```

登录成功后会返回 `token`，调用受保护接口时需要放入请求头：

```http
Authorization: <token>
```

### 汇率

#### 查询汇率列表

```http
GET /api/exchangeRates
```

#### 创建汇率

需要登录。

```http
POST /api/exchangeRates
Authorization: <token>
Content-Type: application/json
```

请求示例：

```json
{
  "from_currency": "USD",
  "to_currency": "CNY",
  "rate": 7.25000000
}
```

### 文章

以下文章接口需要登录。

#### 创建文章

```http
POST /api/articles
Authorization: <token>
Content-Type: application/json
```

请求示例：

```json
{
  "title": "汇率市场观察",
  "content": "这里是文章正文内容",
  "preview": "这里是文章摘要",
  "author_id": 1,
  "status": 1
}
```

#### 分页查询文章

```http
GET /api/articles?page=1&page_size=10
Authorization: <token>
```

#### 查询文章详情

```http
GET /api/articles/{id}
Authorization: <token>
```

#### 点赞文章

```http
POST /api/articles/{id}/like
Authorization: <token>
```

#### 查询文章点赞数

```http
GET /api/articles/{id}/likes
Authorization: <token>
```

## 常用命令

```bash
make build        # 构建项目
make run          # 运行项目
make test         # 运行测试
make fmt          # 格式化代码
make vet          # 静态检查
make docker-build # 构建 Docker 镜像
make docker-run   # 启动 Docker Compose
make docker-stop  # 停止 Docker Compose
make mod-tidy     # 整理 Go 依赖
```

## 数据表

应用启动时会通过 GORM 自动迁移以下数据表：

- `users`
- `exchange_rates`
- `articles`

## 注意事项

- 当前 CORS 默认允许 `http://localhost:5173`，如前端地址不同，需要在 `router/router.go` 中调整。
- Docker 环境下配置里的 `mysql`、`redis` 是 Compose 服务名；本地运行时请改为 `127.0.0.1` 或实际服务地址。
- `Authorization` 请求头目前直接传 token 字符串，不需要加 `Bearer` 前缀。
