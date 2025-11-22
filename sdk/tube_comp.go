package sdk

// tube_comp.go - 钢管（Tube）枚举验证方案
// 使用枚举所有可能的3连续数字组合来验证钢管，统一处理循环情况（Q-K-A, A-2-3）

// getAllConsecutiveTriples 返回所有可能的钢管类型（12种）
// 包括：10种普通连续 + A-2-3 循环 + Q-K-A 循环
func getAllConsecutiveTriples() [][]int {
	triples := make([][]int, 0, 12)
	
	// 普通连续：2-3-4, 3-4-5, ..., 11-12-13 (J-Q-K)
	for start := 2; start <= 11; start++ {
		triples = append(triples, []int{start, start + 1, start + 2})
	}
	
	// A-2-3 循环
	triples = append(triples, []int{1, 2, 3})
	
	// Q-K-A 循环
	triples = append(triples, []int{12, 13, 1})
	
	return triples
}

// canFormTubeWithWildcards 检查是否能用给定的牌和变化牌凑成指定的钢管组合
// triple: 目标钢管的3个数字（如 [12, 13, 1] 表示 Q-K-A）
// cardCounts: 非变化牌的数字计数
// wildcardCount: 可用的变化牌数量
func canFormTubeWithWildcards(triple []int, cardCounts map[int]int, wildcardCount int) bool {
	needed := 0
	
	for _, num := range triple {
		have := cardCounts[num]
		if have < 2 {
			needed += (2 - have)
		}
	}
	
	return needed == wildcardCount
}

// findAllValidTubeTriples 找出所有可以凑成的钢管组合
func findAllValidTubeTriples(cardCounts map[int]int, wildcardCount int) [][]int {
	allTriples := getAllConsecutiveTriples()
	validTriples := make([][]int, 0)
	
	for _, triple := range allTriples {
		if canFormTubeWithWildcards(triple, cardCounts, wildcardCount) {
			validTriples = append(validTriples, triple)
		}
	}
	
	return validTriples
}

// getTubeTripleComparisonKey 获取钢管组合的比较键值
// Q-K-A 最大（返回15），A-2-3 最小（返回1），其他返回第一个数字
func getTubeTripleComparisonKey(triple []int) int {
	// Q-K-A 循环：[12, 13, 1]
	if len(triple) == 3 && triple[0] == 12 && triple[1] == 13 && triple[2] == 1 {
		return 15 // 最大
	}
	
	// A-2-3 循环：[1, 2, 3]
	if len(triple) == 3 && triple[0] == 1 && triple[1] == 2 && triple[2] == 3 {
		return 1 // 最小
	}
	
	// 普通钢管：返回第一个数字
	return triple[0]
}

// selectBestTubeTriple 从多个可能的钢管组合中选择最大的
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
func normalizeTubeWithTriple(sortedCards []*Card, bestTriple []int, cardCounts map[int]int) []*Card {
	result := make([]*Card, 0, 6)
	wildcardPool := make([]*Card, 0)
	
	// 收集变化牌
	for _, card := range sortedCards {
		if card.IsWildcard() {
			wildcardPool = append(wildcardPool, card)
		}
	}
	
	// 按照 bestTriple 的顺序构建结果
	for _, num := range bestTriple {
		// 收集该数字的所有非变化牌
		cardsOfNum := make([]*Card, 0)
		for _, card := range sortedCards {
			if !card.IsWildcard() && card.RawNumber == num {
				cardsOfNum = append(cardsOfNum, card)
			}
		}
		
		// 添加该数字的卡牌（最多2张）
		for i := 0; i < 2; i++ {
			if i < len(cardsOfNum) {
				result = append(result, cardsOfNum[i])
			} else if len(wildcardPool) > 0 {
				// 用变化牌补齐
				result = append(result, wildcardPool[0])
				wildcardPool = wildcardPool[1:]
			}
		}
	}
	
	return result
}

// tubeSatisfyNew 新版钢管验证逻辑（枚举方案）
// 通过枚举所有可能的3连续数字组合来验证钢管，修复循环情况识别问题
func tubeSatisfyNew(cards []*Card) (bool, []*Card) {
	// 检查长度
	if len(cards) != TUBE_CARD_COUNT {
		return failWithSortedCards(cards)
	}
	
	// 排序
	sortedCards := sortCardsForConsecutive(cards)
	
	// 检查大小王
	if hasJokers(sortedCards) {
		return false, sortedCards
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
	
	// 找出所有可能的钢管组合
	validTriples := findAllValidTubeTriples(cardCounts, wildcardCount)
	
	if len(validTriples) == 0 {
		return false, sortedCards
	}
	
	// 选择最大的组合
	bestTriple := selectBestTubeTriple(validTriples)
	
	// 用最佳组合规范化卡牌
	normalizedCards := normalizeTubeWithTriple(sortedCards, bestTriple, cardCounts)
	
	return true, normalizedCards
}
