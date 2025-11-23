package main

import (
	"encoding/json"
	"fmt"
	"guandan-world/sdk"
	"os"
	"path/filepath"
	"strings"
)

type CompData struct {
	Cards [][]interface{} `json:"cards"`
	Type  string          `json:"type"`
}

type ComparisonTestCase struct {
	TestID                int      `json:"test_id"`
	ComparisonType        string   `json:"comparison_type"`
	CompType              string   `json:"comp_type"`
	Comp1                 CompData `json:"comp1"`
	Comp2                 CompData `json:"comp2"`
	Comp1GreaterThanComp2 bool     `json:"comp1_greater_than_comp2"`
	Comp2GreaterThanComp1 bool     `json:"comp2_greater_than_comp1"`
}

type TestData struct {
	Level       int                   `json:"level"`
	Comparisons []ComparisonTestCase `json:"comparisons"`
}

func convertJSONToCards(jsonCards [][]interface{}, level int) []*sdk.Card {
	cards := make([]*sdk.Card, 0)
	for _, cardData := range jsonCards {
		number := int(cardData[0].(float64))
		color := cardData[1].(string)
		card, _ := sdk.NewCard(number, color, level)
		cards = append(cards, card)
	}
	return cards
}

func printCard(card *sdk.Card) string {
	wc := ""
	if card.IsWildcard() {
		wc = "[WC]"
	}
	return fmt.Sprintf("%s%s", card.Name, wc)
}

func printCards(cards []*sdk.Card) string {
	strs := make([]string, len(cards))
	for i, card := range cards {
		strs[i] = printCard(card)
	}
	return strings.Join(strs, ", ")
}

func debugTestCase(testID int, testData TestData) {
	var testCase ComparisonTestCase
	found := false
	for _, tc := range testData.Comparisons {
		if tc.TestID == testID {
			testCase = tc
			found = true
			break
		}
	}
	
	if !found {
		fmt.Printf("❌ TestID %d not found\n", testID)
		return
	}
	
	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("TestID %d - %s (%s)\n", testID, testCase.CompType, testCase.ComparisonType)
	fmt.Println(strings.Repeat("=", 100))
	
	// Create Comp1
	comp1Cards := convertJSONToCards(testCase.Comp1.Cards, testData.Level)
	comp1 := sdk.FromCardList(comp1Cards, nil)
	
	fmt.Println("\n【Comp1】")
	fmt.Printf("  原始牌: %v\n", testCase.Comp1.Cards)
	fmt.Printf("  Cards: %s\n", printCards(comp1.GetCards()))
	fmt.Printf("  Type: %v (期望: %s)\n", comp1.GetType(), testCase.Comp1.Type)
	fmt.Printf("  Valid: %v\n", comp1.IsValid())
	if comp1.GetType().String() == "FullHouse" || comp1.GetType().String() == "Plate" || comp1.GetType().String() == "Tube" {
		if fh, ok := comp1.(*sdk.FullHouse); ok && fh.NormalizedCards != nil {
			fmt.Printf("  NormalizedCards: %s\n", printCards(fh.NormalizedCards))
		} else if p, ok := comp1.(*sdk.Plate); ok && p.NormalizedCards != nil {
			fmt.Printf("  NormalizedCards: %s\n", printCards(p.NormalizedCards))
			fmt.Printf("  ComparisonKey: %d\n", p.ComparisonKey)
		} else if t, ok := comp1.(*sdk.Tube); ok && t.NormalizedCards != nil {
			fmt.Printf("  NormalizedCards: %s\n", printCards(t.NormalizedCards))
			fmt.Printf("  ComparisonKey: %d\n", t.ComparisonKey)
		}
	}
	
	// Create Comp2
	comp2Cards := convertJSONToCards(testCase.Comp2.Cards, testData.Level)
	comp2 := sdk.FromCardList(comp2Cards, nil)
	
	fmt.Println("\n【Comp2】")
	fmt.Printf("  原始牌: %v\n", testCase.Comp2.Cards)
	fmt.Printf("  Cards: %s\n", printCards(comp2.GetCards()))
	fmt.Printf("  Type: %v (期望: %s)\n", comp2.GetType(), testCase.Comp2.Type)
	fmt.Printf("  Valid: %v\n", comp2.IsValid())
	if comp2.GetType().String() == "FullHouse" || comp2.GetType().String() == "Plate" || comp2.GetType().String() == "Tube" {
		if fh, ok := comp2.(*sdk.FullHouse); ok && fh.NormalizedCards != nil {
			fmt.Printf("  NormalizedCards: %s\n", printCards(fh.NormalizedCards))
		} else if p, ok := comp2.(*sdk.Plate); ok && p.NormalizedCards != nil {
			fmt.Printf("  NormalizedCards: %s\n", printCards(p.NormalizedCards))
			fmt.Printf("  ComparisonKey: %d\n", p.ComparisonKey)
		} else if t, ok := comp2.(*sdk.Tube); ok && t.NormalizedCards != nil {
			fmt.Printf("  NormalizedCards: %s\n", printCards(t.NormalizedCards))
			fmt.Printf("  ComparisonKey: %d\n", t.ComparisonKey)
		}
	}
	
	// Compare
	fmt.Println("\n【比较结果】")
	actualComp1Greater := comp1.GreaterThan(comp2)
	actualComp2Greater := comp2.GreaterThan(comp1)
	
	fmt.Printf("  期望: Comp1>Comp2=%v, Comp2>Comp1=%v\n", 
		testCase.Comp1GreaterThanComp2, testCase.Comp2GreaterThanComp1)
	fmt.Printf("  实际: Comp1>Comp2=%v, Comp2>Comp1=%v\n", 
		actualComp1Greater, actualComp2Greater)
	
	if actualComp1Greater == testCase.Comp1GreaterThanComp2 && 
	   actualComp2Greater == testCase.Comp2GreaterThanComp1 {
		fmt.Println("  结果: ✅ PASS")
	} else {
		fmt.Println("  结果: ❌ FAIL")
	}
	fmt.Println()
}

func main() {
	// Read test data
	testDataPath := filepath.Join("test-data", "comparison_test_data.json")
	data, err := os.ReadFile(testDataPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	
	var testData TestData
	if err := json.Unmarshal(data, &testData); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}
	
	fmt.Printf("Level: %d\n", testData.Level)
	fmt.Printf("Total test cases: %d\n\n", len(testData.Comparisons))
	
	// Failed test cases
	failedCases := []int{65, 190, 192, 193, 195, 807, 808}
	
	for _, testID := range failedCases {
		debugTestCase(testID, testData)
	}
}
