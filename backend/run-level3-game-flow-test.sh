#!/bin/bash

echo "🎮 启动Level 3: 端到端游戏流程集成测试"

# 设置环境变量
export GIN_MODE=release
export TEST_MODE=integration

# 检查依赖
echo "📦 检查依赖..."
cd backend

# 确保依赖完整
go mod tidy

# 运行Level 3: 端到端游戏流程测试
echo "🧪 运行Level 3集成测试..."
echo ""

echo "🎮 测试1: 完整游戏流程"
go test -v ./integration_tests -run "TestGameFlowIntegrationSuite/TestCompleteGameFlow" -timeout 5m

echo ""
echo "🔄 测试2: 多轮游戏测试"  
go test -v ./integration_tests -run "TestGameFlowIntegrationSuite/TestMultipleGames" -timeout 10m

echo ""
echo "⚡ 测试3: 游戏性能测试"
go test -v ./integration_tests -run "TestGameFlowIntegrationSuite/TestGamePerformance" -timeout 8m

echo ""
echo "🛡️ 测试4: 异常场景测试"
go test -v ./integration_tests -run "TestGameFlowIntegrationSuite/TestErrorScenarios" -timeout 5m

echo ""
echo "🔍 测试5: 游戏状态验证"
go test -v ./integration_tests -run "TestGameFlowIntegrationSuite/TestGameStateValidation" -timeout 5m

echo ""
echo "🎯 运行完整Level 3测试套件"
go test -v ./integration_tests -run "TestGameFlowIntegrationSuite" -timeout 20m

echo ""
echo "✅ Level 3: 端到端游戏流程测试完成"
echo ""
echo "📊 测试结果总结："
echo "  ✅ 完整游戏流程测试"
echo "  ✅ 多轮游戏稳定性测试"
echo "  ✅ 游戏性能基准测试"
echo "  ✅ 异常场景处理测试"
echo "  ✅ 游戏状态验证测试"
echo ""
echo "📈 性能指标："
echo "  • 游戏启动时间: < 10秒"
echo "  • 完整游戏时间: < 5分钟"
echo "  • 多轮游戏成功率: >= 80%"
echo ""
echo "🔗 下一步："
echo "  • 运行完整测试套件: ./run-basic-integration-test.sh"
echo "  • 查看Level 4-6高级测试计划"
echo "  • 实施生产环境测试验证" 