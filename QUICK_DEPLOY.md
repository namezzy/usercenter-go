# 快速部署到云平台 - 5分钟上线

## 🚀 最快部署方案

### 方案一：Railway（推荐新手）⭐

**优势**: 免费、自动部署、无需配置

#### 1. 准备 GitHub 仓库
```bash
cd go-version
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/YOUR_USERNAME/user-center.git
git push -u origin main
```

#### 2. 部署到 Railway
1. 访问 https://railway.app
2. 点击 "Start a New Project"
3. 选择 "Deploy from GitHub repo"
4. 授权并选择你的仓库
5. Railway 自动检测 Dockerfile 并部署

#### 3. 添加数据库（在 Railway Dashboard）
```
点击项目 → New → Database → MySQL
点击项目 → New → Database → Redis
```

#### 4. 配置环境变量
```bash
# 在 Railway 项目设置中添加
DB_HOST: ${{MySQL.PRIVATE_URL}}
REDIS_HOST: ${{Redis.PRIVATE_URL}}
```

#### 5. 访问应用
```
Railway 自动生成域名: https://your-app.railway.app
```

**总时间**: ⏱️ **5分钟**  
**费用**: 💰 **$0-5/月**

---

### 方案二：Render（最简单免费方案）

**优势**: 完全免费、自动 HTTPS

#### 1. 推送到 GitHub（同上）

#### 2. 部署到 Render
1. 访问 https://render.com
2. 点击 "New +" → "Web Service"
3. 连接 GitHub 仓库
4. 配置:
   ```
   Name: user-center
   Environment: Docker
   Docker Command: 留空（使用 Dockerfile CMD）
   ```

#### 3. 添加数据库
```
Dashboard → New + → MySQL
Dashboard → New + → Redis
```

#### 4. 连接数据库
```bash
# 在 Web Service 环境变量中
DATABASE_URL: (从 MySQL 服务复制)
REDIS_URL: (从 Redis 服务复制)
```

**总时间**: ⏱️ **5-10分钟**  
**费用**: 💰 **$0/月**（免费层）

---

### 方案三：阿里云（国内推荐）

**优势**: 国内访问快、稳定

#### 1. 购买服务器
```
产品: 轻量应用服务器
配置: 2核2G（¥60/月）
系统: Ubuntu 22.04
地域: 选择离你最近的
```

#### 2. 连接服务器
```bash
ssh root@YOUR_SERVER_IP
```

#### 3. 一键安装脚本
```bash
# 下载安装脚本
wget https://raw.githubusercontent.com/YOUR_REPO/deploy-aliyun.sh
chmod +x deploy-aliyun.sh

# 执行安装
./deploy-aliyun.sh
```

#### 4. 安全组配置
```
控制台 → 防火墙 → 添加规则
开放端口: 80, 443, 8080
```

**总时间**: ⏱️ **15分钟**  
**费用**: 💰 **¥60/月**

---

## 📦 创建一键部署脚本

### Railway 部署配置
```bash
# 创建 railway.json
cat > railway.json << 'EOF'
{
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "Dockerfile"
  },
  "deploy": {
    "numReplicas": 1,
    "restartPolicyType": "ON_FAILURE"
  }
}
EOF
```

### Render 部署配置
```bash
# 创建 render.yaml
cat > render.yaml << 'EOF'
services:
  - type: web
    name: user-center
    env: docker
    dockerfilePath: ./Dockerfile
    plan: free
    healthCheckPath: /api/health
    
databases:
  - name: user-center-db
    databaseName: yupi
    plan: free

  - name: user-center-redis
    plan: free
EOF
```

### 阿里云一键部署脚本
```bash
cat > deploy-aliyun.sh << 'EOFSCRIPT'
#!/bin/bash
set -e

echo "🚀 开始部署用户中心到阿里云..."

# 1. 安装 Docker
echo "📦 安装 Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# 2. 安装 Docker Compose
echo "📦 安装 Docker Compose..."
curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# 3. 克隆项目
echo "📥 下载项目..."
cd /opt
git clone https://github.com/YOUR_USERNAME/user-center.git
cd user-center/go-version

# 4. 生成随机密码
DB_PASSWORD=$(openssl rand -base64 32)
SESSION_SECRET=$(openssl rand -base64 32)

# 5. 修改配置
echo "⚙️ 配置系统..."
cat > config.docker.yaml << EOF
server:
  port: 8080
  context_path: /api

database:
  host: mysql
  port: 3306
  database: yupi
  username: root
  password: $DB_PASSWORD

redis:
  host: redis
  port: 6379
  password: ""
  db: 0

session:
  secret: $SESSION_SECRET
  timeout: 86400
EOF

# 6. 修改 docker-compose.yml 密码
sed -i "s/MYSQL_ROOT_PASSWORD: 123456/MYSQL_ROOT_PASSWORD: $DB_PASSWORD/" docker-compose.yml

# 7. 启动服务
echo "🚀 启动服务..."
docker-compose up -d --build

# 8. 等待服务启动
echo "⏳ 等待服务启动..."
sleep 30

# 9. 检查状态
echo "✅ 检查服务状态..."
docker-compose ps

# 10. 输出访问信息
SERVER_IP=$(curl -s ifconfig.me)
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 部署完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "访问地址: http://$SERVER_IP:8080"
echo "健康检查: http://$SERVER_IP:8080/api/health"
echo ""
echo "数据库密码: $DB_PASSWORD"
echo "Session密钥: $SESSION_SECRET"
echo ""
echo "⚠️ 请保存好上述密码！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
EOFSCRIPT

chmod +x deploy-aliyun.sh
```

---

## 🌐 配置域名（可选）

### 1. 购买域名
- 阿里云: https://wanwang.aliyun.com
- Cloudflare: https://www.cloudflare.com
- GoDaddy: https://www.godaddy.com

### 2. 添加 DNS 解析
```
记录类型: A
主机记录: @  (或 www)
记录值: 你的服务器IP
TTL: 600
```

### 3. 配置 Nginx 反向代理
```bash
# 安装 Nginx
apt install nginx -y

# 配置
cat > /etc/nginx/sites-available/user-center << 'EOF'
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

# 启用配置
ln -s /etc/nginx/sites-available/user-center /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

### 4. 配置 HTTPS
```bash
# 安装 Certbot
apt install certbot python3-certbot-nginx -y

# 申请证书
certbot --nginx -d your-domain.com

# 自动续期
certbot renew --dry-run
```

---

## 🔍 测试部署

### 健康检查
```bash
curl http://your-domain.com/api/health
```

预期响应：
```json
{
  "status": "UP",
  "message": "User Center Service is running",
  "components": {
    "database": "UP",
    "redis": "UP"
  }
}
```

### 测试用户注册
```bash
curl -X POST http://your-domain.com/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "userAccount": "testuser",
    "userPassword": "12345678",
    "checkPassword": "12345678",
    "planetCode": "12345"
  }'
```

---

## 📊 平台对比速查表

| 特性 | Railway | Render | 阿里云 | 腾讯云 |
|------|---------|--------|--------|--------|
| 部署难度 | ⭐ | ⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| 免费额度 | ✅ | ✅ | ❌ | ❌ |
| 自动部署 | ✅ | ✅ | ❌ | ❌ |
| 国内速度 | 一般 | 较慢 | 很快 | 很快 |
| 自定义域名 | ✅ | ✅ | ✅ | ✅ |
| 数据库 | 内置 | 内置 | 需单独购买 | 需单独购买 |
| 价格/月 | $5-20 | $0-7 | ¥60+ | ¥50+ |

---

## 💡 推荐选择

### 个人学习/演示
→ **Render** (完全免费)

### 个人项目/快速原型
→ **Railway** (部署简单，$5/月)

### 国内生产环境
→ **阿里云**或**腾讯云** (稳定快速)

### 国际化项目
→ **AWS**或**Fly.io** (全球部署)

---

## 🆘 遇到问题？

### 常见问题

**Q: Railway 部署失败？**
```bash
# 检查 Dockerfile 是否正确
# 确保 go.mod 和 go.sum 存在
# 查看 Railway 日志
```

**Q: 数据库连接失败？**
```bash
# 检查数据库配置
# 确保数据库已启动
# 检查网络连接
```

**Q: 端口访问不了？**
```bash
# 检查安全组/防火墙
# 确保端口已开放
# 检查应用是否启动
```

### 获取帮助
- 查看详细文档: [CLOUD_DEPLOYMENT.md](CLOUD_DEPLOYMENT.md)
- GitHub Issues
- 官方文档

---

## 🎉 部署成功后

1. ✅ 访问健康检查接口
2. ✅ 测试用户注册登录
3. ✅ 配置监控告警
4. ✅ 设置自动备份
5. ✅ 优化性能配置

---

**下一步**: 查看完整的 [云平台部署指南](CLOUD_DEPLOYMENT.md)
