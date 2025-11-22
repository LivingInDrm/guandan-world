package sdk

// fullhouse_comp_new.go - 葫芦（FullHouse）新验证方案
//
// 核心思想：
// 1. 必须有至少2种不同number的牌（排除炸弹情况）
// 2. 如果有joker，必须是2个同number的joker（wildcard无法替换成joker）
// 3. 枚举所有可能的3+2组合，选择triple最大的
// 4. joker不能作为triple，只能作为pair
//
// normalizedCards排序规则：
// - triple普通牌 → triple wildcard → pair普通牌 → pair wildcard

// FullHouseSatisfyNew 新版葫芦验证逻辑（枚举方案）
// 参数:
//   cards: 待验证的卡牌列表
// 返回值:
//   isValid: 是否为合法的葫芦
//   normalizedCards: 规范化后的卡牌列表（3张+2张顺序，普通牌优先）
func FullHouseSatisfyNew(cards []*Card) (bool, []*Card) {
	// Step 1: 输入验证
	if cards == nil || len(cards) != FULL_HOUSE_CARD_COUNT {
		return false, sortCards(cards)
	}

	sortedCards := sortCards(cards)

	// Step 2: 分类统计
	wildcards, numberGroups := classifyFullHouseCards(sortedCards)
	wildcardCount := len(wildcards)

	// Step 3: 前置验证 - 必须有至少2种不同number
	distinctNumbers := []int{}
	for num := range numberGroups {
		distinctNumbers = append(distinctNumbers, num)
	}

	if len(distinctNumbers) < 2 {
		// 只有1种number会形成炸弹，0种全是wildcard也invalid
		return false, sortedCards
	}

	// Step 4: joker规则验证 - 如果有joker，必须是joker pair
	smallJokerCount := len(numberGroups[15])
	bigJokerCount := len(numberGroups[16])
	totalJokers := smallJokerCount + bigJokerCount

	if totalJokers > 0 {
		// 必须有joker pair（2个小王 或 2个大王）
		if smallJokerCount != 2 && bigJokerCount != 2 {
			return false, sortedCards
		}
	}

	// Step 5: 枚举所有可能的3+2组合
	validCombos := enumerateFullHouseCombos(numberGroups, distinctNumbers, wildcardCount)

	if len(validCombos) == 0 {
		return false, sortedCards
	}

	// Step 6: 选择triple最大的组合
	bestCombo := selectBestFullHouseCombo(validCombos, numberGroups)

	// Step 7: 构建normalized cards
	normalizedCards := buildFullHouseNormalizedCards(
		bestCombo.tripleNum,
		bestCombo.pairNum,
		numberGroups,
		wildcards,
	)

	return true, normalizedCards
}

// classifyFullHouseCards 分类统计卡牌
// 返回值:
//   wildcards: 所有wildcard
//   numberGroups: 按Number分组的普通牌（包括joker）
func classifyFullHouseCards(cards []*Card) ([]*Card, map[int][]*Card) {
	wildcards := []*Card{}
	numberGroups := make(map[int][]*Card)

	for _, card := range cards {
		if card.IsWildcard() {
			wildcards = append(wildcards, card)
		} else {
			num := card.Number
			numberGroups[num] = append(numberGroups[num], card)
		}
	}

	return wildcards, numberGroups
}

// fullHouseCombo 葫芦组合
type fullHouseCombo struct {
	tripleNum int
	pairNum   int
}

// enumerateFullHouseCombos 枚举所有可能的3+2组合
func enumerateFullHouseCombos(
	numberGroups map[int][]*Card,
	distinctNumbers []int,
	wildcardCount int,
) []fullHouseCombo {
	validCombos := []fullHouseCombo{}

	for _, tripleNum := range distinctNumbers {
		// joker不能作为triple（15=小王, 16=大王）
		if tripleNum == 15 || tripleNum == 16 {
			continue
		}

		for _, pairNum := range distinctNumbers {
			if tripleNum == pairNum {
				continue
			}

			tripleNormalCount := len(numberGroups[tripleNum])
			pairNormalCount := len(numberGroups[pairNum])

			// 检查能否用wildcard补足
			wildcardNeeded := (3 - tripleNormalCount) + (2 - pairNormalCount)

			// 边界检查
			if wildcardNeeded <= wildcardCount &&
				tripleNormalCount <= 3 &&
				pairNormalCount <= 2 {
				validCombos = append(validCombos, fullHouseCombo{
					tripleNum: tripleNum,
					pairNum:   pairNum,
				})
			}
		}
	}

	return validCombos
}

// selectBestFullHouseCombo 选择triple最大的组合
func selectBestFullHouseCombo(
	combos []fullHouseCombo,
	numberGroups map[int][]*Card,
) fullHouseCombo {
	if len(combos) == 0 {
		return fullHouseCombo{}
	}

	bestCombo := combos[0]
	bestTripleCard := numberGroups[bestCombo.tripleNum][0]

	for _, combo := range combos[1:] {
		tripleCard := numberGroups[combo.tripleNum][0]
		if tripleCard.GreaterThan(bestTripleCard) {
			bestCombo = combo
			bestTripleCard = tripleCard
		}
	}

	return bestCombo
}

// buildFullHouseNormalizedCards 构建规范化的葫芦牌组
// 顺序：[triple普通牌... triple wildcard...] [pair普通牌... pair wildcard...]
func buildFullHouseNormalizedCards(
	tripleNum int,
	pairNum int,
	numberGroups map[int][]*Card,
	wildcards []*Card,
) []*Card {
	result := []*Card{}

	// 1. 添加triple部分的普通牌
	tripleNormalCards := numberGroups[tripleNum]
	result = append(result, tripleNormalCards...)

	// 2. 添加triple部分的wildcard
	tripleNormalCount := len(tripleNormalCards)
	wildcardForTriple := 3 - tripleNormalCount
	wildcardUsed := 0

	for i := 0; i < wildcardForTriple && i < len(wildcards); i++ {
		result = append(result, wildcards[i])
		wildcardUsed++
	}

	// 3. 添加pair部分的普通牌
	pairNormalCards := numberGroups[pairNum]
	result = append(result, pairNormalCards...)

	// 4. 添加pair部分的wildcard
	for i := wildcardUsed; i < len(wildcards); i++ {
		result = append(result, wildcards[i])
	}

	return result
}
