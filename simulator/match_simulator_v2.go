package simulator

import (
	"fmt"
	"time"

	"guandan-world/ai"
	"guandan-world/sdk"
)

// MatchSimulatorV2 重构后的比赛模拟器
// 新架构：专注于输入提供和事件观察，游戏循环由SDK内部处理
type MatchSimulatorV2 struct {
	driver        *sdk.GameDriver          // 游戏驱动器
	inputProvider *SimulatingInputProvider // 模拟输入提供者
	observer      *MatchSimulatorObserver  // 事件观察者
	verbose       bool                     // 是否详细输出
	eventLog      []string                 // 事件日志
}

// NewMatchSimulatorV2 创建新的比赛模拟器V2
func NewMatchSimulatorV2(verbose bool) *MatchSimulatorV2 {
	// 创建游戏引擎
	engine := sdk.NewGameEngine()

	// 创建游戏驱动器
	driver := sdk.NewGameDriver(engine, sdk.DefaultGameDriverConfig())

	// 创建模拟输入提供者
	inputProvider := NewSimulatingInputProvider()

	// 创建事件观察者
	observer := NewMatchSimulatorObserver(engine, verbose, nil)

	simulator := &MatchSimulatorV2{
		driver:        driver,
		inputProvider: inputProvider,
		observer:      observer,
		verbose:       verbose,
		eventLog:      make([]string, 0),
	}

	// 设置自定义日志函数
	observer.SetLogger(simulator.log)

	// 设置输入提供者
	driver.SetInputProvider(inputProvider)

	// 添加观察者
	driver.AddObserver(observer)

	return simulator
}

// SimulateMatch 模拟完整比赛
// 新架构下，这个方法变得非常简洁
func (ms *MatchSimulatorV2) SimulateMatch() error {
	// 创建4个模拟玩家
	players := ms.createPlayers()

	// 设置玩家算法
	if err := ms.setupPlayerAlgorithms(); err != nil {
		return fmt.Errorf("failed to setup player algorithms: %w", err)
	}

	ms.log("Match started with 4 players")

	// 在比赛开始前输出初始状态
	ms.observer.logTeamStatus()

	// 运行比赛（所有游戏逻辑都在SDK内部）
	result, err := ms.driver.RunMatch(players)
	if err != nil {
		return fmt.Errorf("failed to run match: %w", err)
	}

	// 打印最终结果
	ms.printMatchSummary(result)

	return nil
}

// createPlayers 创建4个玩家
func (ms *MatchSimulatorV2) createPlayers() []sdk.Player {
	players := make([]sdk.Player, 4)

	for i := 0; i < 4; i++ {
		players[i] = sdk.Player{
			ID:       fmt.Sprintf("player_%d", i),
			Username: fmt.Sprintf("Player %d", i+1),
			Seat:     i,
			Online:   true,
			AutoPlay: true,
		}
	}

	return players
}

// setupPlayerAlgorithms 设置玩家算法
func (ms *MatchSimulatorV2) setupPlayerAlgorithms() error {
	algorithms := make([]ai.AutoPlayAlgorithm, 4)

	for i := 0; i < 4; i++ {
		algorithms[i] = ai.NewSimpleAutoPlayAlgorithm(2) // 从2级开始
	}

	return ms.inputProvider.BatchSetAlgorithms(algorithms)
}

// printMatchSummary 打印比赛总结
func (ms *MatchSimulatorV2) printMatchSummary(result *sdk.GameDriverResult) {
	fmt.Println("\n========== Match Summary ==========")
	fmt.Printf("Total Deals: %d\n", result.DealCount)
	fmt.Printf("Duration: %v\n", result.Duration)

	if result.MatchResult != nil {
		fmt.Printf("Winner: Team %d\n", result.Winner)
		fmt.Printf("Final Levels: Team 0: Level %d, Team 1: Level %d\n",
			result.FinalLevels[0], result.FinalLevels[1])

		if result.Statistics != nil {
			fmt.Printf("Total Tricks: %d\n", result.Statistics.TotalDeals*25) // 估算
		}
	}

	// 打印玩家统计
	if result.PlayerStats != nil {
		fmt.Println("Player Statistics:")
		for i := 0; i < 4; i++ {
			if stats, exists := result.PlayerStats[i]; exists {
				fmt.Printf("  Player %d: %d cards played, %d tricks won\n",
					i, stats.CardsPlayed, stats.TricksWon)
			}
		}
	}

	fmt.Println("===================================")
}

// log 记录日志
func (ms *MatchSimulatorV2) log(message string) {
	ms.eventLog = append(ms.eventLog, message)
	if ms.verbose {
		fmt.Println(message)
	}
}

// GetEventLog 获取事件日志
func (ms *MatchSimulatorV2) GetEventLog() []string {
	return ms.eventLog
}

// SetVerbose 设置详细输出模式
func (ms *MatchSimulatorV2) SetVerbose(verbose bool) {
	ms.verbose = verbose
	ms.observer.SetVerbose(verbose)
}

// GetDriver 获取游戏驱动器（用于高级用法）
func (ms *MatchSimulatorV2) GetDriver() *sdk.GameDriver {
	return ms.driver
}

// GetEngine 获取游戏引擎（用于高级用法）
func (ms *MatchSimulatorV2) GetEngine() sdk.GameEngineInterface {
	return ms.driver.GetEngine()
}

// SetPlayerAlgorithm 设置特定玩家的算法
func (ms *MatchSimulatorV2) SetPlayerAlgorithm(playerSeat int, algorithm ai.AutoPlayAlgorithm) {
	ms.inputProvider.SetPlayerAlgorithm(playerSeat, algorithm)
}

// AddObserver 添加额外的事件观察者
func (ms *MatchSimulatorV2) AddObserver(observer sdk.EventObserver) {
	ms.driver.AddObserver(observer)
}

// RemoveObserver 移除事件观察者
func (ms *MatchSimulatorV2) RemoveObserver(observer sdk.EventObserver) {
	ms.driver.RemoveObserver(observer)
}

// 便捷函数

// RunMatchSimulationV2 运行比赛模拟的便捷函数（新版本）
func RunMatchSimulationV2(verbose bool) error {
	simulator := NewMatchSimulatorV2(verbose)
	return simulator.SimulateMatch()
}

// RunVerboseDemoV2 运行详细模式演示（新版本）
func RunVerboseDemoV2() error {
	fmt.Println("🎮 掼蛋比赛模拟器 V2 - 详细模式演示")
	fmt.Println("=====================================")
	fmt.Println("🚀 开始模拟比赛（详细模式）...")
	fmt.Println("✨ 使用新的GameDriver架构")

	simulator := NewMatchSimulatorV2(true) // 启用详细模式
	err := simulator.SimulateMatch()

	if err != nil {
		fmt.Printf("❌ 模拟失败: %v\n", err)
		return err
	}

	fmt.Println("\n✅ 模拟完成!")
	fmt.Printf("📊 生成了 %d 条事件日志\n", len(simulator.GetEventLog()))

	return nil
}

// CompareSimulators 已删除旧架构，仅测试新模拟器性能
func CompareSimulators() error {
	fmt.Println("🔄 测试新架构模拟器性能...")

	// 测试新版本
	fmt.Println("\n📊 测试新版本模拟器...")
	startTime := time.Now()
	err := RunMatchSimulationV2(false)
	duration := time.Since(startTime)
	if err != nil {
		fmt.Printf("❌ 新版本失败: %v\n", err)
		return err
	} else {
		fmt.Printf("✅ 新版本完成，耗时: %v\n", duration)
	}

	fmt.Println("\n✨ 新架构模拟器运行正常!")
	return nil
}
