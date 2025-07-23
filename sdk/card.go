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

// GetID 返回牌的唯一标识符
func (c *Card) GetID() string {
	return fmt.Sprintf("%s_%d", c.Color, c.Number)
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

// ToShortString 简化表示，用于方便阅读的输出格式
// 格式：点数+花色首字母，如 9H, QS, SJ(小王), BJ(大王)
func (c *Card) ToShortString() string {
	// 处理王牌
	if c.Color == "Joker" {
		if c.Number == 15 {
			return "SJ" // Small Joker (黑王/小王)
		} else if c.Number == 16 {
			return "BJ" // Big Joker (红王/大王)
		}
		return "?J" // 未知王牌
	}

	// 点数转换
	var numberStr string
	switch c.Number {
	case 11:
		numberStr = "J"
	case 12:
		numberStr = "Q"
	case 13:
		numberStr = "K"
	case 14:
		numberStr = "A"
	default:
		numberStr = fmt.Sprintf("%d", c.Number)
	}

	// 花色转换（取首字母）
	var suitStr string
	switch c.Color {
	case "Spade":
		suitStr = "S"
	case "Heart":
		suitStr = "H"
	case "Diamond":
		suitStr = "D"
	case "Club":
		suitStr = "C"
	default:
		suitStr = "?"
	}

	return numberStr + suitStr
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
