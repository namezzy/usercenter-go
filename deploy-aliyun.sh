#!/bin/bash
# 阿里云一键部署脚本

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 用户中心 - 阿里云一键部署"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 检查是否为 root 用户
if [ "$EUID" -ne 0 ]; then 
    echo "❌ 请使用 root 用户运行此脚本"
    echo "   sudo su"
    exit 1
fi

# 1. 更新系统
echo "📦 更新系统..."
apt update && apt upgrade -y

# 2. 安装 Docker
if ! command -v docker &> /dev/null; then
    echo "📦 安装 Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
else
    echo "✅ Docker 已安装"
fi

# 3. 安装 Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "📦 安装 Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
else
    echo "✅ Docker Compose 已安装"
fi

# 4. 安装 Git
if ! command -v git &> /dev/null; then
    echo "📦 安装 Git..."
    apt install git -y
else
    echo "✅ Git 已安装"
fi

# 5. 获取项目
cd /opt
if [ -d "user-center" ]; then
    echo "📥 更新项目..."
    cd user-center
    git pull
else
    echo "📥 克隆项目..."
    read -p "请输入 Git 仓库地址: " REPO_URL
    git clone $REPO_URL user-center
    cd user-center
fi

cd go-version

# 6. 生成随机密码
echo "🔐 生成安全密码..."
DB_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
REDIS_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
SESSION_SECRET=$(openssl rand -base64 48 | tr -d "=+/" | cut -c1-40)

# 7. 创建配置文件
echo "⚙️ 生成配置文件..."
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
  password: $REDIS_PASSWORD
  db: 0

session:
  secret: $SESSION_SECRET
  timeout: 86400
EOF

# 8. 修改 docker-compose.yml
echo "⚙️ 配置 Docker Compose..."
sed -i "s/MYSQL_ROOT_PASSWORD: 123456/MYSQL_ROOT_PASSWORD: $DB_PASSWORD/" docker-compose.yml
sed -i "s/command: redis-server --appendonly yes/command: redis-server --appendonly yes --requirepass $REDIS_PASSWORD/" docker-compose.yml

# 9. 停止旧服务
if docker-compose ps | grep -q "Up"; then
    echo "🛑 停止旧服务..."
    docker-compose down
fi

# 10. 启动服务
echo "🚀 启动服务..."
docker-compose up -d --build

# 11. 等待服务启动
echo "⏳ 等待服务启动（30秒）..."
sleep 30

# 12. 检查服务状态
echo ""
echo "✅ 检查服务状态..."
docker-compose ps

# 13. 获echoIP
SERVER_IP=$(curl -s ifconfig.me || hostname -I | awk '{print $1}')

# 14. 测试健康检查
echo ""
echo "🧪 测试服务..."
HEALTH_CHECK=$(curl -s http://localhost:8080/api/health || echo "failed")
if echo "$HEALTH_CHECK" | grep -q "UP"; then
    echo "✅ 服务运行正常"
else
    echo "⚠️ 服务可能未正常启动，请检查日志: docker-compose logs"
fi

# 15. 输出信息
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 部署完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📍 访问地址:"
echo "   http://$SERVER_IP:8080"
echo "   http://$SERVER_IP:8080/api/health"
echo ""
echo "🔐 数据库信息:"
echo "   MySQL 密码: $DB_PASSWORD"
echo "   Redis 密码: $REDIS_PASSWORD"
echo "   Session 密钥: $SESSION_SECRET"
echo ""
echo "⚠️ 重要："
echo "   1. 请保存好上述密码！"
echo "   2. 配置安全组，开放端口: 80, 443, 8080"
 HTTPS"
echo ""
echo "📝 常用命令:"
echo "   查看日志: docker-compose logs -f"
echo "   重启服务: docker-compose restart"
echo "   停止服务: docker-compose down"
echo ""
echo "📚 详细文档: https://github.com/YOUR_REPO/README.md"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 16. 保存密码到文件
cat > /opt/user-center/credentials.txt << EOF
:::: $(date)
MySQL 密码: $DB_PASSWORD
Redis 密码: $REDIS_PASSWORD
Session 密钥: $SESSION_SECRET
EOF
chmod 600 /opt/user-center/credentials.txt

echo ""
echo "💾 密码已保存到: /opt/user-center/credentials.txt"
echo ""
