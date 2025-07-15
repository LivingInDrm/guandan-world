package sdk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestData 表示 JSON 测试数据的结构
type TestData struct {
	Level    int        `json:"level"`
	CompList []CompCase `json:"comp_list"`
}

// CompCase 表示单个测试用例
type CompCase struct {
	Cards [][]interface{} `json:"cards"`
	Type  string          `json:"type"`
}

// 将 JSON 中的牌型名称映射到 Go 的 CompType
func stringToCompType(typeStr string) CompType {
	switch typeStr {
	case "Fold":
		return TypeFold
	case "IllegalComp":
		return TypeIllegal
	case "Single":
		return TypeSingle
	case "Pair":
		return TypePair
	case "Triple":
		return TypeTriple
	case "FullHouse":
		return TypeFullHouse
	case "Straight":
		return TypeStraight
	case "Plate":
		return TypePlate
	case "Tube":
		return TypeTube
	case "JokerBomb":
		return TypeJokerBomb
	case "NaiveBomb":
		return TypeNaiveBomb
	case "StraightFlush":
		return TypeStraightFlush
	default:
		return TypeIllegal
	}
}

// 将 JSON 中的牌数据转换为 Card 结构
func jsonToCard(cardData []interface{}, level int) (*Card, error) {
	if len(cardData) != 2 {
		return nil, nil
	}

	var number int
	var color string

	// 处理 number 字段（可能是 int 或 float64）
	switch v := cardData[0].(type) {
	case int:
		number = v
	case float64:
		number = int(v)
	default:
		return nil, nil
	}

	// 处理 color 字段
	color = cardData[1].(string)

	return NewCard(number, color, level)
}

// TestFromCardListBatch 使用 comp.json 数据进行批量测试
func TestFromCardListBatch(t *testing.T) {
	// 读取测试数据文件
	testDataPath := filepath.Join("..", "test-data", "comp.json")
	data, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("无法读取测试数据文件: %v", err)
	}

	// 解析 JSON 数据
	var testData TestData
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("无法解析测试数据: %v", err)
	}

	t.Logf("正在使用级别 %d 进行测试", testData.Level)
	t.Logf("总共有 %d 个测试用例", len(testData.CompList))

	// 统计测试结果
	passCount := 0
	failCount := 0
	totalCount := len(testData.CompList)

	// 遍历所有测试用例
	for i, testCase := range testData.CompList {
		t.Run(formatTestName(i, testCase), func(t *testing.T) {
			// 将 JSON 卡片数据转换为 Card 结构
			var cards []*Card
			for _, cardData := range testCase.Cards {
				card, err := jsonToCard(cardData, testData.Level)
				if err != nil {
					t.Errorf("无法创建牌: %v", err)
					failCount++
					return
				}
				if card != nil {
					cards = append(cards, card)
				}
			}

			// 调用 FromCardList 函数
			result := FromCardList(cards, nil)

			// 获取期望的牌型
			expectedType := stringToCompType(testCase.Type)
			actualType := result.GetType()

			// 验证结果
			if actualType != expectedType {
				t.Errorf("测试用例 %d 失败:\n  输入牌: %s\n  期望类型: %s\n  实际类型: %s\n  结果: %s",
					i+1, formatCards(cards), testCase.Type, actualType.String(), result.String())
				failCount++
			} else {
				t.Logf("测试用例 %d 通过: %s -> %s", i+1, formatCards(cards), actualType.String())
				passCount++
			}
		})
	}

	// 输出统计信息
	t.Logf("\n=== 批量测试结果统计 ===")
	t.Logf("总测试用例数: %d", totalCount)
	t.Logf("通过: %d (%.1f%%)", passCount, float64(passCount)/float64(totalCount)*100)
	t.Logf("失败: %d (%.1f%%)", failCount, float64(failCount)/float64(totalCount)*100)

	if failCount > 0 {
		t.Errorf("有 %d 个测试用例失败", failCount)
	}
}

// 格式化测试名称
func formatTestName(index int, testCase CompCase) string {
	return fmt.Sprintf("Case%d_%s_%s", index+1, testCase.Type, formatCards(jsonToCards(testCase.Cards, 5)))
}

// 格式化牌组显示
func formatCards(cards []*Card) string {
	if len(cards) == 0 {
		return "Empty"
	}

	var cardStrs []string
	for _, card := range cards {
		if card.Color == "Joker" {
			cardStrs = append(cardStrs, fmt.Sprintf("%s", card.Name))
		} else {
			cardStrs = append(cardStrs, fmt.Sprintf("%d%s", card.RawNumber, card.Color[:1]))
		}
	}

	return fmt.Sprintf("[%s]", strings.Join(cardStrs, ","))
}

// 将 JSON 卡片数据转换为 Card 数组（用于格式化显示）
func jsonToCards(cardDataList [][]interface{}, level int) []*Card {
	var cards []*Card
	for _, cardData := range cardDataList {
		card, err := jsonToCard(cardData, level)
		if err == nil && card != nil {
			cards = append(cards, card)
		}
	}
	return cards
}

// TestFromCardListBatchSpecific 测试特定的测试用例
func TestFromCardListBatchSpecific(t *testing.T) {
	// 读取测试数据文件
	testDataPath := filepath.Join("..", "test-data", "comp.json")
	data, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("无法读取测试数据文件: %v", err)
	}

	// 解析 JSON 数据
	var testData TestData
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("无法解析测试数据: %v", err)
	}

	// 测试特定的几个用例
	testCases := []int{0, 1, 2, 3, 4, 5} // 测试前几个用例

	for _, caseIndex := range testCases {
		if caseIndex >= len(testData.CompList) {
			continue
		}

		testCase := testData.CompList[caseIndex]
		t.Run(fmt.Sprintf("SpecificCase_%d", caseIndex+1), func(t *testing.T) {
			// 转换卡片数据
			var cards []*Card
			for _, cardData := range testCase.Cards {
				card, err := jsonToCard(cardData, testData.Level)
				if err != nil {
					t.Errorf("无法创建牌: %v", err)
					return
				}
				if card != nil {
					cards = append(cards, card)
				}
			}

			// 调用 FromCardList 函数
			result := FromCardList(cards, nil)

			// 获取期望的牌型
			expectedType := stringToCompType(testCase.Type)
			actualType := result.GetType()

			// 详细输出
			t.Logf("测试用例 %d:", caseIndex+1)
			t.Logf("  输入牌: %s", formatCards(cards))
			t.Logf("  期望类型: %s", testCase.Type)
			t.Logf("  实际类型: %s", actualType.String())
			t.Logf("  结果: %s", result.String())
			t.Logf("  是否有效: %v", result.IsValid())

			// 验证结果
			if actualType != expectedType {
				t.Errorf("测试用例 %d 失败: 期望 %s, 实际 %s", caseIndex+1, testCase.Type, actualType.String())
			}
		})
	}
}

// TestFromCardListBatchAnalysis 分析失败的测试用例
func TestFromCardListBatchAnalysis(t *testing.T) {
	// 读取测试数据文件
	testDataPath := filepath.Join("..", "test-data", "comp.json")
	data, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("无法读取测试数据文件: %v", err)
	}

	// 解析 JSON 数据
	var testData TestData
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("无法解析测试数据: %v", err)
	}

	// 重点分析失败的测试用例
	failedCases := []int{
		39, 43, 45, 47, 48, 50, 52, 57, 71, 72, 78, 79, 97, 98, 99, 100, // 基于之前的测试结果
	}

	fmt.Printf("\n=== 详细分析失败的测试用例 ===\n")

	for _, caseIndex := range failedCases {
		if caseIndex >= len(testData.CompList) {
			continue
		}

		testCase := testData.CompList[caseIndex]
		fmt.Printf("\n--- 测试用例 %d ---\n", caseIndex+1)

		// 转换卡片数据
		var cards []*Card
		for _, cardData := range testCase.Cards {
			card, err := jsonToCard(cardData, testData.Level)
			if err != nil {
				fmt.Printf("无法创建牌: %v\n", err)
				continue
			}
			if card != nil {
				cards = append(cards, card)
			}
		}

		// 调用 FromCardList 函数
		result := FromCardList(cards, nil)

		// 获取期望的牌型
		expectedType := stringToCompType(testCase.Type)
		actualType := result.GetType()

		// 详细分析
		fmt.Printf("期望类型: %s\n", testCase.Type)
		fmt.Printf("实际类型: %s\n", actualType.String())
		fmt.Printf("输入牌: %s\n", formatCardsDetailed(cards))
		fmt.Printf("结果: %s\n", result.String())
		fmt.Printf("是否有效: %v\n", result.IsValid())

		// 根据不同类型进行分析
		if expectedType == TypeStraight || expectedType == TypeStraightFlush {
			fmt.Printf("顺子分析:\n")
			analyzeStraight(cards, testData.Level)
		} else if expectedType == TypePlate {
			fmt.Printf("钢板分析:\n")
			analyzePlate(cards, testData.Level)
		} else if expectedType == TypeTube {
			fmt.Printf("钢管分析:\n")
			analyzeTube(cards, testData.Level)
		}
	}
}

// formatCardsDetailed 详细格式化牌组显示
func formatCardsDetailed(cards []*Card) string {
	if len(cards) == 0 {
		return "Empty"
	}

	var cardStrs []string
	for _, card := range cards {
		if card.Color == "Joker" {
			cardStrs = append(cardStrs, fmt.Sprintf("%s", card.Name))
		} else {
			cardStrs = append(cardStrs, fmt.Sprintf("%s(%d/%d)", card.String(), card.Number, card.RawNumber))
		}
	}

	return fmt.Sprintf("[%s]", strings.Join(cardStrs, ", "))
}

// analyzeStraight 分析顺子
func analyzeStraight(cards []*Card, level int) {
	fmt.Printf("  卡片数量: %d\n", len(cards))
	if len(cards) != 5 {
		fmt.Printf("  错误: 顺子必须是5张牌\n")
		return
	}

	wildcardCount := 0
	normalCards := []*Card{}
	for _, card := range cards {
		if card.IsWildcard() {
			wildcardCount++
			fmt.Printf("  变化牌: %s\n", card.String())
		} else {
			normalCards = append(normalCards, card)
		}
	}

	fmt.Printf("  变化牌数量: %d\n", wildcardCount)
	fmt.Printf("  正常牌数量: %d\n", len(normalCards))

	if len(normalCards) > 0 {
		// 按 RawNumber 排序
		sort.Slice(normalCards, func(i, j int) bool {
			return normalCards[i].RawNumber < normalCards[j].RawNumber
		})

		fmt.Printf("  正常牌序列 (RawNumber): ")
		for i, card := range normalCards {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%d", card.RawNumber)
		}
		fmt.Printf("\n")

		// 检查间隔
		if len(normalCards) > 1 {
			gaps := 0
			for i := 1; i < len(normalCards); i++ {
				gap := normalCards[i].RawNumber - normalCards[i-1].RawNumber - 1
				if gap > 0 {
					gaps += gap
				}
			}
			fmt.Printf("  间隔总数: %d\n", gaps)
		}
	}
}

// analyzePlate 分析钢板
func analyzePlate(cards []*Card, level int) {
	fmt.Printf("  卡片数量: %d\n", len(cards))
	if len(cards) != 6 {
		fmt.Printf("  错误: 钢板必须是6张牌\n")
		return
	}

	cardCounts := make(map[int]int)
	wildcardCount := 0

	for _, card := range cards {
		if card.IsWildcard() {
			wildcardCount++
		} else {
			cardCounts[card.Number]++
		}
	}

	fmt.Printf("  变化牌数量: %d\n", wildcardCount)
	fmt.Printf("  牌数统计: ")
	for num, count := range cardCounts {
		fmt.Printf("%d:%d张 ", num, count)
	}
	fmt.Printf("\n")

	numbers := make([]int, 0)
	for num := range cardCounts {
		numbers = append(numbers, num)
	}
	sort.Ints(numbers)

	if len(numbers) >= 2 {
		fmt.Printf("  是否连续: %v\n", numbers[1] == numbers[0]+1)
	}
}

// analyzeTube 分析钢管
func analyzeTube(cards []*Card, level int) {
	fmt.Printf("  卡片数量: %d\n", len(cards))
	if len(cards) != 6 {
		fmt.Printf("  错误: 钢管必须是6张牌\n")
		return
	}

	cardCounts := make(map[int]int)
	wildcardCount := 0

	for _, card := range cards {
		if card.IsWildcard() {
			wildcardCount++
		} else {
			cardCounts[card.Number]++
		}
	}

	fmt.Printf("  变化牌数量: %d\n", wildcardCount)
	fmt.Printf("  牌数统计: ")
	for num, count := range cardCounts {
		fmt.Printf("%d:%d张 ", num, count)
	}
	fmt.Printf("\n")

	numbers := make([]int, 0)
	for num := range cardCounts {
		numbers = append(numbers, num)
	}
	sort.Ints(numbers)

	if len(numbers) >= 3 {
		fmt.Printf("  是否连续: %v\n", len(numbers) == 3 && numbers[1] == numbers[0]+1 && numbers[2] == numbers[1]+1)
	}
}
