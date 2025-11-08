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
MODE=${1:-dev}
PROJECT_NAME="guandan-world"

# 根据模式选择配置文件
case $MODE in
    "dev"|"development")
        MODE="dev"
        DOCKER_COMPOSE_FILE="docker-compose.dev.yml"
        ;;
    "prod"|"production")
        MODE="prod"
        DOCKER_COMPOSE_FILE="docker-compose.production.yml"
        ;;
    *)
        echo "错误: 未知模式 '$MODE'"
        echo "请使用: dev 或 prod"
        exit 1
        ;;
esac

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
    echo "用法: $0 [MODE] [COMMAND]"
    echo ""
    echo "模式:"
    echo "  dev          开发模式 (默认) - 后端Docker + 前端本地"
    echo "  prod         生产模式 - 优化构建 + 监控栈"
    echo ""
    echo "命令:"
    echo "  deploy       部署应用 (默认)"
    echo "  start        启动服务"
    echo "  stop         停止服务"
    echo "  restart      重启服务"
    echo "  rebuild      重新构建镜像"
    echo "  logs         查看日志"
    echo "  status       查看状态"
    echo "  health       健康检查"
    echo "  clean        清理资源"
    echo "  backup       备份数据 (仅生产)"
    echo "  restore      恢复数据 (仅生产)"
    echo ""
    echo "示例:"
    echo "  $0 dev deploy          # 启动开发环境"
    echo "  $0 dev logs backend    # 查看后端日志"
    echo "  $0 prod deploy         # 部署生产环境"
    echo "  $0 prod backup         # 备份生产数据"
    echo ""
    echo "重要提示:"
    echo "  - dev模式: 后端在Docker中，前端需本地运行 (npm run dev)"
    echo "  - prod模式: 所有服务都在Docker中运行"
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
    log_info "设置 $MODE 模式环境变量..."

    if [ "$MODE" = "prod" ]; then
        # 检查生产环境必需的环境变量
        if [ ! -f ".env.production" ]; then
            log_info "创建生产环境配置文件..."
            cat > .env.production << EOF
# 生产环境配置
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY=24h
CORS_ORIGINS=https://yourdomain.com
LOG_LEVEL=info
REDIS_PASSWORD=your-redis-password
POSTGRES_PASSWORD=your-postgres-password
GRAFANA_PASSWORD=admin123
REACT_APP_API_URL=https://yourdomain.com
REACT_APP_WS_URL=wss://yourdomain.com
EOF
            log_warning "⚠️  请编辑 .env.production 文件并设置正确的配置值"
            log_warning "生产环境必须修改所有密码和密钥！"
        fi
    else
        log_info "开发模式使用默认配置，无需额外设置"
    fi

    log_success "环境变量设置完成"
}

# 构建镜像
build_images() {
    log_info "构建 $MODE 模式 Docker 镜像..."

    if [ "$MODE" = "prod" ]; then
        # 生产环境构建（无缓存，确保最新）
        docker-compose -f $DOCKER_COMPOSE_FILE build --no-cache
    else
        # 开发环境构建（使用缓存，加快速度）
        docker-compose -f $DOCKER_COMPOSE_FILE build
    fi

    log_success "镜像构建完成"
}

# 重新构建镜像
rebuild_images() {
    log_info "重新构建 $MODE 模式镜像..."
    docker-compose -f $DOCKER_COMPOSE_FILE build --no-cache
    log_success "镜像重建完成"
}

# 部署应用
deploy_application() {
    log_info "部署 $MODE 模式应用..."

    # 创建必要的目录
    mkdir -p logs/{backend,frontend}
    
    if [ "$MODE" = "prod" ]; then
        mkdir -p logs/nginx
        mkdir -p data/{prometheus,grafana,loki}
        chmod -R 755 data
    fi

    chmod -R 755 logs

    # 启动服务
    log_info "启动 Docker 容器..."
    docker-compose -f $DOCKER_COMPOSE_FILE up -d
    
    if [ $? -ne 0 ]; then
        log_error "Docker Compose 启动失败"
        return 1
    fi

    # 等待服务启动
    log_info "等待服务启动..."
    sleep 5
    
    # 检查容器状态
    check_container_status
    
    if [ $? -ne 0 ]; then
        log_error "容器启动失败，请查看日志"
        return 1
    fi

    # 健康检查
    check_health

    log_success "$MODE 模式部署完成"
}

# 检查容器状态
check_container_status() {
    log_info "检查容器状态..."
    
    local failed_containers=""
    local services="backend frontend postgres redis"
    
    for service in $services; do
        local container_name=$(docker-compose -f $DOCKER_COMPOSE_FILE ps -q $service 2>/dev/null)
        
        if [ -z "$container_name" ]; then
            log_warning "$service 容器未找到"
            continue
        fi
        
        local status=$(docker inspect -f '{{.State.Status}}' $container_name 2>/dev/null)
        
        if [ "$status" = "running" ]; then
            log_success "$service 容器运行中"
        else
            log_error "$service 容器状态异常: $status"
            failed_containers="$failed_containers $service"
            
            # 显示容器日志的最后几行
            log_error "$service 容器日志（最后20行）:"
            docker logs --tail 20 $container_name 2>&1 | sed 's/^/  /'
        fi
    done
    
    if [ -n "$failed_containers" ]; then
        log_error "以下容器启动失败:$failed_containers"
        log_info "查看完整日志: ./deploy.sh $MODE logs [service-name]"
        return 1
    fi
    
    return 0
}

# 健康检查
check_health() {
    log_info "执行健康检查..."

    # 后端健康检查
    log_info "检查后端服务..."
    for i in {1..10}; do
        if curl -f -s "http://localhost:8080/healthz" > /dev/null 2>&1; then
            log_success "后端服务健康"
            break
        else
            if [ $i -eq 10 ]; then
                log_error "后端服务健康检查失败"
                log_info "请查看日志: ./deploy.sh $MODE logs backend"
                return 1
            fi
            log_info "等待后端服务启动... ($i/10)"
            sleep 3
        fi
    done

    # 前端健康检查（仅生产模式）
    if [ "$MODE" = "prod" ]; then
        frontend_port=3000
        log_info "检查前端服务..."
        for i in {1..10}; do
            if curl -f -s "http://localhost:$frontend_port" > /dev/null 2>&1; then
                log_success "前端服务健康"
                break
            else
                if [ $i -eq 10 ]; then
                    log_warning "前端服务可能未就绪，请手动检查"
                    break
                fi
                sleep 3
            fi
        done
    else
        log_info "开发模式：前端需要本地运行"
    fi

    log_success "健康检查完成"
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
    echo "🎮 掼蛋在线游戏 - $MODE 模式"
    echo "================================"
    
    if [ "$MODE" = "dev" ]; then
        echo "🔧 后端服务 (Docker):"
        echo "  - API: http://localhost:8080"
        echo "  - WebSocket: ws://localhost:8080/ws"
        echo "  - 健康检查: http://localhost:8080/healthz"
        echo ""
        echo "💻 前端服务 (本地运行):"
        echo "  请在新终端中运行:"
        echo "  $ cd frontend"
        echo "  $ npm install  (首次运行)"
        echo "  $ npm run dev"
        echo "  然后访问: http://localhost:5173"
        echo ""
        echo "✨ 开发模式特性:"
        echo "  - 后端在Docker中，支持热重载"
        echo "  - 前端本地运行，无Service Worker干扰"
        echo "  - 详细的调试日志"
        echo "  - 直接访问服务端口"
    else
        echo "前端地址: http://localhost:3000"
        echo "后端API: http://localhost:8080"
        echo "WebSocket: ws://localhost:8080/ws"
        echo ""
        echo "📊 监控面板:"
        echo "  - Grafana: http://localhost:3001"
        echo "  - Prometheus: http://localhost:9090"
    fi
    
    echo ""
    echo "📋 常用命令:"
    echo "  查看日志: $0 $MODE logs [service]"
    echo "  查看状态: $0 $MODE status"
    echo "  重启服务: $0 $MODE restart"
    echo "  停止服务: $0 $MODE stop"
    
    if [ "$MODE" = "dev" ]; then
        echo "  重新构建: $0 $MODE rebuild"
    fi
    
    echo ""
    if [ "$MODE" = "dev" ]; then
        echo "💡 提示: 后端在Docker，前端本地运行"
        echo "   前端启动后Vite会自动代理 /api 和 /ws 到后端"
    else
        echo "💡 提示: 所有服务都在 Docker 中运行"
    fi
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
        "rebuild")
            rebuild_images
            ;;
        "logs")
            view_logs "$@"
            ;;
        "status")
            check_status
            ;;
        "health")
            check_health
            ;;
        "clean")
            clean_resources
            ;;
        "backup")
            if [ "$MODE" != "prod" ]; then
                log_warning "备份功能仅在生产模式下可用"
                exit 1
            fi
            backup_data
            ;;
        "restore")
            if [ "$MODE" != "prod" ]; then
                log_warning "恢复功能仅在生产模式下可用"
                exit 1
            fi
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