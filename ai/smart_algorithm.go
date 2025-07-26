package ai

import (
	"guandan-world/sdk"
	"sort"
)

// SmartAutoPlayAlgorithm 智能自动出牌算法实现
type SmartAutoPlayAlgorithm struct {
	level int // 当前级别
}

// NewSmartAutoPlayAlgorithm 创建智能算法实例
func NewSmartAutoPlayAlgorithm(level int) AutoPlayAlgorithm {
	return &SmartAutoPlayAlgorithm{
		level: level,
	}
}

// CardGroup 表示一个可能的牌组
type CardGroup struct {
	Cards       []*sdk.Card    // 牌组中的牌
	CompType    sdk.CompType   // 牌型
	Comp        sdk.CardComp   // 牌组对象
	Score       float64        // 综合评分
	CardCount   int            // 牌数量
	Damage      float64        // 破坏度
	Strength    float64        // 牌力大小（越小越好）
}

// CardPattern 用于记录牌型识别的模式
type CardPattern struct {
	Number int         // 牌面数字
	Cards  []*sdk.Card // 该数字的所有牌
}

// SelectCardsToPlay 实现自动出牌逻辑
func (algo *SmartAutoPlayAlgorithm) SelectCardsToPlay(hand []*sdk.Card, trickInfo *sdk.TrickInfo) []*sdk.Card {
	if hand == nil || len(hand) == 0 {
		return nil
	}

	// 如果是首出，选择最佳牌型
	if trickInfo.IsLeader {
		return algo.selectBestFirstPlay(hand)
	}

	// 如果不是首出，尝试跟牌
	return algo.tryToFollowSmart(hand, trickInfo)
}

// selectBestFirstPlay 智能选择最佳首出牌型
func (algo *SmartAutoPlayAlgorithm) selectBestFirstPlay(hand []*sdk.Card) []*sdk.Card {
	// 1. 识别所有可能的牌组（除单牌外）
	groups := algo.identifyAllPossibleGroups(hand)
	
	// 2. 计算每个牌组的破坏度
	algo.calculateDamageForGroups(groups, hand)
	
	// 3. 计算每个牌组的综合评分
	algo.calculateScoresForGroups(groups)
	
	// 4. 选择评分最高的牌组
	bestGroup := algo.selectBestGroup(groups)
	
	// 5. 如果没有找到合适的牌组，返回最小的单牌
	if bestGroup == nil {
		return []*sdk.Card{algo.findSmallestCard(hand)}
	}
	
	return bestGroup.Cards
}

// identifyAllPossibleGroups 识别所有可能的牌组（除单牌外）
func (algo *SmartAutoPlayAlgorithm) identifyAllPossibleGroups(hand []*sdk.Card) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	// 首先对手牌进行分组
	patterns := algo.groupCardsByNumber(hand)
	
	// 1. 识别对子
	groups = append(groups, algo.findAllPairs(patterns)...)
	
	// 2. 识别三张
	groups = append(groups, algo.findAllTriplets(patterns)...)
	
	// 3. 识别炸弹（四张及以上）
	groups = append(groups, algo.findAllBombs(patterns)...)
	
	// 4. 识别顺子
	groups = append(groups, algo.findAllStraights(hand)...)
	
	// 5. 识别同花顺
	groups = append(groups, algo.findAllStraightFlushes(hand)...)
	
	// 6. 识别葫芦（三带二）
	groups = append(groups, algo.findAllFullHouses(patterns)...)
	
	// 7. 识别飞机（连续的三张）
	groups = append(groups, algo.findAllPlanes(patterns)...)
	
	// 8. 识别钢板（连续的三张）
	groups = append(groups, algo.findAllPlates(patterns)...)
	
	// 9. 识别钢管（连续的对子）
	groups = append(groups, algo.findAllTubes(patterns)...)
	
	return groups
}

// groupCardsByNumber 按数字对牌进行分组
func (algo *SmartAutoPlayAlgorithm) groupCardsByNumber(hand []*sdk.Card) []CardPattern {
	// 使用map进行分组
	groupMap := make(map[int][]*sdk.Card)
	for _, card := range hand {
		groupMap[card.Number] = append(groupMap[card.Number], card)
	}
	
	// 转换为CardPattern切片并排序
	patterns := make([]CardPattern, 0, len(groupMap))
	for number, cards := range groupMap {
		patterns = append(patterns, CardPattern{
			Number: number,
			Cards:  cards,
		})
	}
	
	// 按数字大小排序
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Number < patterns[j].Number
	})
	
	return patterns
}

// findAllPairs 找出所有对子
func (algo *SmartAutoPlayAlgorithm) findAllPairs(patterns []CardPattern) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	for _, pattern := range patterns {
		if len(pattern.Cards) >= 2 {
			// 对于每种可能的对子组合
			for i := 0; i < len(pattern.Cards)-1; i++ {
				for j := i + 1; j < len(pattern.Cards); j++ {
					cards := []*sdk.Card{pattern.Cards[i], pattern.Cards[j]}
					comp := sdk.FromCardList(cards, nil)
					if comp != nil && comp.GetType() == sdk.TypePair {
						groups = append(groups, &CardGroup{
							Cards:     cards,
							CompType:  sdk.TypePair,
							Comp:      comp,
							CardCount: 2,
							Strength:  float64(pattern.Number),
						})
					}
				}
			}
		}
	}
	
	return groups
}

// findAllTriplets 找出所有三张
func (algo *SmartAutoPlayAlgorithm) findAllTriplets(patterns []CardPattern) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	for _, pattern := range patterns {
		if len(pattern.Cards) >= 3 {
			// 对于每种可能的三张组合
			for i := 0; i < len(pattern.Cards)-2; i++ {
				for j := i + 1; j < len(pattern.Cards)-1; j++ {
					for k := j + 1; k < len(pattern.Cards); k++ {
						cards := []*sdk.Card{pattern.Cards[i], pattern.Cards[j], pattern.Cards[k]}
						comp := sdk.FromCardList(cards, nil)
						if comp != nil && comp.GetType() == sdk.TypeTriple {
							groups = append(groups, &CardGroup{
								Cards:     cards,
								CompType:  sdk.TypeTriple,
								Comp:      comp,
								CardCount: 3,
								Strength:  float64(pattern.Number),
							})
						}
					}
				}
			}
		}
	}
	
	return groups
}

// findAllBombs 找出所有炸弹
func (algo *SmartAutoPlayAlgorithm) findAllBombs(patterns []CardPattern) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	for _, pattern := range patterns {
		if len(pattern.Cards) >= 4 {
			// 四张炸弹
			for i := 0; i <= len(pattern.Cards)-4; i++ {
				cards := pattern.Cards[i:i+4]
				comp := sdk.FromCardList(cards, nil)
				if comp != nil && comp.IsBomb() {
					groups = append(groups, &CardGroup{
						Cards:     cards,
						CompType:  sdk.TypeNaiveBomb,
						Comp:      comp,
						CardCount: 4,
						Strength:  float64(pattern.Number) * 10, // 炸弹牌力更高
					})
				}
			}
			
			// 五张及以上的炸弹
			for count := 5; count <= len(pattern.Cards); count++ {
				cards := pattern.Cards[:count]
				comp := sdk.FromCardList(cards, nil)
				if comp != nil && comp.IsBomb() {
					groups = append(groups, &CardGroup{
						Cards:     cards,
						CompType:  sdk.TypeNaiveBomb,
						Comp:      comp,
						CardCount: count,
						Strength:  float64(pattern.Number) * 10 * float64(count) / 4, // 更多张数的炸弹更强
					})
				}
			}
		}
	}
	
	// 检查王炸
	jokers := make([]*sdk.Card, 0)
	for _, card := range patterns[0].Cards {
		if card.Color == "Joker" {
			jokers = append(jokers, card)
		}
	}
	for _, pattern := range patterns {
		for _, card := range pattern.Cards {
			if card.Color == "Joker" {
				found := false
				for _, j := range jokers {
					if j == card {
						found = true
						break
					}
				}
				if !found {
					jokers = append(jokers, card)
				}
			}
		}
	}
	
	if len(jokers) == 2 {
		comp := sdk.FromCardList(jokers, nil)
		if comp != nil && comp.IsBomb() {
			groups = append(groups, &CardGroup{
				Cards:     jokers,
				CompType:  sdk.TypeJokerBomb,
				Comp:      comp,
				CardCount: 2,
				Strength:  1000, // 王炸最强
			})
		}
	}
	
	return groups
}

// findAllStraights 找出所有顺子
func (algo *SmartAutoPlayAlgorithm) findAllStraights(hand []*sdk.Card) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	// 对手牌按RawNumber排序
	sortedHand := make([]*sdk.Card, len(hand))
	copy(sortedHand, hand)
	sort.Slice(sortedHand, func(i, j int) bool {
		return sortedHand[i].RawNumber < sortedHand[j].RawNumber
	})
	
	// 尝试找出长度从5到手牌数量的所有顺子
	for length := 5; length <= len(sortedHand); length++ {
		for start := 0; start <= len(sortedHand)-length; start++ {
			cards := sortedHand[start : start+length]
			
			// 检查是否能构成顺子
			if algo.isValidStraight(cards) {
				comp := sdk.FromCardList(cards, nil)
				if comp != nil && comp.GetType() == sdk.TypeStraight {
					avgStrength := 0.0
					for _, card := range cards {
						avgStrength += float64(card.Number)
					}
					avgStrength /= float64(len(cards))
					
					groups = append(groups, &CardGroup{
						Cards:     cards,
						CompType:  sdk.TypeStraight,
						Comp:      comp,
						CardCount: length,
						Strength:  avgStrength,
					})
				}
			}
		}
	}
	
	return groups
}

// isValidStraight 检查牌组是否能构成顺子
func (algo *SmartAutoPlayAlgorithm) isValidStraight(cards []*sdk.Card) bool {
	if len(cards) < 5 {
		return false
	}
	
	// 检查是否连续且每个位置只有一张牌
	numberCount := make(map[int]int)
	for _, card := range cards {
		numberCount[card.RawNumber]++
	}
	
	// 每个数字只能有一张
	for _, count := range numberCount {
		if count != 1 {
			return false
		}
	}
	
	// 检查连续性
	numbers := make([]int, 0, len(numberCount))
	for num := range numberCount {
		numbers = append(numbers, num)
	}
	sort.Ints(numbers)
	
	for i := 1; i < len(numbers); i++ {
		if numbers[i] != numbers[i-1]+1 {
			return false
		}
	}
	
	return true
}

// findAllStraightFlushes 找出所有同花顺
func (algo *SmartAutoPlayAlgorithm) findAllStraightFlushes(hand []*sdk.Card) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	// 按花色分组
	colorGroups := make(map[string][]*sdk.Card)
	for _, card := range hand {
		if card.Color != "Joker" {
			colorGroups[card.Color] = append(colorGroups[card.Color], card)
		}
	}
	
	// 在每个花色组中找顺子
	for _, cards := range colorGroups {
		if len(cards) >= 5 {
			straights := algo.findAllStraights(cards)
			for _, straight := range straights {
				comp := sdk.FromCardList(straight.Cards, nil)
				if comp != nil && comp.GetType() == sdk.TypeStraightFlush {
					straight.CompType = sdk.TypeStraightFlush
					straight.Strength *= 2 // 同花顺比普通顺子强
					groups = append(groups, straight)
				}
			}
		}
	}
	
	return groups
}

// findAllFullHouses 找出所有葫芦
func (algo *SmartAutoPlayAlgorithm) findAllFullHouses(patterns []CardPattern) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	// 找所有三张
	triplets := make([]CardPattern, 0)
	for _, pattern := range patterns {
		if len(pattern.Cards) >= 3 {
			triplets = append(triplets, pattern)
		}
	}
	
	// 找所有对子
	pairs := make([]CardPattern, 0)
	for _, pattern := range patterns {
		if len(pattern.Cards) >= 2 {
			pairs = append(pairs, pattern)
		}
	}
	
	// 组合三张和对子
	for _, triplet := range triplets {
		for _, pair := range pairs {
			if triplet.Number != pair.Number {
				// 选择三张
				tripletCards := triplet.Cards[:3]
				// 选择对子
				pairCards := pair.Cards[:2]
				
				cards := append(tripletCards, pairCards...)
				comp := sdk.FromCardList(cards, nil)
				if comp != nil && comp.GetType() == sdk.TypeFullHouse {
					groups = append(groups, &CardGroup{
						Cards:     cards,
						CompType:  sdk.TypeFullHouse,
						Comp:      comp,
						CardCount: 5,
						Strength:  float64(triplet.Number), // 以三张的大小为主
					})
				}
			}
		}
	}
	
	return groups
}

// findAllPlanes 找出所有飞机（连续的三张）
func (algo *SmartAutoPlayAlgorithm) findAllPlanes(patterns []CardPattern) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	// 找出所有有三张的牌
	tripletPatterns := make([]CardPattern, 0)
	for _, pattern := range patterns {
		if len(pattern.Cards) >= 3 {
			tripletPatterns = append(tripletPatterns, pattern)
		}
	}
	
	// 尝试找连续的三张（至少2个三张）
	for i := 0; i < len(tripletPatterns)-1; i++ {
		consecutiveCount := 1
		cards := tripletPatterns[i].Cards[:3]
		
		for j := i + 1; j < len(tripletPatterns); j++ {
			if tripletPatterns[j].Number == tripletPatterns[j-1].Number+1 {
				consecutiveCount++
				cards = append(cards, tripletPatterns[j].Cards[:3]...)
			} else {
				break
			}
		}
		
		if consecutiveCount >= 2 {
			comp := sdk.FromCardList(cards, nil)
			if comp != nil && comp.GetType() == sdk.TypeTube {
				avgStrength := 0.0
				for k := i; k < i+consecutiveCount; k++ {
					avgStrength += float64(tripletPatterns[k].Number)
				}
				avgStrength /= float64(consecutiveCount)
				
				groups = append(groups, &CardGroup{
					Cards:     cards,
					CompType:  sdk.TypeTube,
					Comp:      comp,
					CardCount: len(cards),
					Strength:  avgStrength,
				})
			}
		}
	}
	
	return groups
}

// findAllPlates 找出所有钢板（连续的三张）
func (algo *SmartAutoPlayAlgorithm) findAllPlates(patterns []CardPattern) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	// 找出所有有三张的牌
	tripletPatterns := make([]CardPattern, 0)
	for _, pattern := range patterns {
		if len(pattern.Cards) >= 3 {
			tripletPatterns = append(tripletPatterns, pattern)
		}
	}
	
	// 如果三张数量不足2个，直接返回
	if len(tripletPatterns) < 2 {
		return groups
	}
	
	// 尝试找连续的两个三张
	for i := 0; i <= len(tripletPatterns)-2; i++ {
		// 检查是否连续
		if tripletPatterns[i+1].Number == tripletPatterns[i].Number+1 {
			
			cards := make([]*sdk.Card, 0, 6)
			cards = append(cards, tripletPatterns[i].Cards[:3]...)
			cards = append(cards, tripletPatterns[i+1].Cards[:3]...)
			
			comp := sdk.FromCardList(cards, nil)
			if comp != nil && comp.GetType() == sdk.TypePlate {
				avgStrength := float64(tripletPatterns[i].Number+tripletPatterns[i+1].Number) / 2
				
				groups = append(groups, &CardGroup{
					Cards:     cards,
					CompType:  sdk.TypePlate,
					Comp:      comp,
					CardCount: 6,
					Strength:  avgStrength,
				})
			}
		}
	}
	
	return groups
}

// findAllTubes 找出所有钢管（连续的对子）
func (algo *SmartAutoPlayAlgorithm) findAllTubes(patterns []CardPattern) []*CardGroup {
	groups := make([]*CardGroup, 0)
	
	// 找出所有有对子的牌
	pairPatterns := make([]CardPattern, 0)
	for _, pattern := range patterns {
		if len(pattern.Cards) >= 2 {
			pairPatterns = append(pairPatterns, pattern)
		}
	}
	
	// 尝试找连续的三对
	for i := 0; i <= len(pairPatterns)-3; i++ {
		isConsecutive := true
		for j := i + 1; j < i+3; j++ {
			if pairPatterns[j].Number != pairPatterns[j-1].Number+1 {
				isConsecutive = false
				break
			}
		}
		
		if isConsecutive {
			cards := make([]*sdk.Card, 0, 6)
			for j := i; j < i+3; j++ {
				cards = append(cards, pairPatterns[j].Cards[:2]...)
			}
			
			comp := sdk.FromCardList(cards, nil)
			if comp != nil && comp.GetType() == sdk.TypeTube {
				avgStrength := 0.0
				for j := i; j < i+3; j++ {
					avgStrength += float64(pairPatterns[j].Number)
				}
				avgStrength /= 3
				
				groups = append(groups, &CardGroup{
					Cards:     cards,
					CompType:  sdk.TypeTube,
					Comp:      comp,
					CardCount: 6,
					Strength:  avgStrength,
				})
			}
		}
	}
	
	return groups
}

// calculateDamageForGroups 计算每个牌组的破坏度
func (algo *SmartAutoPlayAlgorithm) calculateDamageForGroups(groups []*CardGroup, hand []*sdk.Card) {
	for _, group := range groups {
		damage := 0.0
		
		// 创建一个map记录当前牌组使用的牌
		usedCards := make(map[*sdk.Card]bool)
		for _, card := range group.Cards {
			usedCards[card] = true
		}
		
		// 检查对其他牌组的破坏
		for _, otherGroup := range groups {
			if otherGroup == group {
				continue
			}
			
			// 计算重叠的牌数
			overlap := 0
			for _, card := range otherGroup.Cards {
				if usedCards[card] {
					overlap++
				}
			}
			
			// 如果有重叠，计算破坏度
			if overlap > 0 {
				// 破坏度基于被破坏牌组的类型和重叠程度
				baseDamage := float64(overlap) / float64(len(otherGroup.Cards))
				
				// 炸弹被破坏的权重更高
				if otherGroup.CompType == sdk.TypeNaiveBomb || otherGroup.CompType == sdk.TypeJokerBomb {
					baseDamage *= 5.0
				} else if otherGroup.CompType == sdk.TypeTube {
					baseDamage *= 2.5
				} else if otherGroup.CompType == sdk.TypePlate ||
					otherGroup.CompType == sdk.TypeTube ||
					otherGroup.CompType == sdk.TypeStraightFlush {
					baseDamage *= 2.0
				}
				
				damage += baseDamage
			}
		}
		
		group.Damage = damage
	}
}

// calculateScoresForGroups 计算每个牌组的综合评分
func (algo *SmartAutoPlayAlgorithm) calculateScoresForGroups(groups []*CardGroup) {
	// 权重参数
	const (
		cardCountWeight = 5.0  // 出牌数量权重
		damageWeight    = -2.0 // 破坏度权重（负值）
		strengthWeight  = -0.5 // 牌力权重（负值，因为越小越好）
	)
	
	for _, group := range groups {
		// 综合评分 = 出牌数量权重 + 破坏度权重 + 牌力权重
		group.Score = cardCountWeight*float64(group.CardCount) +
			damageWeight*group.Damage +
			strengthWeight*group.Strength/10.0
		
		// 特殊牌型加分
		switch group.CompType {
		case sdk.TypeNaiveBomb:
			// 普通炸弹不需要额外加分
		case sdk.TypePlate:
			group.Score += 5 // 钢板额外加分
		case sdk.TypeTube:
			group.Score += 3 // 钢管额外加分
		}
	}
}

// selectBestGroup 选择评分最高的牌组
func (algo *SmartAutoPlayAlgorithm) selectBestGroup(groups []*CardGroup) *CardGroup {
	if len(groups) == 0 {
		return nil
	}
	
	best := groups[0]
	for _, group := range groups[1:] {
		if group.Score > best.Score {
			best = group
		}
	}
	
	// 只有当评分为正时才返回
	if best.Score > 0 {
		return best
	}
	
	return nil
}

// tryToFollowSmart 智能跟牌
func (algo *SmartAutoPlayAlgorithm) tryToFollowSmart(hand []*sdk.Card, trickInfo *sdk.TrickInfo) []*sdk.Card {
	leadComp := trickInfo.LeadComp
	if leadComp == nil {
		return nil
	}
	
	// 识别所有可能的牌组
	groups := algo.identifyAllPossibleGroups(hand)
	
	// 筛选出能打过当前牌的牌组
	validGroups := make([]*CardGroup, 0)
	for _, group := range groups {
		if group.Comp.GreaterThan(leadComp) {
			validGroups = append(validGroups, group)
		}
	}
	
	// 如果有炸弹，检查是否值得使用
	if leadComp.IsBomb() {
		// 如果对方出的是炸弹，只用更大的炸弹跟
		bombGroups := make([]*CardGroup, 0)
		for _, group := range validGroups {
			if group.CompType == sdk.TypeNaiveBomb || group.CompType == sdk.TypeJokerBomb {
				bombGroups = append(bombGroups, group)
			}
		}
		if len(bombGroups) > 0 {
			// 选择最小的能打过的炸弹
			sort.Slice(bombGroups, func(i, j int) bool {
				return bombGroups[i].Strength < bombGroups[j].Strength
			})
			return bombGroups[0].Cards
		}
	} else {
		// 如果对方不是炸弹，优先不使用炸弹
		nonBombGroups := make([]*CardGroup, 0)
		for _, group := range validGroups {
			if group.CompType != sdk.TypeNaiveBomb && group.CompType != sdk.TypeJokerBomb {
				nonBombGroups = append(nonBombGroups, group)
			}
		}
		
		if len(nonBombGroups) > 0 {
			// 计算破坏度和评分
			algo.calculateDamageForGroups(nonBombGroups, hand)
			algo.calculateScoresForGroups(nonBombGroups)
			
			// 选择评分最高的
			best := algo.selectBestGroup(nonBombGroups)
			if best != nil {
				return best.Cards
			}
		}
		
		// 如果没有非炸弹牌能跟，考虑是否使用炸弹
		// 这里可以加入更复杂的判断逻辑
	}
	
	return nil // 过牌
}

// findSmallestCard 找最小的牌
func (algo *SmartAutoPlayAlgorithm) findSmallestCard(hand []*sdk.Card) *sdk.Card {
	if len(hand) == 0 {
		return nil
	}
	
	smallest := hand[0]
	for _, card := range hand[1:] {
		if !card.GreaterThan(smallest) {
			smallest = card
		}
	}
	
	return smallest
}

// SelectTributeCard 选择贡牌
func (algo *SmartAutoPlayAlgorithm) SelectTributeCard(hand []*sdk.Card, excludeHeartTrump bool) *sdk.Card {
	if len(hand) == 0 {
		return nil
	}
	
	// 选择最大的非关键牌
	var bestCard *sdk.Card
	for _, card := range hand {
		// 如果需要排除红桃主牌，跳过红桃主牌
		if excludeHeartTrump && card.Color == "H" && card.Number == algo.level {
			continue
		}
		
		// 避免进贡大小王
		if card.Color == "Joker" {
			continue
		}
		
		// 选择最大的牌
		if bestCard == nil || card.GreaterThan(bestCard) {
			bestCard = card
		}
	}
	
	// 如果没有找到合适的牌，选择最大的非王牌
	if bestCard == nil {
		for _, card := range hand {
			if card.Color != "Joker" {
				if bestCard == nil || card.GreaterThan(bestCard) {
					bestCard = card
				}
			}
		}
	}
	
	return bestCard
}

// SelectReturnTributeCard 选择还贡牌
func (algo *SmartAutoPlayAlgorithm) SelectReturnTributeCard(hand []*sdk.Card, receivedCard *sdk.Card) *sdk.Card {
	if len(hand) == 0 {
		return nil
	}
	
	// 选择最小的牌进行还贡
	smallestCard := hand[0]
	for _, card := range hand[1:] {
		if !card.GreaterThan(smallestCard) {
			smallestCard = card
		}
	}
	
	return smallestCard
}