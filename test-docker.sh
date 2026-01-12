#!/bin/bash

# Docker 部署测试脚本

echo "======================================"
echo "Docker 部署测试"
echo "======================================"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 基础 URL
BASE_URL="http://localhost:8080/api"

# 测试健康检查
test_health() {
    echo -n "测试健康检查... "
    response=$(curl -s "${BASE_URL}/health")
    if echo "$response" | grep -q "UP"; then
        echo -e "${GREEN}✓ 通过${NC}"
        return 0
    else
        echo -e "${RED}✗ 失败${NC}"
        return 1
    fi
}

# 测试用户注册
test_register() {
    echo -n "测试用户注册... "
    timestamp=$(date +%s)
    response=$(curl -s -X POST "${BASE_URL}/user/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"userAccount\": \"test${timestamp}\",
            \"userPassword\": \"12345678\",
            \"checkPassword\": \"12345678\",
            \"planetCode\": \"${timestamp}\"
        }")
    
    if echo "$response" | grep -q '"code":0'; then
        echo -e "${GREEN}✓ 通过${NC}"
        echo "  账号: test${timestamp}"
        return 0
    else
        echo -e "${RED}✗ 失败${NC}"
        echo "  响应: $response"
        return 1
    fi
}

# 测试用户登录
test_login() {
    echo -n "测试用户登录... "
    timestamp=$(date +%s)
    
    # 先注册
    curl -s -X POST "${BASE_URL}/user/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"userAccount\": \"login${timestamp}\",
            \"userPassword\": \"12345678\",
            \"checkPassword\": \"12345678\",
            \"planetCode\": \"${timestamp}\"
        }" > /dev/null
    
    # 再登录
    response=$(curl -s -X POST "${BASE_URL}/user/login" \
        -H "Content-Type: application/json" \
        -c /tmp/cookie_${timestamp}.txt \
        -d "{
            \"userAccount\": \"login${timestamp}\",
            \"userPassword\": \"12345678\"
        }")
    
    if echo "$response" | grep -q '"code":0'; then
        echo -e "${GREEN}✓ 通过${NC}"
        
        # 测试获取当前用户
        echo -n "测试获取当前用户... "
        current_response=$(curl -s -X GET "${BASE_URL}/user/current" \
            -b /tmp/cookie_${timestamp}.txt)
        
        if echo "$current_response" | grep -q '"code":0'; then
            echo -e "${GREEN}✓ 通过${NC}"
        else
            echo -e "${RED}✗ 失败${NC}"
        fi
        
        # 清理 cookie
        rm -f /tmp/cookie_${timestamp}.txt
        return 0
    else
        echo -e "${RED}✗ 失败${NC}"
        echo "  响应: $response"
        return 1
    fi
}

# 主流程
echo ""
echo "等待服务启动（5秒）..."
sleep 5

echo ""
test_health
echo ""
test_register
echo ""
test_login

echo ""
echo "======================================"
echo "测试完成"
echo "======================================"
