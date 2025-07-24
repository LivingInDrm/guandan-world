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
			// 创建第一个牌组
			comp1Cards := convertJSONToCards(testCase.Comp1.Cards, testData.Level)
			comp1 := FromCardList(comp1Cards, nil)

			// 创建第二个牌组
			comp2Cards := convertJSONToCards(testCase.Comp2.Cards, testData.Level)
			comp2 := FromCardList(comp2Cards, nil)

			// 执行比较
			actualComp1Greater := comp1.GreaterThan(comp2)
			actualComp2Greater := comp2.GreaterThan(comp1)

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
				t.Logf("✓ [TestID:%d] 比较成功: %s vs %s", testCase.TestID, formatCompForLog(comp1), formatCompForLog(comp2))
			} else {
				failCount++
				t.Errorf("🚨 [TestID:%d] 比较失败:", testCase.TestID)
				t.Errorf("📍 快速定位: %s", testCase.GetDebugCommand())
				t.Errorf("  Comp1: %s", formatCompForLog(comp1))
				t.Errorf("  Comp2: %s", formatCompForLog(comp2))
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

// TestComparisonByType 按牌型分类测试比较功能
func TestComparisonByType(t *testing.T) {
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

	// 按牌型分组测试用例
	typeGroups := make(map[string][]ComparisonTestCase)
	for _, testCase := range testData.Comparisons {
		if testCase.ComparisonType == "intra_type" || testCase.ComparisonType == "intra_type_cross_wildcard" {
			typeGroups[testCase.CompType] = append(typeGroups[testCase.CompType], testCase)
		}
	}

	// 测试各个牌型
	for compType, cases := range typeGroups {
		t.Run(fmt.Sprintf("Type_%s", compType), func(t *testing.T) {
			passCount := 0
			failCount := 0

			for _, testCase := range cases {
				// 创建牌组
				comp1Cards := convertJSONToCards(testCase.Comp1.Cards, testData.Level)
				comp1 := FromCardList(comp1Cards, nil)

				comp2Cards := convertJSONToCards(testCase.Comp2.Cards, testData.Level)
				comp2 := FromCardList(comp2Cards, nil)

				// 执行比较
				actualComp1Greater := comp1.GreaterThan(comp2)
				actualComp2Greater := comp2.GreaterThan(comp1)

				// 验证结果
				if actualComp1Greater == testCase.Comp1GreaterThanComp2 &&
					actualComp2Greater == testCase.Comp2GreaterThanComp1 {
					passCount++
				} else {
					failCount++
					t.Errorf("🚨 [TestID:%d] %s类型测试失败:", testCase.TestID, compType)
					t.Errorf("📍 快速定位: %s", testCase.GetDebugCommand())
					t.Errorf("  Comp1: %s", formatCompForLog(comp1))
					t.Errorf("  Comp2: %s", formatCompForLog(comp2))
					t.Errorf("  期望: comp1>comp2=%v, comp2>comp1=%v", testCase.Comp1GreaterThanComp2, testCase.Comp2GreaterThanComp1)
					t.Errorf("  实际: comp1>comp2=%v, comp2>comp1=%v", actualComp1Greater, actualComp2Greater)
				}
			}

			t.Logf("%s: 通过 %d/%d (%.1f%%)", compType, passCount, len(cases), float64(passCount)/float64(len(cases))*100)
			if failCount > 0 {
				t.Errorf("%s 有 %d 个测试用例失败", compType, failCount)
			}
		})
	}
}

// TestInterTypeComparison 测试不同类型之间的比较
func TestInterTypeComparison(t *testing.T) {
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

	passCount := 0
	failCount := 0

	// 只测试不同类型之间的比较
	for _, testCase := range testData.Comparisons {
		if testCase.ComparisonType != "inter_type" {
			continue
		}

		t.Run(fmt.Sprintf("TestID_%d_InterType", testCase.TestID), func(t *testing.T) {
			// 创建牌组
			comp1Cards := convertJSONToCards(testCase.Comp1.Cards, testData.Level)
			comp1 := FromCardList(comp1Cards, nil)

			comp2Cards := convertJSONToCards(testCase.Comp2.Cards, testData.Level)
			comp2 := FromCardList(comp2Cards, nil)

			// 执行比较
			actualComp1Greater := comp1.GreaterThan(comp2)
			actualComp2Greater := comp2.GreaterThan(comp1)

			// 验证结果
			if actualComp1Greater == testCase.Comp1GreaterThanComp2 &&
				actualComp2Greater == testCase.Comp2GreaterThanComp1 {
				passCount++
				t.Logf("✓ [TestID:%d] %s vs %s", testCase.TestID, formatCompForLog(comp1), formatCompForLog(comp2))
			} else {
				failCount++
				t.Errorf("🚨 [TestID:%d] 不同类型比较失败:", testCase.TestID)
				t.Errorf("📍 快速定位: %s", testCase.GetDebugCommand())
				t.Errorf("  Comp1: %s", formatCompForLog(comp1))
				t.Errorf("  Comp2: %s", formatCompForLog(comp2))
				t.Errorf("  期望: comp1>comp2=%v, comp2>comp1=%v", testCase.Comp1GreaterThanComp2, testCase.Comp2GreaterThanComp1)
				t.Errorf("  实际: comp1>comp2=%v, comp2>comp1=%v", actualComp1Greater, actualComp2Greater)
			}
		})
	}

	t.Logf("不同类型比较: 通过 %d, 失败 %d", passCount, failCount)
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
	cards := comp.GetCards()
	if len(cards) == 0 {
		return fmt.Sprintf("%s: Empty", comp.GetType().String())
	}

	var cardStrs []string
	for _, card := range cards {
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
