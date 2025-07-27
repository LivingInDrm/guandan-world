#!/bin/bash

# 掼蛋在线游戏部署脚本
# 支持开发环境和生产环境部署

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
ENVIRONMENT=${1:-development}
PROJECT_NAME="guandan-world"
DOCKER_COMPOSE_FILE="docker-compose.yml"
PRODUCTION_COMPOSE_FILE="docker-compose.production.yml"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "掼蛋在线游戏部署脚本"
    echo ""
    echo "用法: $0 [ENVIRONMENT] [COMMAND]"
    echo ""
    echo "环境:"
    echo "  development  开发环境 (默认)"
    echo "  production   生产环境"
    echo ""
    echo "命令:"
    echo "  deploy       部署应用 (默认)"
    echo "  start        启动服务"
    echo "  stop         停止服务"
    echo "  restart      重启服务"
    echo "  logs         查看日志"
    echo "  status       查看状态"
    echo "  clean        清理资源"
    echo "  backup       备份数据"
    echo "  restore      恢复数据"
    echo ""
    echo "示例:"
    echo "  $0 development deploy"
    echo "  $0 production start"
    echo "  $0 production logs backend"
}

# 检查依赖
check_dependencies() {
    log_info "检查部署依赖..."

    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi

    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi

    # 检查Git
    if ! command -v git &> /dev/null; then
        log_warning "Git未安装，某些功能可能不可用"
    fi

    log_success "依赖检查通过"
}

# 设置环境变量
setup_environment() {
    log_info "设置 $ENVIRONMENT 环境变量..."

    if [ "$ENVIRONMENT" = "production" ]; then
        DOCKER_COMPOSE_FILE="$PRODUCTION_COMPOSE_FILE"
        
        # 检查生产环境必需的环境变量
        required_vars=(
            "JWT_SECRET"
            "REDIS_PASSWORD"
            "GRAFANA_PASSWORD"
        )

        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                log_warning "环境变量 $var 未设置，使用默认值"
            fi
        done

        # 创建生产环境配置文件
        if [ ! -f ".env.production" ]; then
            log_info "创建生产环境配置文件..."
            cat > .env.production << EOF
# 生产环境配置
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY=24h
CORS_ORIGINS=https://yourdomain.com
LOG_LEVEL=info
REDIS_PASSWORD=your-redis-password
GRAFANA_PASSWORD=admin123
REACT_APP_API_URL=https://yourdomain.com
REACT_APP_WS_URL=wss://yourdomain.com
EOF
            log_warning "请编辑 .env.production 文件并设置正确的配置值"
        fi
    else
        # 开发环境配置
        if [ ! -f ".env" ]; then
            log_info "创建开发环境配置文件..."
            cp .env.example .env
        fi
    fi

    log_success "环境变量设置完成"
}

# 构建镜像
build_images() {
    log_info "构建Docker镜像..."

    if [ "$ENVIRONMENT" = "production" ]; then
        # 生产环境构建
        docker-compose -f $DOCKER_COMPOSE_FILE build --no-cache
    else
        # 开发环境构建
        docker-compose build
    fi

    log_success "镜像构建完成"
}

# 部署应用
deploy_application() {
    log_info "部署 $ENVIRONMENT 环境应用..."

    # 创建必要的目录
    mkdir -p logs/{nginx,backend,frontend}
    mkdir -p config
    mkdir -p data/{redis,prometheus,grafana,loki}

    # 设置权限
    chmod -R 755 logs
    chmod -R 755 data

    # 启动服务
    if [ "$ENVIRONMENT" = "production" ]; then
        docker-compose -f $DOCKER_COMPOSE_FILE up -d
    else
        docker-compose up -d
    fi

    # 等待服务启动
    log_info "等待服务启动..."
    sleep 30

    # 健康检查
    check_health

    log_success "$ENVIRONMENT 环境部署完成"
}

# 健康检查
check_health() {
    log_info "执行健康检查..."

    services=("backend" "frontend")
    
    for service in "${services[@]}"; do
        if [ "$service" = "backend" ]; then
            url="http://localhost:8080/healthz"
        else
            url="http://localhost:3000"
        fi

        log_info "检查 $service 服务..."
        
        for i in {1..10}; do
            if curl -f -s "$url" > /dev/null; then
                log_success "$service 服务健康"
                break
            else
                if [ $i -eq 10 ]; then
                    log_error "$service 服务健康检查失败"
                    return 1
                fi
                log_info "等待 $service 服务启动... ($i/10)"
                sleep 5
            fi
        done
    done

    log_success "所有服务健康检查通过"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    docker-compose -f $DOCKER_COMPOSE_FILE start
    log_success "服务启动完成"
}

# 停止服务
stop_services() {
    log_info "停止服务..."
    docker-compose -f $DOCKER_COMPOSE_FILE stop
    log_success "服务停止完成"
}

# 重启服务
restart_services() {
    log_info "重启服务..."
    docker-compose -f $DOCKER_COMPOSE_FILE restart
    log_success "服务重启完成"
}

# 查看日志
view_logs() {
    local service=${2:-}
    if [ -n "$service" ]; then
        docker-compose -f $DOCKER_COMPOSE_FILE logs -f "$service"
    else
        docker-compose -f $DOCKER_COMPOSE_FILE logs -f
    fi
}

# 查看状态
check_status() {
    log_info "服务状态:"
    docker-compose -f $DOCKER_COMPOSE_FILE ps
    
    log_info "资源使用情况:"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"
}

# 清理资源
clean_resources() {
    log_warning "这将删除所有容器、镜像和数据卷，确定要继续吗? (y/N)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        log_info "清理Docker资源..."
        docker-compose -f $DOCKER_COMPOSE_FILE down -v --rmi all
        docker system prune -f
        log_success "资源清理完成"
    else
        log_info "取消清理操作"
    fi
}

# 备份数据
backup_data() {
    log_info "备份数据..."
    
    backup_dir="backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$backup_dir"
    
    # 备份Redis数据
    if docker-compose -f $DOCKER_COMPOSE_FILE ps redis | grep -q "Up"; then
        docker exec guandan-redis redis-cli --rdb /data/backup.rdb
        docker cp guandan-redis:/data/backup.rdb "$backup_dir/redis.rdb"
        log_success "Redis数据备份完成"
    fi
    
    # 备份日志
    cp -r logs "$backup_dir/"
    
    # 备份配置
    cp -r config "$backup_dir/"
    
    log_success "数据备份完成: $backup_dir"
}

# 恢复数据
restore_data() {
    local backup_dir=$2
    
    if [ -z "$backup_dir" ]; then
        log_error "请指定备份目录"
        exit 1
    fi
    
    if [ ! -d "$backup_dir" ]; then
        log_error "备份目录不存在: $backup_dir"
        exit 1
    fi
    
    log_info "从 $backup_dir 恢复数据..."
    
    # 恢复Redis数据
    if [ -f "$backup_dir/redis.rdb" ]; then
        docker cp "$backup_dir/redis.rdb" guandan-redis:/data/dump.rdb
        docker-compose -f $DOCKER_COMPOSE_FILE restart redis
        log_success "Redis数据恢复完成"
    fi
    
    log_success "数据恢复完成"
}

# 显示部署信息
show_deployment_info() {
    log_success "部署完成！"
    echo ""
    echo "🎮 掼蛋在线游戏已启动"
    echo "================================"
    echo "前端地址: http://localhost:3000"
    echo "后端API: http://localhost:8080"
    echo "WebSocket: ws://localhost:8080/ws"
    
    if [ "$ENVIRONMENT" = "production" ]; then
        echo ""
        echo "📊 监控面板:"
        echo "Grafana: http://localhost:3001 (admin/admin123)"
        echo "Prometheus: http://localhost:9090"
    fi
    
    echo ""
    echo "📋 常用命令:"
    echo "查看状态: $0 $ENVIRONMENT status"
    echo "查看日志: $0 $ENVIRONMENT logs [service]"
    echo "重启服务: $0 $ENVIRONMENT restart"
    echo "停止服务: $0 $ENVIRONMENT stop"
}

# 主函数
main() {
    local command=${2:-deploy}
    
    case $command in
        "help"|"-h"|"--help")
            show_help
            exit 0
            ;;
        "deploy")
            check_dependencies
            setup_environment
            build_images
            deploy_application
            show_deployment_info
            ;;
        "start")
            start_services
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            restart_services
            ;;
        "logs")
            view_logs "$@"
            ;;
        "status")
            check_status
            ;;
        "clean")
            clean_resources
            ;;
        "backup")
            backup_data
            ;;
        "restore")
            restore_data "$@"
            ;;
        *)
            log_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"