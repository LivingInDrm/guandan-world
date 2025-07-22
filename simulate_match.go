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
	err := simulator.SimulateMatch()
	if err != nil {
		log.Fatalf("模拟失败: %v", err)
	}

	totalDuration := time.Since(startTime)

	// 输出结果
	fmt.Println()
	fmt.Println("🎉 模拟完成！")
	fmt.Printf("⏱️  总耗时: %v\n", totalDuration)
	fmt.Println()

	fmt.Println()
	fmt.Println("✨ 模拟结束")
}
