package sdk

import (
	"errors"
	"fmt"
)

// 定义花色常量
var Colors = []string{"Spade", "Club", "Heart", "Diamond"}

// 特殊牌的名称映射
var NameMap = map[int]string{
	11: "Jack",
	12: "Queen",
	13: "King",
	14: "Ace",
	15: "Black Joker",
	16: "Red Joker",
}

// Card 结构体定义
type Card struct {
	Number    int    // 牌的数字值 (1-16)
	RawNumber int    // 原始数字值 (用于顺子比较)
	Color     string // 花色
	Level     int    // 当前级别
	Name      string // 牌的名称
}

// NewCard 创建新的牌
func NewCard(number int, color string, level int) (*Card, error) {
	// 验证数字范围
	if number < 1 || number > 16 {
		return nil, errors.New("number must be between 1 and 16")
	}

	card := &Card{
		Level: level,
	}

	// 处理 Ace 的特殊情况
	card.RawNumber = number
	if number == 14 {
		card.RawNumber = 1
	}
	if number == 1 {
		number = 14 // Ace 转换为 14
	}

	// 验证花色
	if number >= 2 && number <= 14 {
		if !contains(Colors, color) {
			return nil, errors.New("invalid color for regular card")
		}
	} else {
		if color != "Joker" {
			return nil, errors.New("joker cards must have 'Joker' color")
		}
	}

	card.Number = number
	card.Color = color

	// 设置牌的名称
	if number >= 2 && number <= 10 {
		card.Name = fmt.Sprintf("%d", number)
	} else {
		card.Name = NameMap[number]
	}

	return card, nil
}

// Clone 克隆牌
func (c *Card) Clone() *Card {
	newCard, _ := NewCard(c.Number, c.Color, c.Level)
	return newCard
}

// IsWildcard 判断是否是变化牌（红桃且数字等于级别）
func (c *Card) IsWildcard() bool {
	return c.Number == c.Level && c.Color == "Heart"
}

// GreaterThan 比较牌的大小
func (c *Card) GreaterThan(other *Card) bool {
	// 级别牌：当前级别的特殊数字，作为除了王之外的最大数字
	if other.Number == c.Level {
		if c.Number >= 15 {
			return true
		} else {
			return false
		}
	} else {
		if c.Number == c.Level {
			if other.Number <= 14 {
				return true
			} else {
				return false
			}
		} else {
			return c.Number > other.Number
		}
	}
}

// ConsecutiveGreaterThan 按原始数字比较（用于顺子）
func (c *Card) ConsecutiveGreaterThan(other *Card) bool {
	return c.RawNumber > other.RawNumber
}

// LessThan 小于比较
func (c *Card) LessThan(other *Card) bool {
	if other.GreaterThan(c) {
		return true
	} else if c.Equals(other) && other.Color == "Heart" && c.Color != "Heart" {
		return true
	}
	return false
}

// GreaterThanOrEqual 大于比较
func (c *Card) GreaterThanOrEqual(other *Card) bool {
	if c.GreaterThan(other) {
		return true
	} else if c.Equals(other) && c.Color == "Heart" && other.Color != "Heart" {
		return true
	}
	return false
}

// Equals 判断相等
func (c *Card) Equals(other *Card) bool {
	return c.Number == other.Number
}

// String 字符串表示
func (c *Card) String() string {
	if c.Color != "Joker" {
		return fmt.Sprintf("%s of %s", c.Name, c.Color)
	} else {
		return c.Name
	}
}

// JSONEncode 转换为 JSON 格式
func (c *Card) JSONEncode() map[string]interface{} {
	return map[string]interface{}{
		"color":    c.Color,
		"number":   c.RawNumber,
		"selected": false,
	}
}

// 辅助函数：检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
