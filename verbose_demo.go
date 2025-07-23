package main

import (
	"fmt"
	"./sdk"
)

func main() {
	fmt.Println("🎮 掼蛋比赛模拟器 - 详细模式演示")
	fmt.Println("=====================================")
	
	// 运行详细模式演示
	err := sdk.RunVerboseDemoV2()
	if err != nil {
		fmt.Printf("❌ 演示失败: %v
", err)
		return
	}
	
	fmt.Println("✅ 演示完成!")
}
