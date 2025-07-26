package sdk

import (
	"testing"
)

func TestNewCard(t *testing.T) {
	// 测试创建普通牌
	card, err := NewCard(3, "Spade", 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if card.Number != 3 || card.Color != "Spade" || card.Level != 2 {
		t.Errorf("Card creation failed: %v", card)
	}
	if card.Name != "3" {
		t.Errorf("Expected name '3', got '%s'", card.Name)
	}

	// 测试创建 Ace (A -> 14)
	card, err = NewCard(1, "Heart", 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if card.Number != 14 || card.RawNumber != 1 {
		t.Errorf("Ace conversion failed: Number=%d, RawNumber=%d", card.Number, card.RawNumber)
	}
	if card.Name != "Ace" {
		t.Errorf("Expected name 'Ace', got '%s'", card.Name)
	}

	// 测试创建 Jack
	card, err = NewCard(11, "Diamond", 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if card.Name != "Jack" {
		t.Errorf("Expected name 'Jack', got '%s'", card.Name)
	}

	// 测试创建大王
	card, err = NewCard(16, "Joker", 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if card.Name != "Red Joker" {
		t.Errorf("Expected name 'Red Joker', got '%s'", card.Name)
	}

	// 测试无效输入
	_, err = NewCard(17, "Spade", 2)
	if err == nil {
		t.Error("Expected error for invalid number")
	}

	_, err = NewCard(3, "InvalidColor", 2)
	if err == nil {
		t.Error("Expected error for invalid color")
	}

	_, err = NewCard(15, "Spade", 2)
	if err == nil {
		t.Error("Expected error for joker with wrong color")
	}
}

func TestClone(t *testing.T) {
	original, _ := NewCard(5, "Heart", 2)
	cloned := original.Clone()

	if original.Number != cloned.Number || original.Color != cloned.Color || original.Level != cloned.Level {
		t.Error("Clone failed")
	}

	// 确保是不同的对象
	if original == cloned {
		t.Error("Clone should create a new object")
	}
}

func TestIsWildcard(t *testing.T) {
	// 测试变化牌（红桃且数字等于级别）
	card, _ := NewCard(3, "Heart", 3)
	if !card.IsWildcard() {
		t.Error("Expected wildcard")
	}

	// 测试非变化牌
	card, _ = NewCard(3, "Spade", 3)
	if card.IsWildcard() {
		t.Error("Expected not wildcard")
	}

	card, _ = NewCard(4, "Heart", 3)
	if card.IsWildcard() {
		t.Error("Expected not wildcard")
	}
}

func TestGreaterThan(t *testing.T) {
	// 测试基本数字比较
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(3, "Heart", 2)
	if !card1.GreaterThan(card2) {
		t.Error("5 should be greater than 3")
	}

	// 测试级别牌
	levelCard, _ := NewCard(2, "Heart", 2)   // 级别为2的级别牌
	normalCard, _ := NewCard(13, "Spade", 2) // 普通的K
	if !levelCard.GreaterThan(normalCard) {
		t.Error("Level card should be greater than normal card")
	}

	// 测试大王小王
	bigJoker, _ := NewCard(16, "Joker", 2)
	smallJoker, _ := NewCard(15, "Joker", 2)
	if !bigJoker.GreaterThan(smallJoker) {
		t.Error("Big joker should be greater than small joker")
	}

	// 测试王大于级别牌
	joker, _ := NewCard(15, "Joker", 2)
	levelCard2, _ := NewCard(2, "Diamond", 2)
	if !joker.GreaterThan(levelCard2) {
		t.Error("Joker should be greater than level card")
	}
}

func TestConsecutiveGreaterThan(t *testing.T) {
	// 测试顺子比较（使用原始数字）
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(3, "Heart", 2)
	if !card1.ConsecutiveGreaterThan(card2) {
		t.Error("5 should be consecutive greater than 3")
	}

	// 测试 Ace 的特殊情况
	ace, _ := NewCard(1, "Heart", 2) // A -> Number=14, RawNumber=1
	two, _ := NewCard(2, "Spade", 2) // 2 -> Number=2, RawNumber=2
	if ace.ConsecutiveGreaterThan(two) {
		t.Error("Ace (raw=1) should not be consecutive greater than 2")
	}
}

func TestLessThan(t *testing.T) {
	card1, _ := NewCard(3, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	if !card1.LessThan(card2) {
		t.Error("3 should be less than 5")
	}

	// 测试相等但红桃更大的情况
	card3, _ := NewCard(5, "Heart", 2)
	card4, _ := NewCard(5, "Spade", 2)
	if !card4.LessThan(card3) {
		t.Error("Same number but Heart should be greater")
	}
}

func TestEquals(t *testing.T) {
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	card3, _ := NewCard(3, "Spade", 2)

	if !card1.Equals(card2) {
		t.Error("Cards with same number should be equal")
	}

	if card1.Equals(card3) {
		t.Error("Cards with different numbers should not be equal")
	}
}

func TestString(t *testing.T) {
	// 测试普通牌的字符串表示
	card, _ := NewCard(3, "Spade", 2)
	expected := "3 of Spade"
	if card.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, card.String())
	}

	// 测试人头牌的字符串表示
	card, _ = NewCard(11, "Heart", 2)
	expected = "Jack of Heart"
	if card.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, card.String())
	}

	// 测试大王的字符串表示
	card, _ = NewCard(16, "Joker", 2)
	expected = "Red Joker"
	if card.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, card.String())
	}
}

func TestContains(t *testing.T) {
	slice := []string{"Spade", "Heart", "Diamond", "Club"}

	if !contains(slice, "Spade") {
		t.Error("Should contain 'Spade'")
	}

	if contains(slice, "Joker") {
		t.Error("Should not contain 'Joker'")
	}
}

func TestParseCardFromID(t *testing.T) {
	// Test valid card IDs
	testCases := []struct {
		cardID   string
		level    int
		expected Card
	}{
		{"Heart_5", 2, Card{Number: 5, Color: "Heart", Level: 2, Name: "5", RawNumber: 5}},
		{"Spade_14", 3, Card{Number: 14, Color: "Spade", Level: 3, Name: "Ace", RawNumber: 1}},
		{"Joker_15", 2, Card{Number: 15, Color: "Joker", Level: 2, Name: "Black Joker", RawNumber: 15}},
		{"Joker_16", 2, Card{Number: 16, Color: "Joker", Level: 2, Name: "Red Joker", RawNumber: 16}},
		{"Club_11", 5, Card{Number: 11, Color: "Club", Level: 5, Name: "Jack", RawNumber: 11}},
	}

	for _, tc := range testCases {
		t.Run(tc.cardID, func(t *testing.T) {
			card, err := ParseCardFromID(tc.cardID, tc.level)
			if err != nil {
				t.Errorf("ParseCardFromID(%s, %d) returned error: %v", tc.cardID, tc.level, err)
				return
			}

			if card.Number != tc.expected.Number {
				t.Errorf("Expected Number %d, got %d", tc.expected.Number, card.Number)
			}
			if card.Color != tc.expected.Color {
				t.Errorf("Expected Color %s, got %s", tc.expected.Color, card.Color)
			}
			if card.Level != tc.expected.Level {
				t.Errorf("Expected Level %d, got %d", tc.expected.Level, card.Level)
			}
			if card.Name != tc.expected.Name {
				t.Errorf("Expected Name %s, got %s", tc.expected.Name, card.Name)
			}
			if card.RawNumber != tc.expected.RawNumber {
				t.Errorf("Expected RawNumber %d, got %d", tc.expected.RawNumber, card.RawNumber)
			}

			// Test round-trip: ID -> Card -> ID
			if card.GetID() != tc.cardID {
				t.Errorf("Round-trip failed: expected %s, got %s", tc.cardID, card.GetID())
			}
		})
	}

	// Test invalid card IDs
	invalidCases := []struct {
		cardID string
		level  int
	}{
		{"", 2},               // Empty ID
		{"Heart", 2},          // Missing number
		{"Heart_5_Extra", 2},  // Too many parts
		{"Heart_17", 2},       // Invalid number
		{"InvalidColor_5", 2}, // Invalid color
		{"Heart_abc", 2},      // Non-numeric number
	}

	for _, tc := range invalidCases {
		t.Run("invalid_"+tc.cardID, func(t *testing.T) {
			_, err := ParseCardFromID(tc.cardID, tc.level)
			if err == nil {
				t.Errorf("ParseCardFromID(%s, %d) should have returned an error", tc.cardID, tc.level)
			}
		})
	}
}
