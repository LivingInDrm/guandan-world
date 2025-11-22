package sdk

import (
	"fmt"
	"testing"
)

// ============================================================================
// 测试组1: K-A钢板（循环高端）
// 约束：只有Heart级牌是万能牌，最多2个万能牌，2副牌
// ============================================================================

func TestPlateSatisfyNew_KA_Plate(t *testing.T) {
	// ========== 0个万能牌 ==========
	
	t.Run("KA_0wild_case1_standard", func(t *testing.T) {
		level := 3 // 避免K和A是万能牌
		cards := []*Card{
			mustNewCard(13, "Spade", level),   // K♠
			mustNewCard(13, "Heart", level),   // K♥
			mustNewCard(13, "Club", level),    // K♣
			mustNewCard(14, "Spade", level),   // A♠
			mustNewCard(14, "Heart", level),   // A♥
			mustNewCard(14, "Club", level),    // A♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("K-A钢板应该有效")
		}
		if key != 13 {
			t.Errorf("K-A钢板的键值应该是13，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 13, 1, "K-A钢板")
		}
		fmt.Printf("✓ K-A钢板(标准): valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("KA_0wild_case2_second_deck", func(t *testing.T) {
		level := 3
		cards := []*Card{
			mustNewCard(13, "Diamond", level), // K♦
			mustNewCard(13, "Spade", level),   // K♠ (第二副)
			mustNewCard(13, "Heart", level),   // K♥ (第二副)
			mustNewCard(14, "Diamond", level), // A♦
			mustNewCard(14, "Spade", level),   // A♠ (第二副)
			mustNewCard(14, "Heart", level),   // A♥ (第二副)
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("K-A钢板应该有效")
		}
		if key != 13 {
			t.Errorf("K-A钢板的键值应该是13，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 13, 1, "K-A钢板(第二副)")
		}
		fmt.Printf("✓ K-A钢板(第二副): valid=%v, key=%d\n", isValid, key)
	})
	
	// ========== 1个万能牌 ==========
	
	t.Run("KA_1wild_case1_missing_K", func(t *testing.T) {
		level := 13 // K是万能牌
		cards := []*Card{
			mustNewCard(13, "Spade", level),   // K♠
			mustNewCard(13, "Diamond", level), // K♦
			mustNewCard(level, "Heart", level), // 13♥ 万能牌
			mustNewCard(14, "Spade", level),   // A♠
			mustNewCard(14, "Heart", level),   // A♥
			mustNewCard(14, "Club", level),    // A♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("K-A钢板+1个万能牌应该有效")
		}
		if key != 13 {
			t.Errorf("K-A钢板的键值应该是13，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 13, 1, "K-A钢板+1万能牌(缺K)")
		}
		fmt.Printf("✓ K-A钢板+1万能牌(缺K): valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("KA_1wild_case2_missing_A", func(t *testing.T) {
		level := 3 // 3♥是万能牌
		cards := []*Card{
			mustNewCard(13, "Spade", level),   // K♠
			mustNewCard(13, "Heart", level),   // K♥
			mustNewCard(13, "Club", level),    // K♣
			mustNewCard(14, "Spade", level),   // A♠
			mustNewCard(14, "Diamond", level), // A♦
			mustNewCard(level, "Heart", level), // 3♥ 万能牌
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("K-A钢板+1个万能牌应该有效")
		}
		if key != 13 {
			t.Errorf("K-A钢板的键值应该是13，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 13, 1, "K-A钢板+1万能牌(缺A)")
		}
		fmt.Printf("✓ K-A钢板+1万能牌(缺A): valid=%v, key=%d\n", isValid, key)
	})
	
	// ========== 2个万能牌 ==========
	
	t.Run("KA_2wild_case1_missing_2K", func(t *testing.T) {
		level := 13 // K是万能牌，2副牌有2张13♥
		// 创建2张万能牌（模拟2副牌）
		wild1 := mustNewCard(level, "Heart", level)
		wild1.DeckIndex = 0
		wild2 := mustNewCard(level, "Heart", level)
		wild2.DeckIndex = 1
		
		cards := []*Card{
			mustNewCard(13, "Spade", level),   // K♠
			wild1,                              // 13♥ 万能牌1
			wild2,                              // 13♥ 万能牌2
			mustNewCard(14, "Spade", level),   // A♠
			mustNewCard(14, "Heart", level),   // A♥
			mustNewCard(14, "Club", level),    // A♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("K-A钢板+2个万能牌应该有效")
		}
		if key != 13 {
			t.Errorf("K-A钢板的键值应该是13，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 13, 1, "K-A钢板+2万能牌(缺2K)")
		}
		fmt.Printf("✓ K-A钢板+2万能牌(缺2K): valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("KA_2wild_case2_missing_2A", func(t *testing.T) {
		level := 3 // 3♥是万能牌
		wild1 := mustNewCard(level, "Heart", level)
		wild1.DeckIndex = 0
		wild2 := mustNewCard(level, "Heart", level)
		wild2.DeckIndex = 1
		
		cards := []*Card{
			mustNewCard(13, "Spade", level),   // K♠
			mustNewCard(13, "Heart", level),   // K♥
			mustNewCard(13, "Club", level),    // K♣
			mustNewCard(14, "Spade", level),   // A♠
			wild1,                              // 3♥ 万能牌1
			wild2,                              // 3♥ 万能牌2
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("K-A钢板+2个万能牌应该有效")
		}
		if key != 13 {
			t.Errorf("K-A钢板的键值应该是13，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 13, 1, "K-A钢板+2万能牌(缺2A)")
		}
		fmt.Printf("✓ K-A钢板+2万能牌(缺2A): valid=%v, key=%d\n", isValid, key)
	})
}

// ============================================================================
// 测试组2: A-2钢板（循环低端）
// ============================================================================

func TestPlateSatisfyNew_A2_Plate(t *testing.T) {
	// ========== 0个万能牌 ==========
	
	t.Run("A2_0wild_case1_standard", func(t *testing.T) {
		level := 5 // 避免A和2是万能牌
		cards := []*Card{
			mustNewCard(14, "Spade", level),   // A♠
			mustNewCard(14, "Heart", level),   // A♥
			mustNewCard(14, "Club", level),    // A♣
			mustNewCard(2, "Spade", level),    // 2♠
			mustNewCard(2, "Heart", level),    // 2♥
			mustNewCard(2, "Club", level),     // 2♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("A-2钢板应该有效")
		}
		if key != 1 {
			t.Errorf("A-2钢板的键值应该是1，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 1, 2, "A-2钢板")
		}
		fmt.Printf("✓ A-2钢板(标准): valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("A2_0wild_case2_second_deck", func(t *testing.T) {
		level := 5
		cards := []*Card{
			mustNewCard(14, "Diamond", level), // A♦
			mustNewCard(14, "Spade", level),   // A♠ (第二副)
			mustNewCard(14, "Heart", level),   // A♥ (第二副)
			mustNewCard(2, "Diamond", level),  // 2♦
			mustNewCard(2, "Spade", level),    // 2♠ (第二副)
			mustNewCard(2, "Heart", level),    // 2♥ (第二副)
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("A-2钢板应该有效")
		}
		if key != 1 {
			t.Errorf("A-2钢板的键值应该是1，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 1, 2, "A-2钢板(第二副)")
		}
		fmt.Printf("✓ A-2钢板(第二副): valid=%v, key=%d\n", isValid, key)
	})
	
	// ========== 1个万能牌 ==========
	
	t.Run("A2_1wild_case1_missing_2", func(t *testing.T) {
		level := 3 // 3♥是万能牌
		cards := []*Card{
			mustNewCard(14, "Spade", level),   // A♠
			mustNewCard(14, "Heart", level),   // A♥
			mustNewCard(14, "Club", level),    // A♣
			mustNewCard(2, "Spade", level),    // 2♠
			mustNewCard(2, "Diamond", level),  // 2♦
			mustNewCard(level, "Heart", level), // 3♥ 万能牌
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("A-2钢板+1个万能牌应该有效")
		}
		if key != 1 {
			t.Errorf("A-2钢板的键值应该是1，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 1, 2, "A-2钢板+1万能牌(缺2)")
		}
		fmt.Printf("✓ A-2钢板+1万能牌(缺2): valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("A2_1wild_case2_missing_A", func(t *testing.T) {
		level := 2 // 2♥是万能牌
		cards := []*Card{
			mustNewCard(14, "Spade", level),   // A♠
			mustNewCard(14, "Diamond", level), // A♦
			mustNewCard(level, "Heart", level), // 2♥ 万能牌（作为A）
			mustNewCard(2, "Spade", level),    // 2♠
			mustNewCard(2, "Club", level),     // 2♣
			mustNewCard(2, "Diamond", level),  // 2♦
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("A-2钢板+1个万能牌应该有效")
		}
		if key != 1 {
			t.Errorf("A-2钢板的键值应该是1，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 1, 2, "A-2钢板+1万能牌(缺A)")
		}
		fmt.Printf("✓ A-2钢板+1万能牌(缺A): valid=%v, key=%d\n", isValid, key)
	})
	
	// ========== 2个万能牌 ==========
	
	t.Run("A2_2wild_case1_missing_2x2", func(t *testing.T) {
		level := 3
		wild1 := mustNewCard(level, "Heart", level)
		wild1.DeckIndex = 0
		wild2 := mustNewCard(level, "Heart", level)
		wild2.DeckIndex = 1
		
		cards := []*Card{
			mustNewCard(14, "Spade", level),   // A♠
			mustNewCard(14, "Heart", level),   // A♥
			mustNewCard(14, "Club", level),    // A♣
			mustNewCard(2, "Spade", level),    // 2♠
			wild1,                              // 3♥ 万能牌1
			wild2,                              // 3♥ 万能牌2
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("A-2钢板+2个万能牌应该有效")
		}
		if key != 1 {
			t.Errorf("A-2钢板的键值应该是1，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 1, 2, "A-2钢板+2万能牌(缺2个2)")
		}
		fmt.Printf("✓ A-2钢板+2万能牌(缺2个2): valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("A2_2wild_case2_missing_2A", func(t *testing.T) {
		level := 7
		wild1 := mustNewCard(level, "Heart", level)
		wild1.DeckIndex = 0
		wild2 := mustNewCard(level, "Heart", level)
		wild2.DeckIndex = 1
		
		cards := []*Card{
			mustNewCard(14, "Spade", level),   // A♠
			wild1,                              // 7♥ 万能牌1
			wild2,                              // 7♥ 万能牌2
			mustNewCard(2, "Spade", level),    // 2♠
			mustNewCard(2, "Heart", level),    // 2♥
			mustNewCard(2, "Club", level),     // 2♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("A-2钢板+2个万能牌应该有效")
		}
		if key != 1 {
			t.Errorf("A-2钢板的键值应该是1，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 1, 2, "A-2钢板+2万能牌(缺2个A)")
		}
		fmt.Printf("✓ A-2钢板+2万能牌(缺2个A): valid=%v, key=%d\n", isValid, key)
	})
}

// ============================================================================
// 测试组3: 普通钢板
// ============================================================================

func TestPlateSatisfyNew_Normal_Plate(t *testing.T) {
	// ========== 0个万能牌 ==========
	
	t.Run("Normal_0wild_case1_3-4", func(t *testing.T) {
		level := 7
		cards := []*Card{
			mustNewCard(3, "Spade", level),    // 3♠
			mustNewCard(3, "Heart", level),    // 3♥
			mustNewCard(3, "Club", level),     // 3♣
			mustNewCard(4, "Spade", level),    // 4♠
			mustNewCard(4, "Heart", level),    // 4♥
			mustNewCard(4, "Club", level),     // 4♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("3-4钢板应该有效")
		}
		if key != 3 {
			t.Errorf("3-4钢板的键值应该是3，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 3, 4, "3-4钢板")
		}
		fmt.Printf("✓ 3-4钢板: valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("Normal_0wild_case2_7-8", func(t *testing.T) {
		level := 2
		cards := []*Card{
			mustNewCard(7, "Spade", level),    // 7♠
			mustNewCard(7, "Heart", level),    // 7♥
			mustNewCard(7, "Club", level),     // 7♣
			mustNewCard(8, "Spade", level),    // 8♠
			mustNewCard(8, "Heart", level),    // 8♥
			mustNewCard(8, "Club", level),     // 8♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("7-8钢板应该有效")
		}
		if key != 7 {
			t.Errorf("7-8钢板的键值应该是7，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 7, 8, "7-8钢板")
		}
		fmt.Printf("✓ 7-8钢板: valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("Normal_0wild_case3_10-J", func(t *testing.T) {
		level := 5
		cards := []*Card{
			mustNewCard(10, "Spade", level),   // 10♠
			mustNewCard(10, "Heart", level),   // 10♥
			mustNewCard(10, "Club", level),    // 10♣
			mustNewCard(11, "Spade", level),   // J♠
			mustNewCard(11, "Heart", level),   // J♥
			mustNewCard(11, "Club", level),    // J♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("10-J钢板应该有效")
		}
		if key != 10 {
			t.Errorf("10-J钢板的键值应该是10，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 10, 11, "10-J钢板")
		}
		fmt.Printf("✓ 10-J钢板: valid=%v, key=%d\n", isValid, key)
	})
	
	// ========== 1个万能牌 ==========
	
	t.Run("Normal_1wild_case1_4-5", func(t *testing.T) {
		level := 5 // 5♥是万能牌
		cards := []*Card{
			mustNewCard(4, "Spade", level),    // 4♠
			mustNewCard(4, "Heart", level),    // 4♥
			mustNewCard(4, "Club", level),     // 4♣
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Diamond", level),  // 5♦
			mustNewCard(level, "Heart", level), // 5♥ 万能牌
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("4-5钢板+1个万能牌应该有效")
		}
		if key != 4 {
			t.Errorf("4-5钢板的键值应该是4，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 4, 5, "4-5钢板+1万能牌")
		}
		fmt.Printf("✓ 4-5钢板+1万能牌: valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("Normal_1wild_case2_6-7", func(t *testing.T) {
		level := 6
		cards := []*Card{
			mustNewCard(6, "Spade", level),    // 6♠
			mustNewCard(6, "Diamond", level),  // 6♦
			mustNewCard(level, "Heart", level), // 6♥ 万能牌
			mustNewCard(7, "Spade", level),    // 7♠
			mustNewCard(7, "Heart", level),    // 7♥
			mustNewCard(7, "Club", level),     // 7♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("6-7钢板+1个万能牌应该有效")
		}
		if key != 6 {
			t.Errorf("6-7钢板的键值应该是6，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 6, 7, "6-7钢板+1万能牌")
		}
		fmt.Printf("✓ 6-7钢板+1万能牌: valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("Normal_1wild_case3_9-10", func(t *testing.T) {
		level := 10
		cards := []*Card{
			mustNewCard(9, "Spade", level),    // 9♠
			mustNewCard(9, "Heart", level),    // 9♥
			mustNewCard(9, "Club", level),     // 9♣
			mustNewCard(10, "Spade", level),   // 10♠
			mustNewCard(10, "Diamond", level), // 10♦
			mustNewCard(level, "Heart", level), // 10♥ 万能牌
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("9-10钢板+1个万能牌应该有效")
		}
		if key != 9 {
			t.Errorf("9-10钢板的键值应该是9，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 9, 10, "9-10钢板+1万能牌")
		}
		fmt.Printf("✓ 9-10钢板+1万能牌: valid=%v, key=%d\n", isValid, key)
	})
	
	// ========== 2个万能牌 ==========
	
	t.Run("Normal_2wild_case1_5-6", func(t *testing.T) {
		level := 6
		wild1 := mustNewCard(level, "Heart", level)
		wild1.DeckIndex = 0
		wild2 := mustNewCard(level, "Heart", level)
		wild2.DeckIndex = 1
		
		cards := []*Card{
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Heart", level),    // 5♥
			mustNewCard(5, "Club", level),     // 5♣
			mustNewCard(6, "Spade", level),    // 6♠
			wild1,                              // 6♥ 万能牌1
			wild2,                              // 6♥ 万能牌2
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("5-6钢板+2个万能牌应该有效")
		}
		if key != 5 {
			t.Errorf("5-6钢板的键值应该是5，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 5, 6, "5-6钢板+2万能牌")
		}
		fmt.Printf("✓ 5-6钢板+2万能牌: valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("Normal_2wild_case2_8-9", func(t *testing.T) {
		level := 3
		wild1 := mustNewCard(level, "Heart", level)
		wild1.DeckIndex = 0
		wild2 := mustNewCard(level, "Heart", level)
		wild2.DeckIndex = 1
		
		cards := []*Card{
			mustNewCard(8, "Spade", level),    // 8♠
			mustNewCard(8, "Heart", level),    // 8♥
			mustNewCard(8, "Club", level),     // 8♣
			mustNewCard(9, "Spade", level),    // 9♠
			wild1,                              // 3♥ 万能牌1
			wild2,                              // 3♥ 万能牌2
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("8-9钢板+2个万能牌应该有效")
		}
		if key != 8 {
			t.Errorf("8-9钢板的键值应该是8，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 8, 9, "8-9钢板+2万能牌")
		}
		fmt.Printf("✓ 8-9钢板+2万能牌: valid=%v, key=%d\n", isValid, key)
	})
	
	t.Run("Normal_2wild_case3_J-Q", func(t *testing.T) {
		level := 11
		wild1 := mustNewCard(level, "Heart", level)
		wild1.DeckIndex = 0
		wild2 := mustNewCard(level, "Heart", level)
		wild2.DeckIndex = 1
		
		cards := []*Card{
			mustNewCard(11, "Spade", level),   // J♠
			wild1,                              // 11♥ 万能牌1
			wild2,                              // 11♥ 万能牌2
			mustNewCard(12, "Spade", level),   // Q♠
			mustNewCard(12, "Heart", level),   // Q♥
			mustNewCard(12, "Club", level),    // Q♣
		}
		
		isValid, normalizedCards, key := plateSatisfyNew(cards)
		
		if !isValid {
			t.Errorf("J-Q钢板+2个万能牌应该有效")
		}
		if key != 11 {
			t.Errorf("J-Q钢板的键值应该是11，实际是%d", key)
		}
		if isValid {
			verifyNormalizedPlate(t, normalizedCards, 11, 12, "J-Q钢板+2万能牌")
		}
		fmt.Printf("✓ J-Q钢板+2万能牌: valid=%v, key=%d\n", isValid, key)
	})
}

// ============================================================================
// 测试组4: 非Plate情况
// ============================================================================

func TestPlateSatisfyNew_Invalid_Cases(t *testing.T) {
	// ========== 0个万能牌 ==========
	
	t.Run("Invalid_0wild_case1_not_consecutive_3-5", func(t *testing.T) {
		level := 7
		cards := []*Card{
			mustNewCard(3, "Spade", level),    // 3♠
			mustNewCard(3, "Heart", level),    // 3♥
			mustNewCard(3, "Club", level),     // 3♣
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Heart", level),    // 5♥
			mustNewCard(5, "Club", level),     // 5♣
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("3-5不连续，应该无效")
		}
		fmt.Printf("✓ 3-5不连续: 正确识别为无效\n")
	})
	
	t.Run("Invalid_0wild_case2_number_exceeds_3", func(t *testing.T) {
		level := 7
		cards := []*Card{
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Heart", level),    // 5♥
			mustNewCard(5, "Club", level),     // 5♣
			mustNewCard(5, "Diamond", level),  // 5♦ (第4张5)
			mustNewCard(6, "Spade", level),    // 6♠
			mustNewCard(6, "Heart", level),    // 6♥
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("5有4张，超过3张限制，应该无效")
		}
		fmt.Printf("✓ 数字超过3张: 正确识别为无效\n")
	})
	
	t.Run("Invalid_0wild_case3_contains_joker", func(t *testing.T) {
		level := 5
		cards := []*Card{
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Heart", level),    // 5♥
			mustNewCard(5, "Club", level),     // 5♣
			mustNewCard(15, "Joker", level),   // 小王
			mustNewCard(16, "Joker", level),   // 大王
			mustNewCard(6, "Spade", level),    // 6♠
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("包含王牌，应该无效")
		}
		fmt.Printf("✓ 包含王牌: 正确识别为无效\n")
	})
	
	t.Run("Invalid_0wild_case4_wrong_length_5cards", func(t *testing.T) {
		level := 5
		cards := []*Card{
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Heart", level),    // 5♥
			mustNewCard(5, "Club", level),     // 5♣
			mustNewCard(6, "Spade", level),    // 6♠
			mustNewCard(6, "Heart", level),    // 6♥
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("只有5张牌，应该无效")
		}
		fmt.Printf("✓ 牌数错误(5张): 正确识别为无效\n")
	})
	
	t.Run("Invalid_0wild_case5_wrong_length_7cards", func(t *testing.T) {
		level := 5
		cards := []*Card{
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Heart", level),    // 5♥
			mustNewCard(5, "Club", level),     // 5♣
			mustNewCard(6, "Spade", level),    // 6♠
			mustNewCard(6, "Heart", level),    // 6♥
			mustNewCard(6, "Club", level),     // 6♣
			mustNewCard(6, "Diamond", level),  // 6♦ (第7张)
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("有7张牌，应该无效")
		}
		fmt.Printf("✓ 牌数错误(7张): 正确识别为无效\n")
	})
	
	t.Run("Invalid_0wild_case6_nil_input", func(t *testing.T) {
		isValid, _, _ := plateSatisfyNew(nil)
		
		if isValid {
			t.Errorf("nil输入，应该无效")
		}
		fmt.Printf("✓ nil输入: 正确识别为无效\n")
	})
	
	t.Run("Invalid_0wild_case7_empty_array", func(t *testing.T) {
		cards := []*Card{}
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("空数组，应该无效")
		}
		fmt.Printf("✓ 空数组: 正确识别为无效\n")
	})
	
	// ========== 1个万能牌 ==========
	
	t.Run("Invalid_1wild_case1_not_consecutive_3-5", func(t *testing.T) {
		level := 4 // 4♥是万能牌
		cards := []*Card{
			mustNewCard(3, "Spade", level),    // 3♠
			mustNewCard(3, "Heart", level),    // 3♥
			mustNewCard(3, "Club", level),     // 3♣
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Diamond", level),  // 5♦
			mustNewCard(level, "Heart", level), // 4♥ 万能牌
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("3-5不连续，即使有万能牌也应该无效（需要3张4，只有1个万能牌）")
		}
		fmt.Printf("✓ 3-5+1万能牌(不够): 正确识别为无效\n")
	})
	
	t.Run("Invalid_1wild_case2_number_exceeds_3", func(t *testing.T) {
		level := 7
		cards := []*Card{
			mustNewCard(3, "Spade", level),    // 3♠
			mustNewCard(3, "Heart", level),    // 3♥
			mustNewCard(3, "Club", level),     // 3♣
			mustNewCard(3, "Diamond", level),  // 3♦ (第4张3)
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(level, "Heart", level), // 7♥ 万能牌
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("3有4张，即使有万能牌也应该无效")
		}
		fmt.Printf("✓ 数字超过3张+1万能牌: 正确识别为无效\n")
	})
	
	t.Run("Invalid_1wild_case3_contains_joker", func(t *testing.T) {
		level := 7 // 使用7作为level避免与其他牌冲突
		wild := mustNewCard(level, "Heart", level)
		wild.DeckIndex = 0
		
		cards := []*Card{
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Heart", level),    // 5♥
			mustNewCard(5, "Club", level),     // 5♣
			mustNewCard(15, "Joker", level),   // 小王
			mustNewCard(6, "Spade", level),    // 6♠
			wild,                               // 7♥ 万能牌
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("包含王牌，即使有万能牌也应该无效")
		}
		fmt.Printf("✓ 包含王牌+1万能牌: 正确识别为无效\n")
	})
	
	// ========== 2个万能牌 ==========
	
	t.Run("Invalid_2wild_case1_not_consecutive_2-5", func(t *testing.T) {
		level := 3
		wild1 := mustNewCard(level, "Heart", level)
		wild1.DeckIndex = 0
		wild2 := mustNewCard(level, "Heart", level)
		wild2.DeckIndex = 1
		
		cards := []*Card{
			mustNewCard(2, "Spade", level),    // 2♠
			mustNewCard(2, "Heart", level),    // 2♥
			mustNewCard(2, "Club", level),     // 2♣
			mustNewCard(5, "Spade", level),    // 5♠
			wild1,                              // 3♥ 万能牌1
			wild2,                              // 3♥ 万能牌2
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("2-5不连续，即使有2个万能牌也应该无效（需要3张3或3张4）")
		}
		fmt.Printf("✓ 2-5+2万能牌(不够): 正确识别为无效\n")
	})
	
	t.Run("Invalid_2wild_case2_number_exceeds_3", func(t *testing.T) {
		level := 7
		wild1 := mustNewCard(level, "Heart", level)
		wild1.DeckIndex = 0
		wild2 := mustNewCard(level, "Heart", level)
		wild2.DeckIndex = 1
		
		cards := []*Card{
			mustNewCard(5, "Spade", level),    // 5♠
			mustNewCard(5, "Heart", level),    // 5♥
			mustNewCard(5, "Club", level),     // 5♣
			mustNewCard(5, "Diamond", level),  // 5♦ (第4张5)
			wild1,                              // 7♥ 万能牌1
			wild2,                              // 7♥ 万能牌2
		}
		
		isValid, _, _ := plateSatisfyNew(cards)
		
		if isValid {
			t.Errorf("5有4张，即使有2个万能牌也应该无效")
		}
		fmt.Printf("✓ 数字超过3张+2万能牌: 正确识别为无效\n")
	})
}

// ============================================================================
// 辅助函数
// ============================================================================

// mustNewCard 创建卡牌，如果失败则panic
func mustNewCard(number int, color string, level int) *Card {
	card, err := NewCard(number, color, level)
	if err != nil {
		panic(fmt.Sprintf("Failed to create card: number=%d, color=%s, level=%d, err=%v", 
			number, color, level, err))
	}
	return card
}

// verifyNormalizedPlate 验证规范化后的钢板结构
// expectedNum1: 前3张牌的RawNumber（第一个数字）
// expectedNum2: 后3张牌的RawNumber（第二个数字）
// 注意：万能牌保持原始RawNumber不变，只验证非万能牌
func verifyNormalizedPlate(t *testing.T, normalizedCards []*Card, expectedNum1, expectedNum2 int, testName string) {
	if len(normalizedCards) != 6 {
		t.Errorf("%s: 规范化后应该有6张牌，实际有%d张", testName, len(normalizedCards))
		return
	}
	
	// 验证前3张是否为第一个数字（万能牌可以保持原值）
	for i := 0; i < 3; i++ {
		if !normalizedCards[i].IsWildcard() && normalizedCards[i].RawNumber != expectedNum1 {
			t.Errorf("%s: 前3张的非万能牌应该是%d，但第%d张的RawNumber是%d", 
				testName, expectedNum1, i+1, normalizedCards[i].RawNumber)
		}
	}
	
	// 验证后3张是否为第二个数字（万能牌可以保持原值）
	for i := 3; i < 6; i++ {
		if !normalizedCards[i].IsWildcard() && normalizedCards[i].RawNumber != expectedNum2 {
			t.Errorf("%s: 后3张的非万能牌应该是%d，但第%d张的RawNumber是%d", 
				testName, expectedNum2, i+1, normalizedCards[i].RawNumber)
		}
	}
}
