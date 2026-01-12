# Docker 快速部署指南

## 一键部署

### 1. 使用部署脚本（推荐）
```bash
cd go-version
chmod +x deploy.sh
./deploy.sh
```

### 2. 使用 Makefile
```bash
cd go-version
make docker-up
```

### 3. 使用 Docker Compose
```bash
cd go-version
docker-compose up -d --build
```

## 服务访问

部署成功后，访问：
- **应用**: http://localhost:8080
- **健康检查**: http://localhost:8080/api/health

## 测试 API

### 注册用户
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

### 登录
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

## 常用命令

### 查看服务状态
```bash
docker-compose ps
# 或
make docker-status
```

### 查看日志
```bash
# 应用日志
docker-compose logs -f app
# 或
make docker-logs

# MySQL 日志
docker-compose logs -f mysql
```

### 进入容器
```bash
# 进入应用容器
docker exec -it user-center-app sh
# 或
make docker-exec

# 进入 MySQL
docker exec -it user-center-mysql mysql -uroot -p123456 yupi
# 或
make docker-mysql
```

### 重启服务
```bash
docker-compose restart
# 或
make docker-restart
```

### 停止服务
```bash
# 仅停止容器
docker-compose down
# 或
make docker-down

# 停止并删除数据卷（会清空数据库）
docker-compose down -v
# 或
make docker-down-v
```

## 故障排查

### 应用无法启动
```bash
# 查看应用日志
docker logs user-center-app

# 检查 MySQL 是否就绪
docker exec user-center-mysql mysqladmin ping -h localhost -uroot -p123456
```

### MySQL 连接失败
```bash
# 检查 MySQL 状态
docker-compose ps mysql

# 查看 MySQL 日志
docker logs user-center-mysql

# 手动测试连接
docker exec -it user-center-mysql mysql -uroot -p123456
```

### 端口冲突
如果 8080 或 3306 端口被占用，修改 `docker-compose.yml`:
```yaml
services:
  mysql:
    ports:
      - "3307:3306"  # 修改为其他端口
  
  app:
    ports:
      - "8081:8080"  # 修改为其他端口
```

## 配置说明

### MySQL 配置
在 `docker-compose.yml` 中修改：
```yaml
services:
  mysql:
    environment:
      MYSQL_ROOT_PASSWORD: 你的密码
      MYSQL_DATABASE: 数据库名
```

### 应用配置
修改 `config.docker.yaml`（不要修改 `config.yaml`）

### 数据持久化
MySQL 数据保存在 Docker volume 中，删除容器不会丢失数据。

查看数据卷：
```bash
docker volume ls | grep user-center
```

备份数据：
```bash
docker exec user-center-mysql mysqldump -uroot -p123456 yupi > backup.sql
```

恢复数据：
```bash
docker exec -i user-center-mysql mysql -uroot -p123456 yupi < backup.sql
```

## 生产环境建议

1. **修改密码**: 更改 MySQL root 密码和 Session secret
2. **使用外部数据库**: 建议使用云数据库服务
3. **配置日志**: 添加日志收集和监控
4. **HTTPS**: 使用 Nginx 反向代理配置 SSL
5. **资源限制**: 在 docker-compose.yml 中添加资源限制

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
```
