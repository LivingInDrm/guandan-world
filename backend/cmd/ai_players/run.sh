#!/bin/bash

# AI玩家测试客户端启动脚本

# 默认参数
SERVER="localhost:8080"
ROOM_ID=""
VERBOSE=""
NUM_PLAYERS="3"

# 解析命令行参数
while [[ $# -gt 0 ]]; do
  case $1 in
    -r|--room-id)
      ROOM_ID="$2"
      shift 2
      ;;
    -s|--server)
      SERVER="$2"
      shift 2
      ;;
    -v|--verbose)
      VERBOSE="-verbose"
      shift
      ;;
    -n|--num-players)
      NUM_PLAYERS="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  -r, --room-id ID       房间ID (必需)"
      echo "  -s, --server HOST:PORT 服务器地址 (默认: localhost:8080)"
      echo "  -v, --verbose          启用详细日志"
      echo "  -n, --num-players N    AI玩家数量 (默认: 3)"
      echo "  -h, --help             显示此帮助信息"
      echo ""
      echo "示例:"
      echo "  $0 -r room_1234567890"
      echo "  $0 -r room_1234567890 -v"
      echo "  $0 -r room_1234567890 -s 192.168.1.100:8080 -n 2"
      exit 0
      ;;
    *)
      echo "未知参数: $1"
      echo "使用 -h 或 --help 查看帮助"
      exit 1
      ;;
  esac
done

# 检查房间ID
if [ -z "$ROOM_ID" ]; then
  echo "错误: 必须指定房间ID"
  echo "使用 -h 或 --help 查看帮助"
  exit 1
fi

echo "=== 启动AI玩家测试客户端 ==="
echo "服务器: $SERVER"
echo "房间ID: $ROOM_ID"
echo "AI玩家数量: $NUM_PLAYERS"
if [ -n "$VERBOSE" ]; then
  echo "详细日志: 启用"
fi
echo "=============================="
echo ""

# 运行程序
go run main.go -server "$SERVER" -room-id "$ROOM_ID" -num-players "$NUM_PLAYERS" $VERBOSE

