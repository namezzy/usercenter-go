# 云托管平台部署指南

## 目录
- [云平台选择](#云平台选择)
- [方案一：阿里云](#方案一阿里云)
- [方案二：腾讯云](#方案二腾讯云)
- [方案三：AWS](#方案三aws)
- [方案四：Vercel + Railway](#方案四vercel--railway)
- [方案五：Render](#方案五render)
- [方案六：Fly.io](#方案六flyio)

---

## 云平台选择

### 对比表

| 平台 | 优势 | 劣势 | 适合场景 | 费用 |
|------|------|------|----------|------|
| 阿里云 | 国内访问快，文档中文 | 配置复杂 | 企业级应用 | ¥50-300/月 |
| 腾讯云 | 国内访问快，价格实惠 | 配置较复杂 | 中小型应用 | ¥40-250/月 |
| AWS | 全球部署，生态完善 | 价格较高 | 国际化应用 | $20-100/月 |
| Railway | 部署简单，自动 CI/CD | 价格较高 | 个人项目 | $5-20/月 |
| Render | 免费额度，自动部署 | 冷启动慢 | 演示项目 | 免费-$7/月 |
| Fly.io | 全球边缘部署 | 学习曲线陡 | 高性能需求 | $0-10/月 |

---

## 方案一：阿里云

### 1. 云服务器 ECS + Docker 部署（推荐）

#### 步骤 1：购买服务器
```bash
# 推荐配置
- CPU: 2核
- 内存: 4GB
- 系统盘: 40GB
- 操作系统: Ubuntu 22.04 LTS
- 带宽: 3Mbps
```

#### 步骤 2：连接服务器
```bash
ssh root@your-server-ip
```

#### 步骤 3：安装 Docker
```bash
# 更新系统
apt update && apt upgrade -y

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# 安装 Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker-compose --version
```

#### 步骤 4：上传项目
```bash
# 方式一：使用 Git
cd /opt
git clone https://your-repo.git user-center
cd user-center/go-version

# 方式二：使用 SCP
# 在本地执行
scp -r go-version root@your-server-ip:/opt/user-center
```

#### 步骤 5：配置环境
```bash
cd /opt/user-center/go-version

# 修改配置
vim config.docker.yaml

# 生产环境配置示例
server:
  port: 8080
  context_path: /api

database:
  host: mysql
  port: 3306
  database: yupi
  username: root
  password: YOUR_STRONG_PASSWORD  # 修改为强密码

redis:
  host: redis
  port: 6379
  password: YOUR_REDIS_PASSWORD   # 建议设置密码
  db: 0

session:
  secret: YOUR_SESSION_SECRET_CHANGE_ME  # 修改为随机字符串
  timeout: 86400
```

#### 步骤 6：启动服务
```bash
# 启动
./deploy.sh

# 或手动启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f
```

#### 步骤 7：配置安全组
```bash
# 在阿里云控制台添加安全组规则
# 放行端口：
- 22    (SSH)
- 80    (HTTP)
- 443   (HTTPS)
- 8080  (应用端口)
```

#### 步骤 8：配置域名（可选）
```bash
# 安装 Nginx
apt install nginx -y

# 配置反向代理
vim /etc/nginx/sites-available/user-center

# Nginx 配置
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

# 启用配置
ln -s /etc/nginx/sites-available/user-center /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

#### 步骤 9：配置 HTTPS（推荐）
```bash
# 安装 Certbot
apt install certbot python3-certbot-nginx -y

# 申请证书
certbot --nginx -d your-domain.com

# 自动续期
certbot renew --dry-run
```

---

## 方案二：腾讯云

### 使用腾讯云轻量应用服务器

#### 步骤 1：购买服务器
```bash
# 推荐配置
- CPU: 2核
- 内存: 4GB
- 系统盘: 60GB SSD
- 操作系统: Ubuntu 22.04
- 带宽: 5Mbps
```

#### 步骤 2-9：与阿里云相同
参考阿里云部署步骤。

### 额外：使用腾讯云数据库（推荐生产环境）

```bash
# 购买云数据库 MySQL
- 版本: MySQL 8.0
- 规格: 1核 1GB（起步）
- 存储: 50GB

# 购买云数据库 Redis
- 版本: Redis 7.0
- 规格: 1GB 内存
```

#### 修改配置使用云数据库
```yaml
# config.docker.yaml
database:
  host: your-mysql-instance.tencentcdb.com
  port: 3306
  database: yupi
  username: root
  password: your_password

redis:
  host: your-redis-instance.tencentcloudapi.com
  port: 6379
  password: your_redis_password
  db: 0
```

#### 仅部署应用容器
```yaml
# docker-compose.yml（仅应用）
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - TZ=Asia/Shanghai
    volumes:
      - ./config.docker.yaml:/root/config.yaml:ro
    restart: always
```

---

## 方案三：AWS

### 使用 AWS ECS + Fargate（无服务器容器）

#### 步骤 1：准备 Docker 镜像

```bash
# 构建镜像
cd go-version
docker build -t user-center:latest .

# 推送到 AWS ECR
aws ecr create-repository --repository-name user-center
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_ECR_URI
docker tag user-center:latest YOUR_ECR_URI/user-center:latest
docker push YOUR_ECR_URI/user-center:latest
```

#### 步骤 2：创建 RDS 数据库
```bash
# AWS 控制台创建
- 引擎: MySQL 8.0
- 实例类型: db.t3.micro
- 存储: 20GB
- 公开访问: 否（在 VPC 内）
```

#### 步骤 3：创建 ElastiCache Redis
```bash
# AWS 控制台创建
- 引擎: Redis 7.0
- 节点类型: cache.t3.micro
- 副本数: 0（单节点）
```

#### 步骤 4：创建 ECS 任务定义

```json
{
  "family": "user-center",
  "containerDefinitions": [
    {
      "name": "user-center",
      "image": "YOUR_ECR_URI/user-center:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {"name": "TZ", "value": "Asia/Shanghai"}
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/user-center",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ],
  "requiresCompatibilities": ["FARGATE"],
  "networkMode": "awsvpc",
  "cpu": "512",
  "memory": "1024"
}
```

#### 步骤 5：创建 ECS 服务
```bash
aws ecs create-service \
  --cluster user-center-cluster \
  --service-name user-center-service \
  --task-definition user-center \
  --desired-count 1 \
  --launch-type FARGATE
```

---

## 方案四：Railway（推荐个人项目）

### 特点
- ✅ 自动 CI/CD
- ✅ GitHub 集成
- ✅ 内置数据库
- ✅ 简单配置

### 步骤 1：准备项目

```bash
# 创建 railway.json
{
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "Dockerfile"
  },
  "deploy": {
    "numReplicas": 1,
    "sleepApplication": false,
    "restartPolicyType": "ON_FAILURE"
  }
}

# 创建 Dockerfile（使用项目现有的）
# 已经存在，无需修改
```

### 步骤 2：创建 Railway 项目

1. 访问 https://railway.app
2. 使用 GitHub 登录
3. 点击 "New Project"
4. 选择 "Deploy from GitHub repo"
5. 选择你的仓库

### 步骤 3：添加数据库服务

```bash
# 在 Railway Dashboard
1. 点击 "New" -> "Database" -> "Add MySQL"
2. 点击 "New" -> "Database" -> "Add Redis"
```

### 步骤 4：配置环境变量

```bash
# 在 Railway 应用中添加环境变量
DB_HOST=${{MySQL.RAILWAY_PRIVATE_DOMAIN}}
DB_PORT=${{MySQL.RAILWAY_TCP_PROXY_PORT}}
DB_NAME=railway
DB_USER=root
DB_PASSWORD=${{MySQL.MYSQL_ROOT_PASSWORD}}

REDIS_HOST=${{Redis.RAILWAY_PRIVATE_DOMAIN}}
REDIS_PORT=${{Redis.RAILWAY_TCP_PROXY_PORT}}
```

### 步骤 5：修改配置读取方式

```go
// 在 cmd/main.go 中添加环境变量支持
import "os"

func loadConfig() *config.Config {
    cfg, err := config.LoadConfig("config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 覆盖为环境变量
    if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
        cfg.Database.Host = dbHost
    }
    if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
        port, _ := strconv.Atoi(dbPort)
        cfg.Database.Port = port
    }
    // ... 其他环境变量
    
    return cfg
}
```

### 步骤 6：部署
```bash
# 推送代码到 GitHub
git add .
git commit -m "Add Railway deployment"
git push origin main

# Railway 自动检测并部署
# 几分钟后应用即可访问
```

---

## 方案五：Render（免费额度）

### 特点
- ✅ 免费层（有限制）
- ✅ 自动 HTTPS
- ✅ 自动部署
- ⚠️ 冷启动慢

### 步骤 1：准备项目

```yaml
# 创建 render.yaml
services:
  - type: web
    name: user-center
    env: docker
    dockerfilePath: ./Dockerfile
    dockerContext: ./go-version
    plan: free
    envVars:
      - key: DATABASE_URL
        fromDatabase:
          name: user-center-db
          property: connectionString
      - key: REDIS_URL
        fromService:
          name: user-center-redis
          type: redis
          property: connectionString

databases:
  - name: user-center-db
    databaseName: yupi
    user: root
    plan: free

  - name: user-center-redis
    plan: free
```

### 步骤 2：部署

1. 访问 https://render.com
2. 连接 GitHub 仓库
3. 选择 "New" -> "Blueprint"
4. Render 自动读取 render.yaml 并部署

---

## 方案六：Fly.io（全球边缘部署）

### 特点
- ✅ 全球多区域部署
- ✅ 自动扩缩容
- ✅ 永久免费额度

### 步骤 1：安装 Fly CLI

```bash
# macOS/Linux
curl -L https://fly.io/install.sh | sh

# Windows
powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"

# 登录
fly auth login
```

### 步骤 2：初始化项目

```bash
cd go-version

# 初始化
fly launch

# 回答问题
? Choose an app name: user-center
? Choose a region: hkg (Hong Kong)
? Would you like to set up a MySQL database? Yes
? Would you like to set up a Redis database? Yes
```

### 步骤 3：配置生成的 fly.toml

```toml
# fly.toml
app = "user-center"
primary_region = "hkg"

[build]
  dockerfile = "Dockerfile"

[env]
  PORT = "8080"

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [[services.ports]]
    port = 80
    handlers = ["http"]

  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]

[services.concurrency]
  type = "connections"
  hard_limit = 25
  soft_limit = 20

[[services.tcp_checks]]
  interval = "15s"
  timeout = "2s"
  grace_period = "5s"
```

### 步骤 4：部署

```bash
# 部署
fly deploy

# 查看状态
fly status

# 查看日志
fly logs

# 打开应用
fly open
```

---

## 通用优化建议

### 1. 环境变量管理

```bash
# 创建 .env.production
DB_HOST=your-db-host
DB_PORT=3306
DB_NAME=yupi
DB_USER=root
DB_PASSWORD=strong-password

REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=redis-password

SESSION_SECRET=random-secret-key-change-me
```

### 2. 健康检查

```go
// 在 main.go 中已有
r.GET("/api/health", healthCheckHandler)
```

### 3. 日志管理

```bash
# 使用集中式日志
- 阿里云: SLS 日志服务
- AWS: CloudWatch Logs
- 腾讯云: CLS 日志服务
```

### 4. 监控告警

```bash
# 推荐工具
- Prometheus + Grafana
- 云平台自带监控
- UptimeRobot（免费）
```

### 5. 自动备份

```bash
# 数据库备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec user-center-mysql mysqldump -uroot -p123456 yupi > backup_$DATE.sql

# 添加到 crontab
0 2 * * * /path/to/backup.sh
```

---

## 成本估算

### 个人/演示项目
- **Render（免费）**: $0/月
- **Railway（免费额度）**: $0-5/月
- **Fly.io**: $0-5/月

### 小型生产项目
- **阿里云轻量**: ¥50-100/月
- **腾讯云轻量**: ¥40-80/月
- **Railway**: $10-20/月

### 中型生产项目
- **阿里云 ECS**: ¥200-500/月
- **AWS EC2**: $30-100/月
- **腾讯云**: ¥180-400/月

---

## 常见问题

### Q1: 如何选择云平台？
**A**: 
- 国内用户 -> 阿里云/腾讯云
- 国际用户 -> AWS/GCP
- 个人项目 -> Railway/Render
- 快速原型 -> Fly.io

### Q2: 需要备案吗？
**A**: 国内服务器（阿里云、腾讯云）需要备案，国外服务器不需要。

### Q3: 如何实现高可用？
**A**: 
- 使用云数据库（自动备份）
- 多实例部署（负载均衡）
- Redis 哨兵/集群模式
- 定期备份数据

### Q4: 如何降低成本？
**A**: 
- 使用按量付费
- 竞价实例（AWS Spot）
- 轻量应用服务器
- 免费托管平台

---

## 下一步

1. **选择云平台**: 根据需求选择合适的平台
2. **准备域名**: 购买域名（可选）
3. **配置 SSL**: 使用 Let's Encrypt 免费证书
4. **监控告警**: 设置监控和告警
5. **自动备份**: 配置数据备份策略
6. **性能优化**: CDN、缓存优化
7. **安全加固**: 防火墙、安全组配置

---

## 参考链接

- [阿里云 ECS 文档](https://help.aliyun.com/product/25365.html)
- [腾讯云轻量应用服务器](https://cloud.tencent.com/product/lighthouse)
- [Railway 文档](https://docs.railway.app/)
- [Render 文档](https://render.com/docs)
- [Fly.io 文档](https://fly.io/docs/)
- [AWS ECS 文档](https://docs.aws.amazon.com/ecs/)
