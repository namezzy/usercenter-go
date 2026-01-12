#!/bin/bash

# 用户中心 Docker 部署脚本

set -e

echo "======================================"
echo "用户中心 - Docker 部署"
echo "======================================"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}错误: Docker 未安装${NC}"
        echo "请先安装 Docker: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${RED}错误: Docker Compose 未安装${NC}"
        echo "请先安装 Docker Compose: https://docs.docker.com/compose/install/"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Docker 环境检查通过${NC}"
}

# 停止并删除旧容器
cleanup() {
    echo -e "${YELLOW}清理旧容器...${NC}"
    docker-compose down -v 2>/dev/null || docker compose down -v 2>/dev/null || true
    echo -e "${GREEN}✓ 清理完成${NC}"
}

# 构建并启动服务
start_services() {
    echo -e "${YELLOW}构建并启动服务...${NC}"
    
    # 尝试使用 docker-compose 或 docker compose
    if command -v docker-compose &> /dev/null; then
        docker-compose up -d --build
    else
        docker compose up -d --build
    fi
    
    echo -e "${GREEN}✓ 服务启动成功${NC}"
}

# 等待服务就绪
wait_for_services() {
    echo -e "${YELLOW}等待服务就绪...${NC}"
    
    echo -n "等待 MySQL 启动"
    for i in {1..30}; do
        if docker exec user-center-mysql mysqladmin ping -h localhost -uroot -p123456 --silent 2>/dev/null; then
            echo -e "\n${GREEN}✓ MySQL 已就绪${NC}"
            break
        fi
        echo -n "."
        sleep 2
    done
    
    echo -n "等待应用启动"
    for i in {1..30}; do
        if docker logs user-center-app 2>&1 | grep -q "Server starting"; then
            echo -e "\n${GREEN}✓ 应用已启动${NC}"
            break
        fi
        echo -n "."
        sleep 2
    done
}

# 显示服务状态
show_status() {
    echo ""
    echo "======================================"
    echo "服务状态"
    echo "======================================"
    
    if command -v docker-compose &> /dev/null; then
        docker-compose ps
    else
        docker compose ps
    fi
    
    echo ""
    echo "======================================"
    echo "访问信息"
    echo "======================================"
    echo -e "应用地址: ${GREEN}http://localhost:8080${NC}"
    echo -e "API 前缀: ${GREEN}/api${NC}"
    echo ""
    echo "测试命令:"
    echo -e "${YELLOW}curl -X POST http://localhost:8080/api/user/register -H \"Content-Type: application/json\" -d '{\"userAccount\":\"test123\",\"userPassword\":\"12345678\",\"checkPassword\":\"12345678\",\"planetCode\":\"12345\"}'${NC}"
    echo ""
    echo "查看日志:"
    echo -e "${YELLOW}docker logs -f user-center-app${NC}"
    echo ""
    echo "进入容器:"
    echo -e "${YELLOW}docker exec -it user-center-app sh${NC}"
    echo ""
    echo "停止服务:"
    echo -e "${YELLOW}docker-compose down${NC}"
    echo "======================================"
}

# 主流程
main() {
    check_docker
    
    # 询问是否清理旧容器
    read -p "是否清理旧容器? (y/n) [y]: " clean
    clean=${clean:-y}
    if [[ $clean =~ ^[Yy]$ ]]; then
        cleanup
    fi
    
    start_services
    wait_for_services
    show_status
}

# 执行主流程
main
