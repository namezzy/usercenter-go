# Docker 部署完整指南

## 快速开始

### 方式一：一键部署脚本 ⭐ 推荐
```bash
./deploy.sh
```

### 方式二：Makefile
```bash
make docker-up        # 启动服务
make docker-status    # 查看状态
make docker-logs      # 查看日志
```

### 方式三：Docker Compose
```bash
docker-compose up -d --build
```

## 完整部署流程

### 1. 检查环境
```bash
docker --version
docker-compose --version
```

### 2. 启
```bash
cd go-version
./deploy.sh
```

### 3. 验证部署
```bash
# 检查健康状态
curl http://localhost:8080/api/health

# 运行测试脚本
./test-docker.sh
```

### 4. 查看日志
```bash
# 应用日志
docker logs -f user-center-app

# MySQL 日志
docker logs -f user-center-mysql

# 所有日志
docker-compose logs -f
```

## 架构说明

```

              Docker Compose 环境                     │
                                                      │
 │  ┌──────────────────┐         ┌─────────────────
  │                  │  MySQL  │                  │ │
  │   Go 应用容器     │────────▶│   MySQL 容器      │ │
  │  user-center-app │   3306  │  user-center-db  │ │
  │                  │         │                  │ │
 │  └────────┬─────────┘         └──
           │                            │            │
           │ 8080                       │ 数据持久化  │
           │                            ▼            │
           │                     mysql_data volume   │
cd /root/p_website && git commit "Update page title to 'Levi's Home Page'" -m
            │
            ▼
    http://localhost:8080
```

## 服务说明

### 应用容器 (user-center-app)
- **镜像**: 多阶段构建，基于 Alpine Linux
- **大小**: ~20MB
- **端口**: 8080
- **配置**: config.docker.yaml
- **日志**: 挂载到 ./logs 目录
#- **
**: 每 30 秒检查一次

### MySQL 容器 (user-center-mysql)
- **镜像**: mysql:8.0
- **端口**: 3306
- **数据库**: yupi
- **用户**: root / 123456
- **数据持久化**: Docker volume
- **自动初始化**: 启动时执行 create_table.sql

## 配置文件

### docker-compose.yml
cd /root/p_website && git commit -m Update page title to Levis Home Page
- 服务定义
- 网络配置
- 数据'EOF'
- 健康检查
- 依赖关系

### config.docker.yaml
Docker 环境专用配置：
```yaml
database:
  host: mysql  # 使用容器服务名
```

### Dockerfile
cd /root/p_website && git commit -m "Update page title to 'Levi's Home Page'"
1. **构建阶段**: golang:1.21-alpine
   - 下载依赖
   - 编译二进制
2. **运行阶段**: alpine:latest
   - 仅包含二进制文件
   - 最小化镜像大小

## 常用命令

### 启动和停止
```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 重启服务
docker-compose restart

# 停止并删除数据
docker-compose down -v
```

### 查看状态
```bash
# 查看所有容器
docker-compose ps

# 查看应用日志
docker logs user-center-app

# 实时日志
docker logs -f user-center-app --tail 100
```

### 进入容器
```bash
# 进入应用容器
docker exec -it user-center-app sh

# 进入 MySQL
docker exec -it user-center-mysql mysql -uroot -p123456 yupi
```

### 数据库操作
```bash
# 备份数据库
docker exec user-center-mysql mysqldump -uroot -p123456 yupi > backup.sql

# 恢复数据库
docker exec -i user-center-mysql mysql -uroot -p123456 yupi < backup.sql

# 查看表
docker exec user-center-mysql mysql -uroot -p123456 -e "USE yupi; SHOW TABLES;"
```

## Make 命令速查

```bash
make help              # 查看所有命令
make docker-up         # 启动服
make docker-down       # 停止服务
make docker-down-v     # 停止并删除数据
make docker-restart    # 重启服务
make docker-logs       # 查看日志
make docker-status     # 查看状态
make docker-exec       # 进入应用容器
make docker-mysql      # 进入 MySQL
```

## API 测试

### 测试脚本
```bash
./test-docker.sh
```

### 手动测试
```bash
# 1. 健康检查
curl http://localhost:8080/api/health

# 2. 注册用户
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "userAccount": "testuser",
    "userPassword": "12345678",
    "checkPassword": "12345678",
    "planetCode": "12345"
  }'

# 3. 登录
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -c cookie.txt \
  -d '{
    "userAccount": "testuser",
    "userPassword": "12345678"
  }'

# 4. 获取当前用户
curl -X GET http://localhost:8080/api/user/current \
  -b cookie.txt
```

## 故障排查

### 应用启动失败
```bash
#  1.
docker logs user-center-app

# 2. 检查配置文件
cat config.docker.yaml

# 3. 检查 MySQL 连接
docker exec user-center-app sh -c "nc -zv mysql 3306"
```

### MySQL 连接失败
```bash
# 1. 检查 MySQL 状态
docker ps | grep mysql

# 2. 查看 MySQL 日志
docker logs user-center-mysql

# 3. 测试连接
docker exec user-center-mysql mysqladmin ping -h localhost -uroot -p123456

# 4. 检查数据库
docker exec user-center-mysql mysql -uroot -p123456 -e "SHOW DATABASES;"
```

### 端口冲突
#
cd /root/p_website && git commit "Update page title to 'Levi's Home Page'" -m         docker-compose.yml：
```yaml
services:
  mysql:
    ports:
      - "3307:3306"  # MySQL 改为 3307
  app:
    ports:
      - "8081:8080"  # 应用改为 8081
```

'EOF'         config.docker.yaml：
```yaml
database:
  port: 3307  # 如果修改了 MySQL 端口
```

### 数据丢失
 Docker volume 中，删除容器不会丢失数据。'EOF'

cd /root/p_website && git commit -m "Update page title to 'Levi's Home Page'"
```bash
docker volume ls | grep user-center
docker volume inspect go-version_mysql_data
```

### 重新初始化
```bash
# 完全清理并重新部署
docker-compose down -v
rm -rf logs/*
docker-compose up -d --build
```

## 性能优化

### 1. 镜像优化
- ✅ 多阶段构建
- ✅ Alpine Linux 基础镜像
- ✅ 编译时优化 (-ldflags="-w -s")

### 2. 数据库连
 cmd/main.go 中已配置：
```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 3. 资源限制
 docker-compose.yml 中添加：
```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

## 生产环境建议

### 1. 安全配置
```yaml
# config.docker.yaml
session:
  secret: your-strong-random-secret-key  # 修改为强密码

# docker-compose.yml
mysql:
  environment:
    MYSQL_ROOT_PASSWORD: strong-password  # 修改为强密码
```

### 2. 使用环境变量
```bash
# .env 文件
MYSQL_PASSWORD=your-strong-password
SESSION_SECRET=your-session-secret

# docker-compose.yml 中引用
environment:
  MYSQL_ROOT_PASSWORD: ${MYSQL_PASSWORD}
```

### 3. 外部数据库
cd /root/p_website && git commit -m  AWS RDS、阿里云 RDS）："Update page title to 'Levi's Home 
```yaml
# docker-compose.yml - 移除 mysql 服务
services:
  app:
    environment:
      DB_HOST: your-rds-endpoint
      DB_PASSWORD: your-password
```

### 4. 反向代理
 Nginx 作为反向代理：
```nginx
upstream user_center {
    server localhost:8080;
}

server {
    listen 80;
    server_name yourdomain.com;
    
    location /api {
        proxy_pass http://user_center;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 5. 日志管理
```yaml
# docker-compose.yml
services:
  app:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 监控和维护

### 健康检查
```bash
# 应用健康检查
curl http://localhost:8080/api/health

# Docker 容器健康状态
docker ps --format "table {{.Names}}\t{{.Status}}"
```

### 性能监控
```bash
# 容器资源使用
docker stats user-center-app user-center-mysql

# 查看容器详情
docker inspect user-center-app
```

### 定期维护
```bash
# 清理未使用的镜像
docker image prune -a

# 清理未使用的卷
docker volume prune

# 清理未使用的网络
docker network prune
```

## 版本更新

### 更新应用
```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建
docker-compose up -d --build

# 3. 验证
./test-docker.sh
```

### 更新 MySQL
```bash
# 1. 备份数据
docker exec user-center-mysql mysqldump -uroot -p123456 yupi > backup.sql

# 2. 修改版本
# 编辑 docker-compose.yml: mysql:8.1

# 3. 重启
docker-compose up -d mysql

# 4. 验证
docker logs user-center-mysql
```

## 技术支持

- 原 Java 版本: https://github.com/liyupi/user-center-backend
- Docker 官方文档: https://docs.docker.com
- Go 官方文档: https://golang.org/doc
- Gin 框架: https://gin-gonic.com
