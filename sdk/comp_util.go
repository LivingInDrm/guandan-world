package sdk

import (
	"fmt"
	"sort"
)

// comp_util.go - 牌组公共工具函数和常量
// 提供给所有牌型使用的基础工具函数和常量定义

// 牌型卡片数量常量
const (
	STRAIGHT_CARD_COUNT   = 5
	PLATE_CARD_COUNT      = 6
	TUBE_CARD_COUNT       = 6
	FULL_HOUSE_CARD_COUNT = 5
)

// sortCards 对卡片进行排序
func sortCards(cards []*Card) []*Card {
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].LessThan(sortedCards[j])
	})
	return sortedCards
}

// sortCardsForConsecutive 专门用于连续性判断的排序（按RawNumber排序）
// 用于顺子、钢板、钢管等需要数字连续性的牌型
func sortCardsForConsecutive(cards []*Card) []*Card {
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		// 变化牌排在最后，便于处理
		if sortedCards[i].IsWildcard() && !sortedCards[j].IsWildcard() {
			return false
		}
		if !sortedCards[i].IsWildcard() && sortedCards[j].IsWildcard() {
			return true
		}
		if sortedCards[i].IsWildcard() && sortedCards[j].IsWildcard() {
			return false // 变化牌之间顺序不重要
		}
		// 普通牌按RawNumber排序（保持数学连续性）
		return sortedCards[i].RawNumber < sortedCards[j].RawNumber
	})
	return sortedCards
}

// countWildcards 统计变化牌数量
func countWildcards(cards []*Card) int {
	count := 0
	for _, card := range cards {
		if card.IsWildcard() {
			count++
		}
	}
	return count
}

// hasJokers 检查是否有王牌
func hasJokers(cards []*Card) bool {
	for _, card := range cards {
		if card.Color == "Joker" {
			return true
		}
	}
	return false
}

// extractDeckIndexes 从 Card 列表中提取 DeckIndex 列表
func extractDeckIndexes(cards []*Card) ([]int, error) {
	indexes := make([]int, len(cards))
	for i, card := range cards {
		if card == nil {
			return nil, fmt.Errorf("nil card at index %d", i)
		}
		indexes[i] = card.DeckIndex
	}
	return indexes, nil
}

// findCardsByDeckIndexes 根据 DeckIndex 列表从牌组中查找对应的原始 Card
// source: 数据源（如玩家手牌）
// deckIndexes: 要查找的 DeckIndex 列表
// 返回找到的牌列表，保持与 deckIndexes 相同的顺序
// 如果有任何 DeckIndex 找不到则返回错误
func findCardsByDeckIndexes(source []*Card, deckIndexes []int) ([]*Card, error) {
	result := make([]*Card, 0, len(deckIndexes))
	for _, deckIndex := range deckIndexes {
		found := false
		for _, card := range source {
			if card.DeckIndex == deckIndex {
				result = append(result, card)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("card with deck_index %d not found in source", deckIndex)
		}
	}
	return result, nil
}

// findCardByDeckIndex 根据单个 DeckIndex 从牌组中查找对应的原始 Card
func findCardByDeckIndex(source []*Card, deckIndex int) (*Card, error) {
	for _, card := range source {
		if card.DeckIndex == deckIndex {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with deck_index %d not found in source", deckIndex)
}
