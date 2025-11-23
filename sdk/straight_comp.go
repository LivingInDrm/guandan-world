package sdk

// straight_comp.go - 顺子（Straight）枚举验证方案
// 使用枚举所有可能的5连续数字组合来验证顺子，统一处理循环情况（A-2-3-4-5, 10-J-Q-K-A）

const (
	// STRAIGHT_CARDS_PER_NUMBER 顺子中每个数字需要的牌数
	STRAIGHT_CARDS_PER_NUMBER = 1
)

// allConsecutiveFives 预定义的所有可能的顺子类型（10种）
// 包括：8种普通连续 + A低位循环 + A高位循环
var allConsecutiveFives = [][]int{
	// A-2-3-4-5 (A作为1，最小顺子)
	{1, 2, 3, 4, 5},
	// 普通连续：2-3-4, 3-4-5, ..., 9-10-J-Q-K
	{2, 3, 4, 5, 6},
	{3, 4, 5, 6, 7},
	{4, 5, 6, 7, 8},
	{5, 6, 7, 8, 9},
	{6, 7, 8, 9, 10},
	{7, 8, 9, 10, 11},   // 7-8-9-10-J
	{8, 9, 10, 11, 12},  // 8-9-10-J-Q
	{9, 10, 11, 12, 13}, // 9-10-J-Q-K
	// 10-J-Q-K-A (A作为高位，最大顺子)
	{10, 11, 12, 13, 1},
}

// getAllConsecutiveFives 返回所有可能的顺子类型
func getAllConsecutiveFives() [][]int {
	return allConsecutiveFives
}

// canFormStraightWithWildcards 检查是否能用给定的牌和变化牌凑成指定的顺子组合
// 参数:
//   five: 目标顺子的5个数字（如 [10, 11, 12, 13, 1] 表示 10-J-Q-K-A）
//   cardCounts: 非变化牌的数字计数
//   wildcardCount: 可用的变化牌数量
// 返回值:
//   true: 可以凑成指定顺子; false: 不能凑成
func canFormStraightWithWildcards(five []int, cardCounts map[int]int, wildcardCount int) bool {
	needed := 0
	
	for _, num := range five {
		have := cardCounts[num]
		if have < STRAIGHT_CARDS_PER_NUMBER {
			needed += (STRAIGHT_CARDS_PER_NUMBER - have)
		}
	}
	
	return needed == wildcardCount
}

// findAllValidStraightFives 找出所有可以凑成的顺子组合
// 参数:
//   cardCounts: 非变化牌的数字计数
//   wildcardCount: 可用的变化牌数量
// 返回值:
//   所有可以凑成的顺子组合列表
func findAllValidStraightFives(cardCounts map[int]int, wildcardCount int) [][]int {
	allFives := getAllConsecutiveFives()
	validFives := make([][]int, 0, len(allFives))
	
	for _, five := range allFives {
		if canFormStraightWithWildcards(five, cardCounts, wildcardCount) {
			validFives = append(validFives, five)
		}
	}
	
	return validFives
}

// getStraightFiveComparisonKey 获取顺子组合的比较键值
// 直接返回第一个数字即可，因为在枚举的10种组合中：
// - A-2-3-4-5 [1,2,3,4,5] → 1（最小）
// - 2-3-4-5-6 [2,3,4,5,6] → 2
// - ...
// - 9-10-J-Q-K [9,10,11,12,13] → 9
// - 10-J-Q-K-A [10,11,12,13,1] → 10（最大）
// five[0] 的自然顺序就是顺子的牌力顺序
// 参数:
//   five: 顺子组合的5个数字
// 返回值:
//   比较键值，值越大牌力越大
func getStraightFiveComparisonKey(five []int) int {
	return five[0]
}

// selectBestStraightFive 从多个可能的顺子组合中选择最大的
// 参数:
//   validFives: 所有合法的顺子组合
// 返回值:
//   最大的顺子组合，如果输入为空则返回 nil
func selectBestStraightFive(validFives [][]int) []int {
	if len(validFives) == 0 {
		return nil
	}
	
	bestFive := validFives[0]
	bestKey := getStraightFiveComparisonKey(bestFive)
	
	for i := 1; i < len(validFives); i++ {
		currentKey := getStraightFiveComparisonKey(validFives[i])
		if currentKey > bestKey {
			bestKey = currentKey
			bestFive = validFives[i]
		}
	}
	
	return bestFive
}

// normalizeStraightWithFive 根据最佳顺子组合规范化卡牌顺序
// 将卡牌按照 bestFive 的顺序排列，变化牌放在对应位置
// 参数:
//   sortedCards: 已排序的卡牌列表
//   bestFive: 最佳顺子组合的5个数字
// 返回值:
//   规范化后的卡牌列表，按照 bestFive 的顺序排列
func normalizeStraightWithFive(sortedCards []*Card, bestFive []int) []*Card {
	result := make([]*Card, 0, STRAIGHT_CARD_COUNT)
	wildcardPool := make([]*Card, 0)
	
	// 第一步：遍历一次，按数字分组非变化牌，同时收集变化牌
	cardsByNum := make(map[int][]*Card)
	for _, card := range sortedCards {
		if card.IsWildcard() {
			wildcardPool = append(wildcardPool, card)
		} else {
			cardsByNum[card.RawNumber] = append(cardsByNum[card.RawNumber], card)
		}
	}
	
	// 第二步：按照 bestFive 的顺序构建结果
	for _, num := range bestFive {
		cardsOfNum := cardsByNum[num]
		
		// 添加该数字的卡牌（恰好1张）
		if len(cardsOfNum) > 0 {
			result = append(result, cardsOfNum[0])
		} else if len(wildcardPool) > 0 {
			// 用变化牌补齐
			result = append(result, wildcardPool[0])
			wildcardPool = wildcardPool[1:]
		}
	}
	
	// 验证结果长度（理论上应该总是5张）
	if len(result) != STRAIGHT_CARD_COUNT {
		// 这不应该发生，如果发生了说明逻辑有bug
		// 但为了鲁棒性，返回原始排序结果
		return sortedCards
	}
	
	return result
}

// StraightSatisfy 新版顺子验证逻辑（枚举方案）
// 通过枚举所有可能的5连续数字组合来验证顺子，修复循环情况识别问题
// 参数:
//   cards: 待验证的卡牌列表
// 返回值:
//   isValid: 是否为合法的顺子
//   normalizedCards: 规范化后的卡牌列表
//   comparisonKey: 比较键值（1-10），仅在 isValid=true 时有效，无效时为 0
func StraightSatisfy(cards []*Card) (bool, []*Card, int) {
	// 检查输入有效性
	if cards == nil || len(cards) != STRAIGHT_CARD_COUNT {
		return false, sortCards(cards), 0
	}
	
	// 排序
	sortedCards := sortCardsForConsecutive(cards)
	
	// 检查大小王
	if hasJokers(sortedCards) {
		return false, sortedCards, 0
	}
	
	// 统计变化牌数量
	wildcardCount := countWildcards(sortedCards)
	
	// 统计非变化牌的数字及数量
	cardCounts := make(map[int]int)
	for _, card := range sortedCards {
		if !card.IsWildcard() {
			cardCounts[card.RawNumber]++
		}
	}
	
	// 健壮性检查：顺子定义为5个连续数字，每个数字恰好1张
	// 因此任何数字超过1张都不可能构成合法顺子
	for _, count := range cardCounts {
		if count > STRAIGHT_CARDS_PER_NUMBER {
			return false, sortedCards, 0
		}
	}
	
	// 找出所有可能的顺子组合
	validFives := findAllValidStraightFives(cardCounts, wildcardCount)
	
	if len(validFives) == 0 {
		return false, sortedCards, 0
	}
	
	// 选择最大的组合
	bestFive := selectBestStraightFive(validFives)
	
	// 防御性检查：理论上不应该为 nil，但为了鲁棒性进行检查
	if bestFive == nil {
		return false, sortedCards, 0
	}
	
	// 计算比较键值
	comparisonKey := getStraightFiveComparisonKey(bestFive)
	
	// 用最佳组合规范化卡牌
	normalizedCards := normalizeStraightWithFive(sortedCards, bestFive)
	
	return true, normalizedCards, comparisonKey
}
