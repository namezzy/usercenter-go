#!/bin/bash

echo "====================================="
echo "检查和修复数据库表结构"
echo "====================================="

# 检查 Docker 是否运行
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装或未运行"
    exit 1
fi

# 检查 MySQL 容器是否运行
if docker ps | grep -q user-center-mysql; then
    echo "✅ MySQL 容器正在运行"
    
    echo ""
    echo "检查当前表结构..."
    docker exec user-center-mysql mysql -uroot -padmin@123321 levi -e "SHOW COLUMNS FROM user;" 2>/dev/null || {
        echo "⚠️  无法查询表结构，可能表不存在"
    }
    
    echo ""
    echo "====================================="
    echo "需要重建数据库以修复列名问题"
    echo "====================================="
    echo ""
    echo "执行以下命令："
    echo "1. docker-compose down -v   # 删除容器和数据卷"
    echo "2. docker-compose up -d --build  # 重新构建"
    echo ""
    read -p "是否现在执行重建? (y/n): " confirm
    
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        echo ""
        echo "停止并删除容器和数据卷..."
        docker-compose down -v
        
        echo ""
        echo "重新构建并启动..."
        docker-compose up -d --build
        
        echo ""
        echo "等待 MySQL 就绪..."
        sleep 15
        
        echo ""
        echo "检查新的表结构..."
        docker exec user-center-mysql mysql -uroot -padmin@123321 levi -e "SHOW COLUMNS FROM user;"
        
        echo ""
        echo "✅ 重建完成！"
    fi
else
    echo "⚠️  MySQL 容器未运行"
    echo "请先启动: docker-compose up -d"
fi
