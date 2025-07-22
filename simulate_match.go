package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"guandan-world/sdk"
)

func main() {
	fmt.Println("🀄 掼蛋牌局模拟器 🀄")
	fmt.Println("==================")
	fmt.Println()

	// 检查命令行参数
	verbose := true
	if len(os.Args) > 1 && os.Args[1] == "-q" {
		verbose = false
	}

	// 创建模拟器
	simulator := sdk.NewMatchSimulator(verbose)

	fmt.Println("开始模拟掼蛋牌局...")
	startTime := time.Now()

	// 运行模拟
	result, err := simulator.SimulateMatch()
	if err != nil {
		log.Fatalf("模拟失败: %v", err)
	}

	totalDuration := time.Since(startTime)

	// 输出结果
	fmt.Println()
	fmt.Println("🎉 模拟完成！")
	fmt.Printf("⏱️  总耗时: %v\n", totalDuration)
	fmt.Println()

	if result != nil {
		fmt.Println("📊 最终结果:")
		fmt.Printf("   🏆 获胜队伍: %d\n", result.Winner)
		fmt.Printf("   📈 最终等级: 队伍0=%d级, 队伍1=%d级\n",
			result.FinalLevels[0], result.FinalLevels[1])
		fmt.Printf("   ⏰ 比赛时长: %v\n", result.Duration)

		if result.Statistics != nil {
			fmt.Printf("   🎯 总局数: %d\n", result.Statistics.TotalDeals)
			fmt.Println()
			fmt.Println("   📋 队伍统计:")
			for i, teamStats := range result.Statistics.TeamStats {
				if teamStats != nil {
					fmt.Printf("      队伍%d: 获胜%d局, 升级%d级, 总墩数%d\n",
						i, teamStats.DealsWon, teamStats.Upgrades, teamStats.TotalTricks)
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("✨ 模拟结束")
}
