package sdk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ComparisonTestData 表示比较测试数据的结构
type ComparisonTestData struct {
	Level                             int                  `json:"level"`
	Description                       string               `json:"description"`
	TotalComparisons                  int                  `json:"total_comparisons"`
	IntraTypeComparisons              int                  `json:"intra_type_comparisons"`
	IntraTypeCrossWildcardComparisons int                  `json:"intra_type_cross_wildcard_comparisons"`
	InterTypeComparisons              int                  `json:"inter_type_comparisons"`
	Comparisons                       []ComparisonTestCase `json:"comparisons"`
}

// ComparisonTestCase 表示单个比较测试用例
type ComparisonTestCase struct {
	TestID                int      `json:"test_id"` // 测试编号，用于快速定位
	ComparisonType        string   `json:"comparison_type"`
	CompType              string   `json:"comp_type"`
	WildcardCount         int      `json:"wildcard_count,omitempty"`
	WildcardCount1        int      `json:"wildcard_count_1,omitempty"`
	WildcardCount2        int      `json:"wildcard_count_2,omitempty"`
	Comp1                 CompData `json:"comp1"`
	Comp2                 CompData `json:"comp2"`
	Comp1GreaterThanComp2 bool     `json:"comp1_greater_than_comp2"`
	Comp2GreaterThanComp1 bool     `json:"comp2_greater_than_comp1"`
}

// CompData 表示单个牌组数据
type CompData struct {
	Cards [][]interface{} `json:"cards"`
	Type  string          `json:"type"`
}

// GetDebugCommand 返回快速定位该测试用例的命令
func (tc *ComparisonTestCase) GetDebugCommand() string {
	return fmt.Sprintf("jq '.comparisons[%d]' comparison_test_data.json", tc.TestID)
}

// TestComparisonBatch 批量测试牌组比较功能
func TestComparisonBatch(t *testing.T) {
	// 读取测试数据文件
	testDataPath := filepath.Join("..", "test-data", "comparison_test_data.json")
	data, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("无法读取测试数据文件: %v", err)
	}

	// 解析 JSON 数据
	var testData ComparisonTestData
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("无法解析测试数据: %v", err)
	}

	t.Logf("开始批量比较测试 - 级别: %d", testData.Level)
	t.Logf("总比较数: %d", testData.TotalComparisons)
	t.Logf("同类型比较: %d", testData.IntraTypeComparisons)
	t.Logf("跨变化牌比较: %d", testData.IntraTypeCrossWildcardComparisons)
	t.Logf("不同类型比较: %d", testData.InterTypeComparisons)

	// 统计测试结果
	passCount := 0
	failCount := 0
	totalCount := len(testData.Comparisons)

	// 按比较类型分组统计
	intraTypeStats := make(map[string]int)
	crossWildcardStats := make(map[string]int)
	interTypeStats := make(map[string]int)

	// 遍历所有测试用例
	for _, testCase := range testData.Comparisons {
		t.Run(fmt.Sprintf("TestID_%d_%s_%s", testCase.TestID, testCase.ComparisonType, testCase.CompType), func(t *testing.T) {
			// 创建第二个牌组（被比较方）
			// 根据测试数据中的type创建指定类型的牌组
			comp2Cards := convertJSONToCards(testCase.Comp2.Cards, testData.Level)
			comp2 := CreateCompByType(comp2Cards, testCase.Comp2.Type)

			// 对comp2进行normalize，将万能牌替换为具体的牌
			normalizedComp2 := NormalizeComp(comp2)

			// 创建第一个牌组（主动比较方）
			// 使用FromCardList，传入normalizedComp2作为prev参数
			// 这样comp1会优先尝试与comp2相同的类型
			comp1Cards := convertJSONToCards(testCase.Comp1.Cards, testData.Level)
			comp1 := FromCardList(comp1Cards, normalizedComp2)

			// 执行比较
			// comp1可能包含万能牌，normalizedComp2已经没有万能牌
			actualComp1Greater := comp1.GreaterThan(normalizedComp2)
			actualComp2Greater := normalizedComp2.GreaterThan(comp1)

			// 验证结果
			success := true
			if actualComp1Greater != testCase.Comp1GreaterThanComp2 {
				t.Errorf("comp1 > comp2 比较失败: 期望 %v, 实际 %v", testCase.Comp1GreaterThanComp2, actualComp1Greater)
				success = false
			}
			if actualComp2Greater != testCase.Comp2GreaterThanComp1 {
				t.Errorf("comp2 > comp1 比较失败: 期望 %v, 实际 %v", testCase.Comp2GreaterThanComp1, actualComp2Greater)
				success = false
			}

			if success {
				passCount++
				t.Logf("✓ [TestID:%d] 比较成功: %s vs %s", testCase.TestID, formatCompForLog(comp1), formatCompForLog(normalizedComp2))
			} else {
				failCount++
				t.Errorf("🚨 [TestID:%d] 比较失败:", testCase.TestID)
				t.Errorf("📍 快速定位: %s", testCase.GetDebugCommand())
				t.Errorf("  Comp1: %s", formatCompForLog(comp1))
				t.Errorf("  Comp2 (original): %s", formatCompForLog(comp2))
				t.Errorf("  Comp2 (normalized): %s", formatCompForLog(normalizedComp2))
				t.Errorf("  期望: comp1>comp2=%v, comp2>comp1=%v", testCase.Comp1GreaterThanComp2, testCase.Comp2GreaterThanComp1)
				t.Errorf("  实际: comp1>comp2=%v, comp2>comp1=%v", actualComp1Greater, actualComp2Greater)
			}

			// 统计各类型结果
			switch testCase.ComparisonType {
			case "intra_type":
				if success {
					intraTypeStats[testCase.CompType]++
				} else {
					intraTypeStats[testCase.CompType+"_failed"]++
				}
			case "intra_type_cross_wildcard":
				if success {
					crossWildcardStats[testCase.CompType]++
				} else {
					crossWildcardStats[testCase.CompType+"_failed"]++
				}
			case "inter_type":
				if success {
					interTypeStats["success"]++
				} else {
					interTypeStats["failed"]++
				}
			}
		})
	}

	// 输出统计信息
	t.Logf("\n=== 批量比较测试结果统计 ===")
	t.Logf("总测试用例数: %d", totalCount)
	t.Logf("通过: %d (%.1f%%)", passCount, float64(passCount)/float64(totalCount)*100)
	t.Logf("失败: %d (%.1f%%)", failCount, float64(failCount)/float64(totalCount)*100)

	t.Logf("\n=== 同类型比较统计 ===")
	for compType, count := range intraTypeStats {
		if !strings.HasSuffix(compType, "_failed") {
			failed := intraTypeStats[compType+"_failed"]
			t.Logf("%s: 通过 %d, 失败 %d", compType, count, failed)
		}
	}

	t.Logf("\n=== 跨变化牌比较统计 ===")
	for compType, count := range crossWildcardStats {
		if !strings.HasSuffix(compType, "_failed") {
			failed := crossWildcardStats[compType+"_failed"]
			t.Logf("%s: 通过 %d, 失败 %d", compType, count, failed)
		}
	}

	t.Logf("\n=== 不同类型比较统计 ===")
	t.Logf("通过: %d, 失败: %d", interTypeStats["success"], interTypeStats["failed"])

	if failCount > 0 {
		t.Errorf("有 %d 个比较测试用例失败", failCount)
	}
}

// convertJSONToCards 将 JSON 卡片数据转换为 Card 数组
func convertJSONToCards(cardDataList [][]interface{}, level int) []*Card {
	var cards []*Card
	for _, cardData := range cardDataList {
		card, err := jsonToCard(cardData, level)
		if err == nil && card != nil {
			cards = append(cards, card)
		}
	}
	return cards
}

// formatCompForLog 格式化牌组用于日志输出
func formatCompForLog(comp CardComp) string {
	if comp == nil {
		return "nil"
	}

	cards := comp.GetCards()
	if len(cards) == 0 {
		return fmt.Sprintf("%s: Empty", comp.GetType().String())
	}

	var cardStrs []string
	for _, card := range cards {
		if card == nil {
			cardStrs = append(cardStrs, "nil")
			continue
		}
		if card.Color == "Joker" {
			cardStrs = append(cardStrs, fmt.Sprintf("%s", card.Name))
		} else {
			cardStrs = append(cardStrs, fmt.Sprintf("%d%s", card.RawNumber, card.Color[:1]))
		}
	}

	return fmt.Sprintf("%s: [%s]", comp.GetType().String(), strings.Join(cardStrs, ","))
}

// TestComparisonSpecific 测试特定的比较用例
func TestComparisonSpecific(t *testing.T) {
	// 读取测试数据文件
	testDataPath := filepath.Join("..", "test-data", "comparison_test_data.json")
	data, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("无法读取测试数据文件: %v", err)
	}

	// 解析 JSON 数据
	var testData ComparisonTestData
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("无法解析测试数据: %v", err)
	}

	// 测试特定的test_id用例
	specificTestIDs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} // 测试前10个test_id

	for _, targetTestID := range specificTestIDs {
		// 根据test_id查找对应的测试用例
		var testCase *ComparisonTestCase
		for i := range testData.Comparisons {
			if testData.Comparisons[i].TestID == targetTestID {
				testCase = &testData.Comparisons[i]
				break
			}
		}

		if testCase == nil {
			t.Logf("TestID %d 未找到，跳过", targetTestID)
			continue
		}

		t.Run(fmt.Sprintf("TestID_%d_Specific", testCase.TestID), func(t *testing.T) {
			// 创建牌组
			comp1Cards := convertJSONToCards(testCase.Comp1.Cards, testData.Level)
			comp1 := FromCardList(comp1Cards, nil)

			comp2Cards := convertJSONToCards(testCase.Comp2.Cards, testData.Level)
			comp2 := FromCardList(comp2Cards, nil)

			// 执行比较
			actualComp1Greater := comp1.GreaterThan(comp2)
			actualComp2Greater := comp2.GreaterThan(comp1)

			// 详细输出
			t.Logf("🔍 [TestID:%d] 详细测试 (%s):", testCase.TestID, testCase.ComparisonType)
			t.Logf("📍 快速定位: %s", testCase.GetDebugCommand())
			t.Logf("  Comp1: %s", formatCompForLog(comp1))
			t.Logf("  Comp2: %s", formatCompForLog(comp2))
			t.Logf("  期望: comp1>comp2=%v, comp2>comp1=%v", testCase.Comp1GreaterThanComp2, testCase.Comp2GreaterThanComp1)
			t.Logf("  实际: comp1>comp2=%v, comp2>comp1=%v", actualComp1Greater, actualComp2Greater)

			// 验证结果
			if actualComp1Greater != testCase.Comp1GreaterThanComp2 {
				t.Errorf("🚨 [TestID:%d] comp1 > comp2 比较失败: 期望 %v, 实际 %v", testCase.TestID, testCase.Comp1GreaterThanComp2, actualComp1Greater)
				t.Errorf("📍 快速定位: %s", testCase.GetDebugCommand())
			}
			if actualComp2Greater != testCase.Comp2GreaterThanComp1 {
				t.Errorf("🚨 [TestID:%d] comp2 > comp1 比较失败: 期望 %v, 实际 %v", testCase.TestID, testCase.Comp2GreaterThanComp1, actualComp2Greater)
				t.Errorf("📍 快速定位: %s", testCase.GetDebugCommand())
			}
		})
	}
}

// TestAnalyzeFailedCasesWildcardDistribution 分析失败用例的wildcard分布
func TestAnalyzeFailedCasesWildcardDistribution(t *testing.T) {
	// 读取测试数据文件
	testDataPath := filepath.Join("..", "test-data", "comparison_test_data.json")
	data, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("无法读取测试数据文件: %v", err)
	}

	// 解析 JSON 数据
	var testData ComparisonTestData
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("无法解析测试数据: %v", err)
	}

	// 统计失败用例的wildcard分布
	wildcardDistribution := make(map[string]map[int]int) // compType -> wildcardCount -> failureCount
	totalFailures := make(map[string]int)                // compType -> totalFailures

	t.Logf("=== 开始分析失败用例的Wildcard分布 ===")

	failureCount := 0
	for _, testCase := range testData.Comparisons {
		// 创建牌组
		comp1Cards := convertJSONToCards(testCase.Comp1.Cards, testData.Level)
		comp1 := FromCardList(comp1Cards, nil)

		comp2Cards := convertJSONToCards(testCase.Comp2.Cards, testData.Level)
		comp2 := FromCardList(comp2Cards, nil)

		// 执行比较
		actualComp1Greater := comp1.GreaterThan(comp2)
		actualComp2Greater := comp2.GreaterThan(comp1)

		// 检查是否失败
		isFailed := (actualComp1Greater != testCase.Comp1GreaterThanComp2) ||
			(actualComp2Greater != testCase.Comp2GreaterThanComp1)

		if isFailed {
			failureCount++

			// 分析两个牌组的wildcard数量
			comp1Type := comp1.GetType().String()
			comp2Type := comp2.GetType().String()

			comp1WildcardCount := countWildcards(comp1Cards)
			comp2WildcardCount := countWildcards(comp2Cards)

			// 初始化map
			if wildcardDistribution[comp1Type] == nil {
				wildcardDistribution[comp1Type] = make(map[int]int)
			}
			if wildcardDistribution[comp2Type] == nil {
				wildcardDistribution[comp2Type] = make(map[int]int)
			}

			// 记录失败情况
			wildcardDistribution[comp1Type][comp1WildcardCount]++
			wildcardDistribution[comp2Type][comp2WildcardCount]++

			totalFailures[comp1Type]++
			totalFailures[comp2Type]++

			if failureCount <= 10 { // 只显示前10个失败用例的详细信息
				t.Logf("失败用例 [TestID:%d] %s:", testCase.TestID, testCase.ComparisonType)
				t.Logf("  Comp1: %s (wildcards: %d) - %s", comp1Type, comp1WildcardCount, formatCompForLog(comp1))
				t.Logf("  Comp2: %s (wildcards: %d) - %s", comp2Type, comp2WildcardCount, formatCompForLog(comp2))
			}
		}
	}

	t.Logf("总共发现 %d 个失败用例", failureCount)

	// 输出统计结果
	t.Logf("=== Wildcard分布统计 ===")

	// 按牌型排序输出
	for compType, distribution := range wildcardDistribution {
		t.Logf("%s类型 - 总失败次数: %d", compType, totalFailures[compType])

		for wildcardCount := 0; wildcardCount <= 5; wildcardCount++ {
			if failures, exists := distribution[wildcardCount]; exists {
				percentage := float64(failures) / float64(totalFailures[compType]) * 100
				t.Logf("  %d个wildcard: %d次失败 (%.1f%%)", wildcardCount, failures, percentage)
			}
		}
	}

	// 总体wildcard分布
	t.Logf("=== 总体Wildcard分布 ===")
	totalWildcardCount := make(map[int]int)
	totalAllFailures := 0

	for _, distribution := range wildcardDistribution {
		for wildcardCount, failures := range distribution {
			totalWildcardCount[wildcardCount] += failures
			totalAllFailures += failures
		}
	}

	for wildcardCount := 0; wildcardCount <= 5; wildcardCount++ {
		if failures, exists := totalWildcardCount[wildcardCount]; exists {
			percentage := float64(failures) / float64(totalAllFailures) * 100
			t.Logf("%d个wildcard: %d次失败 (%.1f%%)", wildcardCount, failures, percentage)
		}
	}
}
