package sdk

// tube_comp.go - 钢管（Tube）枚举验证方案
// 使用枚举所有可能的3连续数字组合来验证钢管，统一处理循环情况（Q-K-A, A-2-3）

const (
	// TUBE_CARDS_PER_NUMBER 钢管中每个数字需要的牌数
	TUBE_CARDS_PER_NUMBER = 2
)

// allConsecutiveTriples 预定义的所有可能的钢管类型（12种）
// 包括：10种普通连续 + A-2-3 循环 + Q-K-A 循环
var allConsecutiveTriples = [][]int{
	// 普通连续：2-3-4, 3-4-5, ..., 11-12-13 (J-Q-K)
	{2, 3, 4}, {3, 4, 5}, {4, 5, 6}, {5, 6, 7}, {6, 7, 8},
	{7, 8, 9}, {8, 9, 10}, {9, 10, 11}, {10, 11, 12}, {11, 12, 13},
	// A-2-3 循环
	{1, 2, 3},
	// Q-K-A 循环
	{12, 13, 1},
}

// getAllConsecutiveTriples 返回所有可能的钢管类型
func getAllConsecutiveTriples() [][]int {
	return allConsecutiveTriples
}

// canFormTubeWithWildcards 检查是否能用给定的牌和变化牌凑成指定的钢管组合
// 参数:
//   triple: 目标钢管的3个数字（如 [12, 13, 1] 表示 Q-K-A）
//   cardCounts: 非变化牌的数字计数
//   wildcardCount: 可用的变化牌数量
// 返回值:
//   true: 可以凑成指定钢管; false: 不能凑成
func canFormTubeWithWildcards(triple []int, cardCounts map[int]int, wildcardCount int) bool {
	needed := 0
	
	for _, num := range triple {
		have := cardCounts[num]
		if have < TUBE_CARDS_PER_NUMBER {
			needed += (TUBE_CARDS_PER_NUMBER - have)
		}
	}
	
	return needed == wildcardCount
}

// findAllValidTubeTriples 找出所有可以凑成的钢管组合
// 参数:
//   cardCounts: 非变化牌的数字计数
//   wildcardCount: 可用的变化牌数量
// 返回值:
//   所有可以凑成的钢管组合列表
func findAllValidTubeTriples(cardCounts map[int]int, wildcardCount int) [][]int {
	allTriples := getAllConsecutiveTriples()
	validTriples := make([][]int, 0, len(allTriples))
	
	for _, triple := range allTriples {
		if canFormTubeWithWildcards(triple, cardCounts, wildcardCount) {
			validTriples = append(validTriples, triple)
		}
	}
	
	return validTriples
}

// getTubeTripleComparisonKey 获取钢管组合的比较键值
// 直接返回第一个数字即可，因为在枚举的12种组合中：
// - A-2-3 [1,2,3] → 1（最小）
// - 2-3-4 [2,3,4] → 2
// - ...
// - J-Q-K [11,12,13] → 11
// - Q-K-A [12,13,1] → 12（最大）
// triple[0] 的自然顺序就是钢管的牌力顺序
// 参数:
//   triple: 钢管组合的3个数字
// 返回值:
//   比较键值，值越大牌力越大
func getTubeTripleComparisonKey(triple []int) int {
	return triple[0]
}

// selectBestTubeTriple 从多个可能的钢管组合中选择最大的
// 参数:
//   validTriples: 所有合法的钢管组合
// 返回值:
//   最大的钢管组合，如果输入为空则返回 nil
func selectBestTubeTriple(validTriples [][]int) []int {
	if len(validTriples) == 0 {
		return nil
	}
	
	bestTriple := validTriples[0]
	bestKey := getTubeTripleComparisonKey(bestTriple)
	
	for i := 1; i < len(validTriples); i++ {
		currentKey := getTubeTripleComparisonKey(validTriples[i])
		if currentKey > bestKey {
			bestKey = currentKey
			bestTriple = validTriples[i]
		}
	}
	
	return bestTriple
}

// normalizeTubeWithTriple 根据最佳钢管组合规范化卡牌顺序
// 将卡牌按照 bestTriple 的顺序排列，变化牌放在对应位置
// 参数:
//   sortedCards: 已排序的卡牌列表
//   bestTriple: 最佳钢管组合的3个数字
// 返回值:
//   规范化后的卡牌列表，按照 bestTriple 的顺序排列
func normalizeTubeWithTriple(sortedCards []*Card, bestTriple []int) []*Card {
	result := make([]*Card, 0, TUBE_CARD_COUNT)
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
	
	// 第二步：按照 bestTriple 的顺序构建结果
	for _, num := range bestTriple {
		cardsOfNum := cardsByNum[num]
		
		// 添加该数字的卡牌（最多2张）
		for i := 0; i < TUBE_CARDS_PER_NUMBER; i++ {
			if i < len(cardsOfNum) {
				result = append(result, cardsOfNum[i])
			} else if len(wildcardPool) > 0 {
				// 用变化牌补齐
				result = append(result, wildcardPool[0])
				wildcardPool = wildcardPool[1:]
			}
		}
	}
	
	// 验证结果长度（理论上应该总是6张）
	if len(result) != TUBE_CARD_COUNT {
		// 这不应该发生，如果发生了说明逻辑有bug
		// 但为了鲁棒性，返回原始排序结果
		return sortedCards
	}
	
	return result
}

// tubeSatisfyNew 新版钢管验证逻辑（枚举方案）
// 通过枚举所有可能的3连续数字组合来验证钢管，修复循环情况识别问题
// 参数:
//   cards: 待验证的卡牌列表
// 返回值:
//   isValid: 是否为合法的钢管
//   normalizedCards: 规范化后的卡牌列表
//   comparisonKey: 比较键值（1-12），仅在 isValid=true 时有效，无效时为 0
func tubeSatisfyNew(cards []*Card) (bool, []*Card, int) {
	// 检查输入有效性
	if cards == nil || len(cards) != TUBE_CARD_COUNT {
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
	
	// 健壮性检查：钢管定义为3个连续数字，每个数字恰好2张
	// 因此任何数字超过2张都不可能构成合法钢管
	for _, count := range cardCounts {
		if count > TUBE_CARDS_PER_NUMBER {
			return false, sortedCards, 0
		}
	}
	
	// 找出所有可能的钢管组合
	validTriples := findAllValidTubeTriples(cardCounts, wildcardCount)
	
	if len(validTriples) == 0 {
		return false, sortedCards, 0
	}
	
	// 选择最大的组合
	bestTriple := selectBestTubeTriple(validTriples)
	
	// 防御性检查：理论上不应该为 nil，但为了鲁棒性进行检查
	if bestTriple == nil {
		return false, sortedCards, 0
	}
	
	// 计算比较键值
	comparisonKey := getTubeTripleComparisonKey(bestTriple)
	
	// 用最佳组合规范化卡牌
	normalizedCards := normalizeTubeWithTriple(sortedCards, bestTriple)
	
	return true, normalizedCards, comparisonKey
}
