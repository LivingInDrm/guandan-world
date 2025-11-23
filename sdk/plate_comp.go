package sdk

// plate_comp.go - 钢板（Plate）枚举验证方案
// 使用枚举所有可能的2连续数字组合来验证钢板，统一处理循环情况（A-2, K-A）

const (
	// PLATE_CARDS_PER_NUMBER 钢板中每个数字需要的牌数
	PLATE_CARDS_PER_NUMBER = 3
)

// allConsecutivePairs 预定义的所有可能的钢板类型（13种）
// 包括：12种普通连续 + K-A 循环
var allConsecutivePairs = [][]int{
	// 普通连续：A-2, 2-3, 3-4, ..., Q-K (12种)
	{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7},
	{7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}, {12, 13},
	// K-A 循环 (1种)
	{13, 1},
}

// getAllConsecutivePairs 返回所有可能的钢板类型
func getAllConsecutivePairs() [][]int {
	return allConsecutivePairs
}

// canFormPlateWithWildcards 检查是否能用给定的牌和万能牌凑成指定的钢板组合
// 参数:
//   pair: 目标钢板的2个数字（如 [13, 1] 表示 K-A）
//   cardCounts: 非万能牌的数字计数
//   wildcardCount: 可用的万能牌数量
// 返回值:
//   true: 可以凑成指定钢板; false: 不能凑成
func canFormPlateWithWildcards(pair []int, cardCounts map[int]int, wildcardCount int) bool {
	needed := 0
	
	for _, num := range pair {
		have := cardCounts[num]
		if have < PLATE_CARDS_PER_NUMBER {
			needed += (PLATE_CARDS_PER_NUMBER - have)
		}
	}
	
	return needed == wildcardCount
}

// findAllValidPlatePairs 找出所有可以凑成的钢板组合
// 参数:
//   cardCounts: 非万能牌的数字计数
//   wildcardCount: 可用的万能牌数量
// 返回值:
//   所有可以凑成的钢板组合列表
func findAllValidPlatePairs(cardCounts map[int]int, wildcardCount int) [][]int {
	allPairs := getAllConsecutivePairs()
	validPairs := make([][]int, 0, len(allPairs))
	
	for _, pair := range allPairs {
		if canFormPlateWithWildcards(pair, cardCounts, wildcardCount) {
			validPairs = append(validPairs, pair)
		}
	}
	
	return validPairs
}

// getPlatePairComparisonKey 获取钢板组合的比较键值
// 直接返回第一个数字即可，因为在枚举的13种组合中：
// - A-2 [1,2] → 1（最小）
// - 2-3 [2,3] → 2
// - ...
// - Q-K [12,13] → 12
// - K-A [13,1] → 13（最大）
// pair[0] 的自然顺序就是钢板的牌力顺序
// 参数:
//   pair: 钢板组合的2个数字
// 返回值:
//   比较键值，值越大牌力越大
func getPlatePairComparisonKey(pair []int) int {
	return pair[0]
}

// selectBestPlatePair 从多个可能的钢板组合中选择最大的
// 参数:
//   validPairs: 所有合法的钢板组合
// 返回值:
//   最大的钢板组合，如果输入为空则返回 nil
func selectBestPlatePair(validPairs [][]int) []int {
	if len(validPairs) == 0 {
		return nil
	}
	
	bestPair := validPairs[0]
	bestKey := getPlatePairComparisonKey(bestPair)
	
	for i := 1; i < len(validPairs); i++ {
		currentKey := getPlatePairComparisonKey(validPairs[i])
		if currentKey > bestKey {
			bestKey = currentKey
			bestPair = validPairs[i]
		}
	}
	
	return bestPair
}

// normalizePlateWithPair 根据最佳钢板组合规范化卡牌顺序
// 将卡牌按照 bestPair 的顺序排列，万能牌放在对应位置
// 参数:
//   sortedCards: 已排序的卡牌列表
//   bestPair: 最佳钢板组合的2个数字
// 返回值:
//   规范化后的卡牌列表，按照 bestPair 的顺序排列
func normalizePlateWithPair(sortedCards []*Card, bestPair []int) []*Card {
	result := make([]*Card, 0, PLATE_CARD_COUNT)
	wildcardPool := make([]*Card, 0)
	
	// 第一步：遍历一次，按数字分组非万能牌，同时收集万能牌
	cardsByNum := make(map[int][]*Card)
	for _, card := range sortedCards {
		if card.IsWildcard() {
			wildcardPool = append(wildcardPool, card)
		} else {
			cardsByNum[card.RawNumber] = append(cardsByNum[card.RawNumber], card)
		}
	}
	
	// 第二步：按照 bestPair 的顺序构建结果
	for _, num := range bestPair {
		cardsOfNum := cardsByNum[num]
		
		// 添加该数字的卡牌（恰好3张）
		for i := 0; i < PLATE_CARDS_PER_NUMBER; i++ {
			if i < len(cardsOfNum) {
				result = append(result, cardsOfNum[i])
			} else if len(wildcardPool) > 0 {
				// 用万能牌补齐
				result = append(result, wildcardPool[0])
				wildcardPool = wildcardPool[1:]
			}
		}
	}
	
	// 验证结果长度（理论上应该总是6张）
	if len(result) != PLATE_CARD_COUNT {
		// 这不应该发生，如果发生了说明逻辑有bug
		// 但为了鲁棒性，返回原始排序结果
		return sortedCards
	}
	
	return result
}

// PlateSatisfy 新版钢板验证逻辑（枚举方案）
// 通过枚举所有可能的2连续数字组合来验证钢板，统一处理循环情况（A-2, K-A）
// 参数:
//   cards: 待验证的卡牌列表
// 返回值:
//   isValid: 是否为合法的钢板
//   normalizedCards: 规范化后的卡牌列表
//   comparisonKey: 比较键值（1-13），仅在 isValid=true 时有效，无效时为 0
func PlateSatisfy(cards []*Card) (bool, []*Card, int) {
	// 1. 检查输入有效性
	if cards == nil || len(cards) != PLATE_CARD_COUNT {
		return false, sortCards(cards), 0
	}
	
	// 2. 排序
	sortedCards := sortCardsForConsecutive(cards)
	
	// 3. 检查大小王（钢板不允许王牌）
	if hasJokers(sortedCards) {
		return false, sortedCards, 0
	}
	
	// 4. 统计万能牌数量
	wildcardCount := countWildcards(sortedCards)
	
	// 5. 统计非万能牌的数字及数量
	cardCounts := make(map[int]int)
	for _, card := range sortedCards {
		if !card.IsWildcard() {
			cardCounts[card.RawNumber]++
		}
	}
	
	// 6. 健壮性检查：钢板定义为2个连续数字，每个数字恰好3张
	// 因此任何数字超过3张都不可能构成合法钢板
	for _, count := range cardCounts {
		if count > PLATE_CARDS_PER_NUMBER {
			return false, sortedCards, 0
		}
	}
	
	// 7. 找出所有可能的钢板组合
	validPairs := findAllValidPlatePairs(cardCounts, wildcardCount)
	
	if len(validPairs) == 0 {
		return false, sortedCards, 0
	}
	
	// 8. 选择最大的组合
	bestPair := selectBestPlatePair(validPairs)
	
	// 防御性检查：理论上不应该为 nil，但为了鲁棒性进行检查
	if bestPair == nil {
		return false, sortedCards, 0
	}
	
	// 9. 计算比较键值
	comparisonKey := getPlatePairComparisonKey(bestPair)
	
	// 10. 用最佳组合规范化卡牌
	normalizedCards := normalizePlateWithPair(sortedCards, bestPair)
	
	return true, normalizedCards, comparisonKey
}
