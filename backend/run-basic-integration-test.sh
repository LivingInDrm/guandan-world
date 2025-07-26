#!/bin/bash

echo "🚀 启动后端基础集成测试"

# 设置环境变量
export GIN_MODE=test
export TEST_MODE=integration

# 检查依赖
echo "📦 检查依赖..."
cd backend

# 安装testify（如果需要）
go get github.com/stretchr/testify/suite
go get github.com/stretchr/testify/assert

# 运行基础集成测试
echo "🧪 运行基础集成测试..."
echo ""

echo "📊 Level 1: API接口集成测试"
go test -v ./integration_tests -run "TestBasicAPIFlow" -timeout 5m

echo ""
echo "🎮 Level 2: 完整游戏流程测试" 
go test -v ./integration_tests -run "TestCompleteGameFlow" -timeout 10m

echo ""
echo "⚡ Level 3: 基础性能测试"
go test -v ./integration_tests -run "TestBasicPerformance" -timeout 3m

echo ""
echo "🎯 运行完整基础集成测试套件"
go test -v ./integration_tests -run "TestBasicIntegration" -timeout 15m

echo ""
echo "✅ 基础集成测试完成"
echo ""
echo "📈 查看更多测试选项："
echo "  • 运行单个测试: go test -v ./integration_tests -run TestBasicAPIFlow"
echo "  • 查看测试覆盖率: go test -cover ./integration_tests"
echo "  • 详细输出: go test -v ./integration_tests -args -test.v"
echo ""
echo "🔗 下一步："
echo "  • 查看完整集成测试指南: cat integration-test-guide.md"
echo "  • 实施更多测试级别: Level 4-6 性能和稳定性测试" 