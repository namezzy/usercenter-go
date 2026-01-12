# 用户中心后端 - Go语言版本

> 基于原 Java Spring Boot 项目改写的 Go 语言版本

## 项目简介

这是一个用户中心管理系统的 Go 语言实现版本，包含用户注册、登录、查询等基础功能。

原项目作者：[程序员鱼皮](https://github.com/liyupi)  
Go版本改写：基于原项目迁移

## 技术栈

### 后端
- **Go 1.21+** - 编程语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **MySQL** - 数据库
- **Gin Sessions** - Session 管理

## 项目结构

```
go-version/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序
├── internal/              # 内部代码
│   ├── common/           # 通用组件
│   │   ├── constants.go  # 常量定义
│   │   ├── error_code.go # 错误码
│   │   └── response.go   # 统一响应
│   ├── config/           # 配置管理
│   │   └── config.go     # 配置加载
│   ├── controller/       # 控制器层
│   │   └── user_controller.go
│   ├── middleware/       # 中间件
│   │   └── auth.go       # 认证中间件
│   ├── model/           # 数据模型
│   │   ├── user.go      # 用户实体
│   │   └── request.go   # 请求DTO
│   ├── repository/      # 数据访问层
│   │   └── user_repository.go
│   └── service/         # 业务逻辑层
│       └── user_service.go
├── pkg/                 # 公共库
│   └── utils/          # 工具函数
│       └── crypto.go   # 加密工具
├── config.yaml         # 配置文件
└── go.mod             # 依赖管理
```

## 主要功能

### 用户模块
- ✅ **用户注册** - 账号密码注册，支持星球编号
- ✅ **用户登录** - Session 认证
- ✅ **用户注销** - 清除登录态
- ✅ **获取当前用户** - 查看登录用户信息
- ✅ **搜索用户** - 管理员按用户名搜索（模糊查询）
- ✅ **删除用户** - 管理员删除用户（逻辑删除）

### 安全特性
- 密码 MD5 加密（加盐）
- 用户信息脱敏
- Session 会话管理
- 管理员权限控制
- 参数校验

## 快速开始

### 1. 环境要求
- Go 1.21 或更高版本
- MySQL 5.7 或更高版本

### 2. 数据库初始化
```bash
# 导入数据库表结构
mysql -u root -p yupi < ../sql/create_table.sql
```

### 3. 修改配置
编辑 `config.yaml` 文件，配置数据库连接信息：
```yaml
database:
  host: localhost
  port: 3306
  database: yupi
  username: root
  password: 你的密码
```

### 4. 安装依赖
```bash
cd go-version
go mod download
```

### 5. 运行项目
```bash
go run cmd/main.go
```

服务将启动在 `http://localhost:8080`

## API 接口

### 公开接口
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录

### 需要登录
- `POST /api/user/logout` - 用户注销
- `GET /api/user/current` - 获取当前用户

### 需要管理员权限
- `GET /api/user/search?username=xxx` - 搜索用户
- `POST /api/user/delete` - 删除用户

## API 示例

### 用户注册
```bash
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "userAccount": "testuser",
    "userPassword": "12345678",
    "checkPassword": "12345678",
    "planetCode": "12345"
  }'
```

### 用户登录
```bash
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -c cookie.txt \
  -d '{
    "userAccount": "testuser",
    "userPassword": "12345678"
  }'
```

### 获取当前用户
```bash
curl -X GET http://localhost:8080/api/user/current \
  -b cookie.txt
```

## 响应格式

### 成功响应
```json
{
  "code": 0,
  "data": {...},
  "message": "ok",
  "description": ""
}
```

### 错误响应
```json
{
  "code": 40000,
  "data": null,
  "message": "请求参数错误",
  "description": "参数为空"
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40000 | 请求参数错误 |
| 40001 | 请求数据为空 |
| 40100 | 未登录 |
| 40101 | 无权限 |
| 50000 | 系统内部异常 |

## 与 Java 版本的对应关系

| Java | Go |
|------|------|
| Spring Boot | Gin Framework |
| MyBatis Plus | GORM |
| Controller | Controller |
| Service | Service |
| Mapper | Repository |
| @Autowired | 构造函数注入 |
| Session (HttpServletRequest) | Gin Sessions |
| Lombok @Data | 手动定义 struct |

## 开发说明

### 添加新的接口
1. 在 `model/request.go` 中定义请求结构体
2. 在 `repository` 层添加数据访问方法
3. 在 `service` 层添加业务逻辑
4. 在 `controller` 层添加接口处理
5. 在 `main.go` 中注册路由

### 数据库迁移
建议使用 SQL 文件手动迁移，也可以取消 `main.go` 中的 `AutoMigrate` 注释：
```go
err = db.AutoMigrate(&model.User{})
```

## 注意事项

1. **Session Secret** - 生产环境请修改 `config.yaml` 中的 `session.secret`
2. **密码安全** - 当前使用 MD5+盐值，生产环境建议使用 bcrypt
3. **数据库连接** - 确保数据库配置正确
4. **跨域问题** - 如需前端对接，请添加 CORS 中间件

## 部署

### 方式一：本地编译运行

#### 编译
```bash
make build
# 或
go build -o user-center cmd/main.go
```

#### 运行
```bash
./user-center
```

### 方式二：Docker Compose 一键部署（推荐）

#### 使用部署脚本（最简单）
```bash
chmod +x deploy.sh
./deploy.sh
```

部署脚本会自动：
- ✅ 检查 Docker 环境
- ✅ 清理旧容器
- ✅ 构建镜像
- ✅ 启动 MySQL 和应用
- ✅ 等待服务就绪
- ✅ 显示访问信息

#### 手动部署
```bash
# 启动所有服务（包括 MySQL）
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看应用日志
docker-compose logs -f app

# 查看 MySQL 日志
docker-compose logs -f mysql

# 停止服务
docker-compose down

# 停止服务并删除数据卷
docker-compose down -v
```

#### 使用 Makefile（推荐）
```bash
# 查看所有可用命令
make help

# 启动 Docker 服务
make docker-up

# 查看服务状态
make docker-status

# 查看日志
make docker-logs

# 停止服务
make docker-down

# 重启服务
make docker-restart

# 进入应用容器
make docker-exec

# 进入 MySQL 容器
make docker-mysql
```

### 方式三：仅构建 Docker 镜像

```bash
# 构建镜像
docker build -t user-center:latest .

# 运行容器（需要外部 MySQL）
docker run -d \
  --name user-center-app \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/root/config.yaml:ro \
  user-center:latest
```

## Docker 部署说明

### 架构
```
┌─────────────────────────────────────────┐
│         Docker Compose 环境              │
│                                          │
│  ┌──────────────┐    ┌──────────────┐  │
│  │              │    │              │  │
│  │   Go App     │───▶│    MySQL     │  │
│  │  (Port 8080) │    │  (Port 3306) │  │
│  │              │    │              │  │
│  └──────────────┘    └──────────────┘  │
│         │                    │          │
└─────────┼────────────────────┼──────────┘
          │                    │
          ▼                    ▼
    宿主机:8080           宿主机:3306
```

### 特性
- ✅ **多阶段构建** - 最小化镜像大小（约 20MB）
- ✅ **健康检查** - 自动检测服务状态
- ✅ **数据持久化** - MySQL 数据卷挂载
- ✅ **自动初始化** - 启动时自动执行 SQL 脚本
- ✅ **依赖管理** - 等待 MySQL 就绪后再启动应用
- ✅ **日志管理** - 日志目录挂载到宿主机
- ✅ **网络隔离** - 独立的 Docker 网络

### 环境变量

可以通过修改 `docker-compose.yml` 自定义配置：

```yaml
services:
  mysql:
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: your_database
```

### 端口说明

| 服务 | 容器端口 | 宿主机端口 | 说明 |
|------|----------|------------|------|
| Go App | 8080 | 8080 | HTTP API 服务 |
| MySQL | 3306 | 3306 | 数据库服务 |

### 数据持久化

MySQL 数据存储在 Docker volume `mysql_data` 中，即使删除容器数据也不会丢失。

如需完全清理：
```bash
docker-compose down -v  # 删除容器和数据卷
```

## 原项目链接

- [原 Java 版本](https://github.com/liyupi/user-center-backend)
- [编程导航知识星球](https://yupi.icu)

## License

本项目基于原项目改写，请遵守原项目的版权声明。
