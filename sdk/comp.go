package sdk

import (
	"fmt"
	"sort"
)

// 常量定义
const (
	STRAIGHT_CARD_COUNT   = 5
	PLATE_CARD_COUNT      = 6
	TUBE_CARD_COUNT       = 6
	FULL_HOUSE_CARD_COUNT = 5
)

// 预定义的数组模式常量
var (
	// 顺子模式（4张牌的相对位置）
	STRAIGHT_PATTERN_0123 = []int{0, 1, 2, 3}
	STRAIGHT_PATTERN_0124 = []int{0, 1, 2, 4}
	STRAIGHT_PATTERN_0134 = []int{0, 1, 3, 4}
	STRAIGHT_PATTERN_0234 = []int{0, 2, 3, 4}

	// 顺子模式（3张牌的相对位置）
	STRAIGHT_PATTERN_012 = []int{0, 1, 2}
	STRAIGHT_PATTERN_023 = []int{0, 2, 3}
	STRAIGHT_PATTERN_013 = []int{0, 1, 3}
	STRAIGHT_PATTERN_024 = []int{0, 2, 4}
	STRAIGHT_PATTERN_034 = []int{0, 3, 4}
	STRAIGHT_PATTERN_014 = []int{0, 1, 4}

	// 钢管模式
	TUBE_PATTERN_TRIPLET  = []int{0, 0, 1, 1, 2, 2}
	TUBE_PATTERN_0112     = []int{0, 0, 1, 1, 2}
	TUBE_PATTERN_0122     = []int{0, 0, 1, 2, 2}
	TUBE_PATTERN_1122     = []int{0, 1, 1, 2, 2}
	TUBE_PATTERN_0011     = []int{0, 0, 1, 1}
	TUBE_PATTERN_0022     = []int{0, 0, 2, 2}
	TUBE_PATTERN_0012     = []int{0, 0, 1, 2}
	TUBE_PATTERN_0112_ALT = []int{0, 1, 1, 2}
	TUBE_PATTERN_0122_ALT = []int{0, 1, 2, 2}
)

// 辅助函数
func failWithSortedCards(cards []*Card) (bool, []*Card) {
	return false, sortCards(cards)
}

func matchesPattern(diffs []int, pattern []int) bool {
	return isArrayEqual(diffs, pattern)
}

func createNewOrder(size int) []*Card {
	return make([]*Card, size)
}

func createResult() []*Card {
	return make([]*Card, 0)
}

func combineCardArrays(arrays ...[]*Card) []*Card {
	result := createResult()
	for _, arr := range arrays {
		result = append(result, arr...)
	}
	return result
}

func combineCardsAndSlices(cards []*Card, slices ...[]int) []*Card {
	result := createResult()
	for _, slice := range slices {
		start, end := slice[0], slice[1]
		if end >= 0 {
			result = append(result, cards[start:end+1]...)
		} else {
			result = append(result, cards[start])
		}
	}
	return result
}

func computeRelativeDiffs(numbers []int, count int) []int {
	diffs := make([]int, count)
	baseValue := numbers[0]
	for i := 0; i < count; i++ {
		diffs[i] = numbers[i] - baseValue
	}
	return diffs
}

func anyPatternMatches(diffs []int, patterns ...[]int) bool {
	for _, pattern := range patterns {
		if matchesPattern(diffs, pattern) {
			return true
		}
	}
	return false
}

// CardComp 牌组接口
type CardComp interface {
	GreaterThan(other CardComp) bool
	IsBomb() bool
	GetCards() []*Card
	String() string
	IsValid() bool
	GetType() CompType
}

// CompType 牌组类型
type CompType int

const (
	TypeFold CompType = iota
	TypeIllegal
	TypeSingle
	TypePair
	TypeTriple
	TypeFullHouse
	TypeStraight
	TypePlate
	TypeTube
	TypeJokerBomb
	TypeNaiveBomb
	TypeStraightFlush
)

// String 返回CompType的字符串表示
func (ct CompType) String() string {
	switch ct {
	case TypeFold:
		return "Fold"
	case TypeIllegal:
		return "IllegalComp"
	case TypeSingle:
		return "Single"
	case TypePair:
		return "Pair"
	case TypeTriple:
		return "Triple"
	case TypeFullHouse:
		return "FullHouse"
	case TypeStraight:
		return "Straight"
	case TypePlate:
		return "Plate"
	case TypeTube:
		return "Tube"
	case TypeJokerBomb:
		return "JokerBomb"
	case TypeNaiveBomb:
		return "NaiveBomb"
	case TypeStraightFlush:
		return "StraightFlush"
	default:
		return "Unknown"
	}
}

// 公共工具函数

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

// getNormalCards 获取非变化牌
func getNormalCards(cards []*Card) []*Card {
	normalCards := []*Card{}
	for _, card := range cards {
		if !card.IsWildcard() {
			normalCards = append(normalCards, card)
		}
	}
	return normalCards
}

// countCardNumbers 统计卡片数字出现次数
func countCardNumbers(cards []*Card) map[int]int {
	cardCounts := make(map[int]int)
	for _, card := range cards {
		if !card.IsWildcard() {
			cardCounts[card.Number]++
		}
	}
	return cardCounts
}

// getMaxCardNumber 获取最大卡片数字
func getMaxCardNumber(cards []*Card) int {
	maxNumber := 0
	for _, card := range cards {
		if card.Number > maxNumber {
			maxNumber = card.Number
		}
	}
	return maxNumber
}

// getMaxCardRawNumber 获取最大卡片原始数字（用于连续性检查）
func getMaxCardRawNumber(cards []*Card) int {
	maxRawNumber := 0
	for _, card := range cards {
		if !card.IsWildcard() && card.RawNumber > maxRawNumber {
			maxRawNumber = card.RawNumber
		}
	}
	return maxRawNumber
}

// BaseComp 基础牌组结构
type BaseComp struct {
	Cards           []*Card  // 原始牌组（包含万能牌）
	NormalizedCards []*Card  // 规范化牌组（万能牌已替换为具体牌）
	Valid           bool
	Type            CompType
}

// GetCards 获取牌组中的原始牌（包含万能牌）
func (b *BaseComp) GetCards() []*Card {
	return b.Cards
}

// IsValid 检查牌组是否有效
func (b *BaseComp) IsValid() bool {
	return b.Valid
}

// GetType 获取牌组类型
func (b *BaseComp) GetType() CompType {
	return b.Type
}

// String 字符串表示
func (b *BaseComp) String() string {
	return fmt.Sprintf("%v: %v", b.Type, b.Cards)
}

// FromCardList 从牌列表生成牌组
func FromCardList(cards []*Card, prev CardComp) CardComp {
	if len(cards) == 0 {
		return &Fold{BaseComp: BaseComp{Cards: cards, Valid: true, Type: TypeFold}}
	}

	// 根据牌数判断可能的牌型
	switch len(cards) {
	case 1:
		if comp := NewSingle(cards); comp.IsValid() {
			return comp
		}
		return &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}

	case 2:
		if comp := NewPair(cards); comp.IsValid() {
			return comp
		}
		return &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}

	case 3:
		if comp := NewTriple(cards); comp.IsValid() {
			return comp
		}
		return &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}

	case 4:
		// 王炸 > 炸弹
		if comp := NewJokerBomb(cards); comp.IsValid() {
			return comp
		}
		if comp := NewNaiveBomb(cards); comp.IsValid() {
			return comp
		}
		return &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}

	case 5:
		// 同花顺 > 炸弹 > 葫芦 > 顺子
		if comp := NewStraightFlush(cards); comp.IsValid() {
			return comp
		}
		if comp := NewNaiveBomb(cards); comp.IsValid() {
			return comp
		}
		// 优先级处理：如果前一个牌组不是顺子，优先尝试葫芦
		if prev == nil || prev.GetType() != TypeStraight {
			if comp := NewFullHouse(cards); comp.IsValid() {
				return comp
			}
		}
		if comp := NewStraight(cards); comp.IsValid() {
			return comp
		}

		if comp := NewFullHouse(cards); comp.IsValid() {
			return comp
		}

		return &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}

	case 6:
		// 炸弹 > 钢板/钢管
		if comp := NewNaiveBomb(cards); comp.IsValid() {
			return comp
		}

		// 优先级处理：如果前一个牌组不是钢板，优先尝试钢管
		if prev == nil || prev.GetType() != TypePlate {
			if comp := NewTube(cards); comp.IsValid() {
				return comp
			}
		}

		// 然后尝试钢板
		if comp := NewPlate(cards); comp.IsValid() {
			return comp
		}

		// 如果钢板失败，再试钢管（防止错过）
		if comp := NewTube(cards); comp.IsValid() {
			return comp
		}

		return &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}

	default:
		// 长度大于6的只可能是炸弹
		if comp := NewNaiveBomb(cards); comp.IsValid() {
			return comp
		}
		return &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}
	}
}

// Fold 弃牌
type Fold struct {
	BaseComp
}

func (f *Fold) GreaterThan(other CardComp) bool {
	return false
}

func (f *Fold) IsBomb() bool {
	return false
}

// IllegalComp 非法牌组
type IllegalComp struct {
	BaseComp
}

func (i *IllegalComp) GreaterThan(other CardComp) bool {
	return false
}

func (i *IllegalComp) IsBomb() bool {
	return false
}

// Single 单张
type Single struct {
	BaseComp
}

func NewSingle(cards []*Card) *Single {
	valid := len(cards) == 1
	return &Single{
		BaseComp: BaseComp{
			Cards: cards,
			Valid: valid,
			Type:  TypeSingle,
		},
	}
}

func (s *Single) GreaterThan(other CardComp) bool {
	if other.GetType() != TypeSingle {
		return false
	}
	otherSingle := other.(*Single)
	return s.Cards[0].GreaterThan(otherSingle.Cards[0])
}

func (s *Single) IsBomb() bool {
	return false
}

// Pair 对子
type Pair struct {
	BaseComp
}

func NewPair(cards []*Card) *Pair {
	valid := false
	sortedCards := sortCards(cards)
	var normalizedCards []*Card
	
	if len(cards) == 2 {
		levelCond0 := cards[0].IsWildcard() && cards[1].Color != "Joker"
		levelCond1 := cards[1].IsWildcard() && cards[0].Color != "Joker"
		valid = cards[0].Equals(cards[1]) || levelCond0 || levelCond1
		
		// 如果有效，创建规范化牌组
		if valid {
			normalizedCards = cloneCards(sortedCards)
			// 找到非万能牌作为基准
			var baseCard *Card
			for _, card := range normalizedCards {
				if !card.IsWildcard() {
					baseCard = card
					break
				}
			}
			// 将所有万能牌替换为基准牌
			if baseCard != nil {
				for i, card := range normalizedCards {
					if card.IsWildcard() {
						normalizedCards[i] = createReplacementCard(
							baseCard.Number,
							baseCard.Color,
							card.Level,
						)
					}
				}
			}
		}
	}

	return &Pair{
		BaseComp: BaseComp{
			Cards:           sortedCards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypePair,
		},
	}
}

func (p *Pair) GreaterThan(other CardComp) bool {
	if other.GetType() != TypePair {
		return false
	}
	otherPair := other.(*Pair)
	// 使用规范化牌组进行比较
	myCards := p.NormalizedCards
	if myCards == nil {
		myCards = p.Cards
	}
	otherCards := otherPair.NormalizedCards
	if otherCards == nil {
		otherCards = otherPair.Cards
	}
	return myCards[0].GreaterThan(otherCards[0])
}

func (p *Pair) IsBomb() bool {
	return false
}

// Triple 三张
type Triple struct {
	BaseComp
}

func NewTriple(cards []*Card) *Triple {
	valid := false
	sortedCards := sortCards(cards)
	var normalizedCards []*Card

	if len(cards) == 3 {
		// 如果有王，则非法
		if hasJokers(sortedCards) {
			valid = false
		} else {
			// 检查是否为三张相同或包含变化牌
			valid = true
			baseCard := sortedCards[0]
			for i := 1; i < len(sortedCards); i++ {
				if !sortedCards[i].Equals(baseCard) && !sortedCards[i].IsWildcard() {
					valid = false
					break
				}
			}

			// 如果有效，创建规范化牌组
			if valid {
				normalizedCards = cloneCards(sortedCards)
				// 找到非万能牌作为基准
				var baseNormalCard *Card
				for _, card := range normalizedCards {
					if !card.IsWildcard() {
						baseNormalCard = card
						break
					}
				}
				// 将所有万能牌替换为基准牌
				if baseNormalCard != nil {
					for i, card := range normalizedCards {
						if card.IsWildcard() {
							normalizedCards[i] = createReplacementCard(
								baseNormalCard.Number,
								baseNormalCard.Color,
								card.Level,
							)
						}
					}
				}
			}
		}
	}

	return &Triple{
		BaseComp: BaseComp{
			Cards:           sortedCards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypeTriple,
		},
	}
}

func (t *Triple) GreaterThan(other CardComp) bool {
	if other.GetType() != TypeTriple {
		return false
	}
	otherTriple := other.(*Triple)
	// 使用规范化牌组进行比较
	myCards := t.NormalizedCards
	if myCards == nil {
		myCards = t.Cards
	}
	otherCards := otherTriple.NormalizedCards
	if otherCards == nil {
		otherCards = otherTriple.Cards
	}
	return myCards[0].GreaterThan(otherCards[0])
}

func (t *Triple) IsBomb() bool {
	return false
}

// FullHouse 葫芦（三带二）
type FullHouse struct {
	BaseComp
}

func NewFullHouse(cards []*Card) *FullHouse {
	valid := false
	var sortedCards []*Card
	var normalizedCards []*Card

	if len(cards) == 5 {
		// 使用Python的satisfy逻辑
		var ok bool
		ok, sortedCards = fullHouseSatisfy(cards)
		valid = ok
		
		// 如果有效，创建规范化牌组
		if valid {
			normalizedCards = normalizeFullHouse(sortedCards)
		}
	} else {
		sortedCards = sortCards(cards)
	}

	return &FullHouse{
		BaseComp: BaseComp{
			Cards:           sortedCards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypeFullHouse,
		},
	}
}

// fullHouseSatisfy 实现Python的FullHouse.satisfy逻辑
func fullHouseSatisfy(cards []*Card) (bool, []*Card) {
	if len(cards) != FULL_HOUSE_CARD_COUNT {
		return failWithSortedCards(cards)
	}

	// 排序卡片
	sortedCards := sortCards(cards)

	// 如果最大的卡片是王牌，它必须是一对王牌
	if sortedCards[4].Color == "Joker" {
		if !sortedCards[3].Equals(sortedCards[4]) {
			return false, sortedCards
		}
		// 检查剩余的是否是一个三张
		triple := NewTriple(sortedCards[:3])
		if !triple.Valid {
			return false, sortedCards
		}
		return true, append(triple.Cards, sortedCards[3:]...)
	}

	// 统计变化牌数量
	wildcardCount := countWildcards(sortedCards)

	// 如果有两个变化牌
	if wildcardCount == 2 && sortedCards[3].IsWildcard() && sortedCards[4].IsWildcard() {
		// 检查 1 + 2 优先选择更大的葫芦
		pair := NewPair(sortedCards[1:3])
		if pair.Valid && !sortedCards[0].Equals(pair.Cards[0]) {
			result := createResult()
			result = append(result, sortedCards[0])
			result = append(result, sortedCards[3:]...)
			result = append(result, pair.Cards...)
			return true, result
		}
		// 检查 2 + 1
		pair = NewPair(sortedCards[:2])
		if pair.Valid && !sortedCards[2].Equals(pair.Cards[0]) {
			result := createResult()
			result = append(result, sortedCards[2])
			result = append(result, sortedCards[3:]...)
			result = append(result, pair.Cards...)
			return true, result
		}
		// 检查 3
		triple := NewTriple(sortedCards[:3])
		if triple.Valid {
			return true, append(triple.Cards, sortedCards[3:]...)
		}
		return failWithSortedCards(sortedCards)
	}

	// 如果没有变化牌
	if wildcardCount == 0 {
		// 可以是 3 + 2 或 2 + 3
		// 检查 2 + 3 优先选择更大的葫芦
		triple := NewTriple(sortedCards[2:])
		pair := NewPair(sortedCards[:2])
		if triple.Valid && pair.Valid && !triple.Cards[0].Equals(pair.Cards[0]) {
			return true, append(triple.Cards, pair.Cards...)
		}
		// 检查 3 + 2
		triple = NewTriple(sortedCards[:3])
		pair = NewPair(sortedCards[3:])
		if triple.Valid && pair.Valid && !triple.Cards[0].Equals(pair.Cards[0]) {
			return true, append(triple.Cards, pair.Cards...)
		}
		return failWithSortedCards(sortedCards)
	}

	// 如果有一个变化牌
	if wildcardCount == 1 {
		// 可以是 2 + 2, 3 + 1, 1 + 3
		// 检查 2 + 2
		pair1 := NewPair(sortedCards[:2])
		pair2 := NewPair(sortedCards[2:4])
		if pair1.Valid && pair2.Valid && !pair1.Cards[0].Equals(pair2.Cards[0]) {
			result := createResult()
			result = append(result, pair2.Cards...)
			result = append(result, sortedCards[4])
			result = append(result, pair1.Cards...)
			return true, result
		}
		// 检查 1 + 3 优先选择更大的葫芦
		triple := NewTriple(sortedCards[1:4])
		if triple.Valid && !triple.Cards[0].Equals(sortedCards[0]) {
			result := createResult()
			result = append(result, triple.Cards...)
			result = append(result, sortedCards[0])
			result = append(result, sortedCards[4])
			return true, result
		}
		// 检查 3 + 1
		triple = NewTriple(sortedCards[:3])
		if triple.Valid && !triple.Cards[0].Equals(sortedCards[3]) {
			result := createResult()
			result = append(result, triple.Cards...)
			result = append(result, sortedCards[3])
			result = append(result, sortedCards[4])
			return true, result
		}
		return failWithSortedCards(sortedCards)
	}

	return failWithSortedCards(sortedCards)
}

func (f *FullHouse) GreaterThan(other CardComp) bool {
	if other.GetType() != TypeFullHouse {
		return false
	}
	otherFullHouse := other.(*FullHouse)
	// 使用规范化牌组进行比较（比较三张部分）
	myCards := f.NormalizedCards
	if myCards == nil {
		myCards = f.Cards
	}
	otherCards := otherFullHouse.NormalizedCards
	if otherCards == nil {
		otherCards = otherFullHouse.Cards
	}
	return myCards[0].GreaterThan(otherCards[0])
}

func (f *FullHouse) IsBomb() bool {
	return false
}

// Straight 顺子
type Straight struct {
	BaseComp
}

func NewStraight(cards []*Card) *Straight {
	valid := false
	var sortedCards []*Card
	var normalizedCards []*Card

	if len(cards) == 5 {
		// 使用Python的satisfy逻辑
		var ok bool
		ok, sortedCards = straightSatisfy(cards)
		valid = ok
		
		// 如果有效，创建规范化牌组
		if valid {
			normalizedCards = normalizeStraight(sortedCards)
		}
	} else {
		sortedCards = sortCards(cards)
	}

	return &Straight{
		BaseComp: BaseComp{
			Cards:           sortedCards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypeStraight,
		},
	}
}

// straightSatisfy 实现Python的Straight.satisfy逻辑
func straightSatisfy(cards []*Card) (bool, []*Card) {
	if len(cards) != STRAIGHT_CARD_COUNT {
		return failWithSortedCards(cards)
	}

	// 使用连续性排序卡片
	sortedCards := sortCardsForConsecutive(cards)

	// 统计变化牌数量
	numWildcards := countWildcards(sortedCards)

	// 获取卡片数字（使用RawNumber进行连续性判断）
	cardNumbers := make([]int, len(sortedCards))
	for i, card := range sortedCards {
		if card.IsWildcard() {
			cardNumbers[i] = -1 // 变化牌标记为-1，后续处理
		} else {
			cardNumbers[i] = card.RawNumber
		}
	}

	// 最大牌不能超过A（使用RawNumber检查）
	if getMaxCardRawNumber(sortedCards) > 13 {
		return failWithSortedCards(cards)
	}

	// 没有变化牌
	if numWildcards == 0 {
		// 过滤掉变化牌标记(-1)并获取实际数字
		actualNumbers := []int{}
		for _, num := range cardNumbers {
			if num != -1 {
				actualNumbers = append(actualNumbers, num)
			}
		}

		if len(actualNumbers) == 5 && len(removeDuplicates(actualNumbers)) == 5 {
			// 检查常规顺子：如2-3-4-5-6
			if actualNumbers[4]-actualNumbers[0] == 4 {
				return true, sortedCards
			}

			// 检查A高位顺子：如10-J-Q-K-A (10,11,12,13,1)
			if len(actualNumbers) == 5 && actualNumbers[0] == 1 &&
				actualNumbers[1] == 10 && actualNumbers[2] == 11 &&
				actualNumbers[3] == 12 && actualNumbers[4] == 13 {
				// 重新排序为10-J-Q-K-A
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				for _, card := range sortedCards {
					if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 10 {
						newOrder[0] = card
					} else if card.RawNumber == 11 {
						newOrder[1] = card
					} else if card.RawNumber == 12 {
						newOrder[2] = card
					} else if card.RawNumber == 13 {
						newOrder[3] = card
					}
				}
				return true, newOrder
			}
		}
		return failWithSortedCards(cards)
	}

	// 一个变化牌
	if numWildcards == 1 {
		firstFour := computeRelativeDiffs(cardNumbers, 4)

		// 过滤掉变化牌标记(-1)并获取实际数字
		actualNumbers := []int{}
		for _, num := range cardNumbers {
			if num != -1 {
				actualNumbers = append(actualNumbers, num)
			}
		}

		if len(actualNumbers) == 4 && len(removeDuplicates(actualNumbers)) == 4 {
			// 检查各种A高位顺子的情况

			// J-Q-K-A + 变化牌 (填入10)
			if actualNumbers[0] == 1 && actualNumbers[1] == 11 && actualNumbers[2] == 12 && actualNumbers[3] == 13 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardPlaced := false
				for _, card := range sortedCards {
					if card.IsWildcard() && !wildcardPlaced {
						newOrder[0] = card // 变化牌作为10
						wildcardPlaced = true
					} else if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 11 {
						newOrder[1] = card
					} else if card.RawNumber == 12 {
						newOrder[2] = card
					} else if card.RawNumber == 13 {
						newOrder[3] = card
					}
				}
				return true, newOrder
			}

			// 10-Q-K-A + 变化牌 (填入J)
			if actualNumbers[0] == 1 && actualNumbers[1] == 10 && actualNumbers[2] == 12 && actualNumbers[3] == 13 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardPlaced := false
				for _, card := range sortedCards {
					if card.IsWildcard() && !wildcardPlaced {
						newOrder[1] = card // 变化牌作为J
						wildcardPlaced = true
					} else if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 10 {
						newOrder[0] = card
					} else if card.RawNumber == 12 {
						newOrder[2] = card
					} else if card.RawNumber == 13 {
						newOrder[3] = card
					}
				}
				return true, newOrder
			}

			// 10-J-K-A + 变化牌 (填入Q)
			if actualNumbers[0] == 1 && actualNumbers[1] == 10 && actualNumbers[2] == 11 && actualNumbers[3] == 13 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardPlaced := false
				for _, card := range sortedCards {
					if card.IsWildcard() && !wildcardPlaced {
						newOrder[2] = card // 变化牌作为Q
						wildcardPlaced = true
					} else if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 10 {
						newOrder[0] = card
					} else if card.RawNumber == 11 {
						newOrder[1] = card
					} else if card.RawNumber == 13 {
						newOrder[3] = card
					}
				}
				return true, newOrder
			}

			// 10-J-Q-A + 变化牌 (填入K)
			if actualNumbers[0] == 1 && actualNumbers[1] == 10 && actualNumbers[2] == 11 && actualNumbers[3] == 12 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardPlaced := false
				for _, card := range sortedCards {
					if card.IsWildcard() && !wildcardPlaced {
						newOrder[3] = card // 变化牌作为K
						wildcardPlaced = true
					} else if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 10 {
						newOrder[0] = card
					} else if card.RawNumber == 11 {
						newOrder[1] = card
					} else if card.RawNumber == 12 {
						newOrder[2] = card
					}
				}
				return true, newOrder
			}

			firstFour := computeRelativeDiffs(actualNumbers, 4)

			// i, i+1, i+2, i+3 wild
			if matchesPattern(firstFour, STRAIGHT_PATTERN_0123) {
				if actualNumbers[3] <= 12 {
					return true, sortedCards
				}
			}
		}

		// i, i+1, i+2, i+4 wild
		if matchesPattern(firstFour, STRAIGHT_PATTERN_0124) {
			newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
			copy(newOrder[0:3], sortedCards[0:3])
			newOrder[3] = sortedCards[4] // 变化牌
			newOrder[4] = sortedCards[3]
			return true, newOrder
		}

		// i, i+1, i+3, i+4 wild
		if matchesPattern(firstFour, STRAIGHT_PATTERN_0134) {
			newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
			copy(newOrder[0:2], sortedCards[0:2])
			newOrder[2] = sortedCards[4] // 变化牌
			copy(newOrder[3:], sortedCards[2:4])
			return true, newOrder
		}

		// i, i+2, i+3, i+4 wild
		if matchesPattern(firstFour, STRAIGHT_PATTERN_0234) {
			newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
			newOrder[0] = sortedCards[0]
			newOrder[1] = sortedCards[4] // 变化牌
			copy(newOrder[2:], sortedCards[1:4])
			return true, newOrder
		}

		return failWithSortedCards(cards)
	}

	// 两个变化牌
	if numWildcards == 2 {
		// 过滤掉变化牌标记(-1)并获取实际数字
		actualNumbers := []int{}
		for _, num := range cardNumbers {
			if num != -1 {
				actualNumbers = append(actualNumbers, num)
			}
		}

		if len(actualNumbers) == 3 && len(removeDuplicates(actualNumbers)) == 3 {
			// 检查各种A高位顺子的情况

			// Q-K-A + 两个变化牌 (填入J和10)
			if actualNumbers[0] == 1 && actualNumbers[1] == 12 && actualNumbers[2] == 13 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardIndex := 0
				for _, card := range sortedCards {
					if card.IsWildcard() {
						if wildcardIndex == 0 {
							newOrder[0] = card // 第一个变化牌作为10
						} else {
							newOrder[1] = card // 第二个变化牌作为J
						}
						wildcardIndex++
					} else if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 12 {
						newOrder[2] = card
					} else if card.RawNumber == 13 {
						newOrder[3] = card
					}
				}
				return true, newOrder
			}

			// J-K-A + 两个变化牌 (填入10和Q)
			if actualNumbers[0] == 1 && actualNumbers[1] == 11 && actualNumbers[2] == 13 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardIndex := 0
				for _, card := range sortedCards {
					if card.IsWildcard() {
						if wildcardIndex == 0 {
							newOrder[0] = card // 第一个变化牌作为10
						} else {
							newOrder[2] = card // 第二个变化牌作为Q
						}
						wildcardIndex++
					} else if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 11 {
						newOrder[1] = card
					} else if card.RawNumber == 13 {
						newOrder[3] = card
					}
				}
				return true, newOrder
			}

			// J-Q-A + 两个变化牌 (填入10和K)
			if actualNumbers[0] == 1 && actualNumbers[1] == 11 && actualNumbers[2] == 12 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardIndex := 0
				for _, card := range sortedCards {
					if card.IsWildcard() {
						if wildcardIndex == 0 {
							newOrder[0] = card // 第一个变化牌作为10
						} else {
							newOrder[3] = card // 第二个变化牌作为K
						}
						wildcardIndex++
					} else if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 11 {
						newOrder[1] = card
					} else if card.RawNumber == 12 {
						newOrder[2] = card
					}
				}
				return true, newOrder
			}

			// 10-K-A + 两个变化牌 (填入J和Q)
			if actualNumbers[0] == 1 && actualNumbers[1] == 10 && actualNumbers[2] == 13 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardIndex := 0
				for _, card := range sortedCards {
					if card.IsWildcard() {
						if wildcardIndex == 0 {
							newOrder[1] = card // 第一个变化牌作为J
						} else {
							newOrder[2] = card // 第二个变化牌作为Q
						}
						wildcardIndex++
					} else if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 10 {
						newOrder[0] = card
					} else if card.RawNumber == 13 {
						newOrder[3] = card
					}
				}
				return true, newOrder
			}

			// 10-J-A + 两个变化牌 (填入Q和K)
			if actualNumbers[0] == 1 && actualNumbers[1] == 10 && actualNumbers[2] == 11 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardIndex := 0
				for _, card := range sortedCards {
					if card.IsWildcard() {
						if wildcardIndex == 0 {
							newOrder[2] = card // 第一个变化牌作为Q
						} else {
							newOrder[3] = card // 第二个变化牌作为K
						}
						wildcardIndex++
					} else if card.RawNumber == 1 { // A
						newOrder[4] = card
					} else if card.RawNumber == 10 {
						newOrder[0] = card
					} else if card.RawNumber == 11 {
						newOrder[1] = card
					}
				}
				return true, newOrder
			}

			// J-Q-K + 两个变化牌 (填入10和A，构成10-J-Q-K-A)
			if actualNumbers[0] == 11 && actualNumbers[1] == 12 && actualNumbers[2] == 13 {
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				wildcardIndex := 0
				for _, card := range sortedCards {
					if card.IsWildcard() {
						if wildcardIndex == 0 {
							newOrder[0] = card // 第一个变化牌作为10
						} else {
							newOrder[4] = card // 第二个变化牌作为A
						}
						wildcardIndex++
					} else if card.RawNumber == 11 { // J
						newOrder[1] = card
					} else if card.RawNumber == 12 { // Q
						newOrder[2] = card
					} else if card.RawNumber == 13 { // K
						newOrder[3] = card
					}
				}
				return true, newOrder
			}

			firstThree := computeRelativeDiffs(actualNumbers, 3)

			// i, i+1, i+2, wild, wild
			if matchesPattern(firstThree, STRAIGHT_PATTERN_012) {
				if actualNumbers[2] <= 11 {
					return true, sortedCards
				}
			}

			// 处理其他二变化牌的情况
			if anyPatternMatches(firstThree, STRAIGHT_PATTERN_023, STRAIGHT_PATTERN_013, STRAIGHT_PATTERN_024, STRAIGHT_PATTERN_034, STRAIGHT_PATTERN_014) {
				return true, sortedCards
			}
		}

		return failWithSortedCards(cards)
	}

	return failWithSortedCards(cards)
}

// removeDuplicates 移除重复元素
func removeDuplicates(arr []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, num := range arr {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

// isArrayEqual 检查两个数组是否相等
func isArrayEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (s *Straight) GreaterThan(other CardComp) bool {
	if other.GetType() != TypeStraight {
		return false
	}
	otherStraight := other.(*Straight)

	// 使用规范化牌组进行比较
	myCards := s.NormalizedCards
	if myCards == nil {
		myCards = s.Cards
	}
	otherCards := otherStraight.NormalizedCards
	if otherCards == nil {
		otherCards = otherStraight.Cards
	}

	// 使用getStraightComparisonKey来正确比较顺子大小
	myKey := getStraightComparisonKey(myCards)
	otherKey := getStraightComparisonKey(otherCards)

	return myKey > otherKey
}

// getStraightComparisonKey 获取顺子的比较键值
func getStraightComparisonKey(cards []*Card) int {
	if len(cards) != STRAIGHT_CARD_COUNT {
		return 0
	}

	// cards已经通过straightSatisfy重构，直接分析重构后的序列
	// 构建数字序列，万能牌用前后推导的值填充
	sequence := make([]int, STRAIGHT_CARD_COUNT)
	hasWildcard := false

	// 先填入所有非万能牌的值
	for i, card := range cards {
		if card.IsWildcard() {
			hasWildcard = true
			sequence[i] = -1 // 标记为待填充
		} else {
			sequence[i] = card.Number
		}
	}

	// 如果有万能牌，通过前后推导填充
	if hasWildcard {
		for i := 0; i < STRAIGHT_CARD_COUNT; i++ {
			if sequence[i] == -1 {
				// 查找最近的非万能牌来推导
				expectedValue := -1
				if i > 0 && sequence[i-1] != -1 {
					expectedValue = sequence[i-1] + 1
				} else if i < STRAIGHT_CARD_COUNT-1 && sequence[i+1] != -1 {
					expectedValue = sequence[i+1] - 1
				}

				// 如果还是推导不出，使用位置推导
				if expectedValue == -1 {
					// 找到第一个非万能牌作为参考点
					for j := 0; j < STRAIGHT_CARD_COUNT; j++ {
						if sequence[j] != -1 {
							expectedValue = sequence[j] + (i - j)
							break
						}
					}
				}

				sequence[i] = expectedValue
			}
		}
	}

	// 处理A的循环性
	// 如果序列包含A(14)和小牌，判断是否为A-2-3-4-5类型
	hasAce := false
	hasSmallCard := false
	for _, num := range sequence {
		if num == 14 || num == 1 { // A可能被表示为1或14
			hasAce = true
		}
		if num >= 2 && num <= 5 {
			hasSmallCard = true
		}
	}

	if hasAce && hasSmallCard {
		// A-2-3-4-5 低端顺子
		return 1
	} else if hasAce {
		// 包含A的高端顺子，统一返回相同键值
		return 10
	}

	// 普通顺子，返回最小值
	minNum := sequence[0]
	for _, num := range sequence {
		if num < minNum {
			minNum = num
		}
	}
	return minNum
}

func (s *Straight) IsBomb() bool {
	return false
}

// Plate 钢板（连续三张）
type Plate struct {
	BaseComp
}

func NewPlate(cards []*Card) *Plate {
	valid := false
	var sortedCards []*Card
	var normalizedCards []*Card

	if len(cards) == 6 {
		// 使用Python的satisfy逻辑
		var ok bool
		ok, sortedCards = plateSatisfy(cards)
		valid = ok
		
		// 如果有效，创建规范化牌组
		if valid {
			normalizedCards = normalizePlate(sortedCards)
		}
	} else {
		sortedCards = sortCards(cards)
	}

	return &Plate{
		BaseComp: BaseComp{
			Cards:           sortedCards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypePlate,
		},
	}
}

// plateSatisfy 实现Python的Plate.satisfy逻辑
func plateSatisfy(cards []*Card) (bool, []*Card) {
	if len(cards) != PLATE_CARD_COUNT {
		return failWithSortedCards(cards)
	}

	// 使用连续性排序卡片
	sortedCards := sortCardsForConsecutive(cards)

	// 检查是否有王牌
	if hasJokers(sortedCards) {
		return false, sortedCards
	}

	// 统计变化牌数量
	wildcardCount := countWildcards(sortedCards)

	// 如果没有变化牌
	if wildcardCount == 0 {
		// 检查 3 + 3 模式
		triple1 := NewTriple(sortedCards[:3])
		triple2 := NewTriple(sortedCards[3:])
		if triple1.Valid && triple2.Valid {
			// 使用RawNumber进行连续性判断
			card1Num := triple1.Cards[0].RawNumber
			card2Num := triple2.Cards[0].RawNumber

			// 普通连续情况：如3-4
			if card1Num+1 == card2Num {
				return true, append(triple1.Cards, triple2.Cards...)
			}

			// A的特殊情况：A(1)-2，需要重新排序为A-2
			if card1Num == 1 && card2Num == 2 {
				return true, append(triple1.Cards, triple2.Cards...)
			}

			// K-A的特殊情况：K(13)-A(1)
			if card1Num == 13 && card2Num == 1 {
				return true, append(triple1.Cards, triple2.Cards...)
			}

			// A-K的特殊情况：A(1)-K(13) (连续性排序后的情况)
			if card1Num == 1 && card2Num == 13 {
				return true, append(triple1.Cards, triple2.Cards...)
			}
		}
		return failWithSortedCards(sortedCards)
	}

	// 如果有至少1个变化牌
	if wildcardCount >= 1 {
		// 检查前5张牌是否能构成葫芦
		fullhouse := NewFullHouse(sortedCards[:5])
		if fullhouse.Valid {
			// fullhouse的前3张是triple，后2张是pair
			triple := fullhouse.Cards[:3]
			pair := fullhouse.Cards[3:5]

			// 使用RawNumber进行连续性判断
			tripleNum := triple[0].RawNumber
			pairNum := pair[0].RawNumber

			// 普通连续情况
			if tripleNum+1 == pairNum {
				result := createResult()
				result = append(result, triple...)
				result = append(result, pair...)
				result = append(result, sortedCards[5]) // 变化牌
				return true, result
			}

			if tripleNum-1 == pairNum {
				result := createResult()
				result = append(result, pair...)
				result = append(result, sortedCards[5]) // 变化牌
				result = append(result, triple...)
				return true, result
			}

			// A的特殊情况：A(1)和2
			if tripleNum == 1 && pairNum == 2 {
				result := createResult()
				result = append(result, triple...)
				result = append(result, pair...)
				result = append(result, sortedCards[5]) // 变化牌
				return true, result
			}

			if tripleNum == 2 && pairNum == 1 {
				result := createResult()
				result = append(result, pair...)
				result = append(result, sortedCards[5]) // 变化牌
				result = append(result, triple...)
				return true, result
			}

			// K-A的特殊情况：K(13)和A(1)
			if tripleNum == 13 && pairNum == 1 {
				result := createResult()
				result = append(result, triple...)
				result = append(result, pair...)
				result = append(result, sortedCards[5]) // 变化牌
				return true, result
			}

			if tripleNum == 1 && pairNum == 13 {
				result := createResult()
				result = append(result, pair...)
				result = append(result, sortedCards[5]) // 变化牌
				result = append(result, triple...)
				return true, result
			}
		}
		return failWithSortedCards(sortedCards)
	}

	return failWithSortedCards(sortedCards)
}

func (p *Plate) GreaterThan(other CardComp) bool {
	if other.GetType() != TypePlate {
		return false
	}
	otherPlate := other.(*Plate)

	// 使用规范化牌组进行比较
	myCards := p.NormalizedCards
	if myCards == nil {
		myCards = p.Cards
	}
	otherCards := otherPlate.NormalizedCards
	if otherCards == nil {
		otherCards = otherPlate.Cards
	}

	// 使用getPlateComparisonKey来正确比较钢板大小
	myKey := getPlateComparisonKey(myCards)
	otherKey := getPlateComparisonKey(otherCards)

	return myKey > otherKey
}

// getPlateComparisonKey 获取钢板的比较键值
func getPlateComparisonKey(cards []*Card) int {
	if len(cards) != PLATE_CARD_COUNT {
		return 0
	}

	// 钢板结构：两个三张的连续牌型
	// cards已经通过plateSatisfy重构，前3张为第一个三张，后3张为第二个三张
	firstTriple := cards[:3]
	secondTriple := cards[3:]

	// 获取两个三张的主牌
	firstCard := firstTriple[0]
	secondCard := secondTriple[0]

	// 检查A-2钢板的特殊情况（最小的钢板）
	if isA2Plate(firstCard, secondCard) {
		return 1 // A-2钢板是最小的钢板
	}

	// 检查K-A（A-K）钢板的特殊情况（最大的钢板）
	if isAKPlate(firstCard, secondCard) {
		return 15 // K-A（A-K）钢板是最大的钢板
	}

	// 对于普通钢板，使用第一张牌的Number作为比较键值
	return firstCard.Number
}

// isA2Plate 检查是否为A-2钢板
func isA2Plate(card1, card2 *Card) bool {
	// 检查是否为A-2组合（不管顺序）
	num1 := card1.RawNumber
	num2 := card2.RawNumber

	// A-2或2-A的组合
	return (num1 == 1 && num2 == 2) || (num1 == 2 && num2 == 1)
}

// isAKPlate 检查是否为A-K钢板
func isAKPlate(card1, card2 *Card) bool {
	// 检查是否为A-K组合（不管顺序）
	num1 := card1.RawNumber
	num2 := card2.RawNumber

	// A-K或K-A的组合
	return (num1 == 1 && num2 == 13) || (num1 == 13 && num2 == 1)
}

func (p *Plate) IsBomb() bool {
	return false
}

// Tube 钢管（连续对子）
type Tube struct {
	BaseComp
}

func NewTube(cards []*Card) *Tube {
	valid := false
	var sortedCards []*Card
	var normalizedCards []*Card

	if len(cards) == 6 {
		// 使用Python的satisfy逻辑
		var ok bool
		ok, sortedCards = tubeSatisfy(cards)
		valid = ok
		
		// 如果有效，创建规范化牌组
		if valid {
			normalizedCards = normalizeTube(sortedCards)
		}
	} else {
		sortedCards = sortCards(cards)
	}

	return &Tube{
		BaseComp: BaseComp{
			Cards:           sortedCards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypeTube,
		},
	}
}

// tubeSatisfy 实现Python的Tube.satisfy逻辑
func tubeSatisfy(cards []*Card) (bool, []*Card) {
	if len(cards) != TUBE_CARD_COUNT {
		return failWithSortedCards(cards)
	}

	// 使用连续性排序卡片
	sortedCards := sortCardsForConsecutive(cards)
	wildcardCount := countWildcards(sortedCards)

	// 获取牌的数字（使用RawNumber进行连续性判断）
	cardNumbers := make([]int, len(sortedCards))
	for i, card := range sortedCards {
		if card.IsWildcard() {
			cardNumbers[i] = -1 // 变化牌标记为-1，后续处理
		} else {
			cardNumbers[i] = card.RawNumber
		}
	}

	// 检查最大牌数不超过A（使用RawNumber检查）
	if getMaxCardRawNumber(sortedCards) > 13 {
		return failWithSortedCards(cards)
	}

	// 没有变化牌的情况
	if wildcardCount == 0 {
		// 必须是 i, i, i+1, i+1, i+2, i+2
		uniqueNumbers := make(map[int]bool)
		for _, num := range cardNumbers {
			uniqueNumbers[num] = true
		}
		if len(uniqueNumbers) == 3 {
			temp := computeRelativeDiffs(cardNumbers, len(cardNumbers))
			if matchesPattern(temp, TUBE_PATTERN_TRIPLET) {
				return true, sortedCards
			}
		}
		return failWithSortedCards(cards)
	}

	// 一个变化牌的情况
	if wildcardCount == 1 {
		firstFive := computeRelativeDiffs(cardNumbers, 5)

		// i, i, i+1, i+1, i+2 wild
		if matchesPattern(firstFive, TUBE_PATTERN_0112) {
			return true, sortedCards
		}
		// i, i, i+1, i+2, i+2 wild
		if matchesPattern(firstFive, TUBE_PATTERN_0122) {
			result := createResult()
			result = append(result, sortedCards[0:3]...)
			result = append(result, sortedCards[5])
			result = append(result, sortedCards[3:5]...)
			return true, result
		}
		// i, i+1, i+1, i+2, i+2 wild
		if matchesPattern(firstFive, TUBE_PATTERN_1122) {
			result := createResult()
			result = append(result, sortedCards[0:1]...)
			result = append(result, sortedCards[5])
			result = append(result, sortedCards[1:5]...)
			return true, result
		}
		return failWithSortedCards(cards)
	}

	// 两个变化牌的情况
	if wildcardCount == 2 {
		firstFour := computeRelativeDiffs(cardNumbers, 4)

		// A-K循环连续性的特殊情况：A(1)-A(1)-K(13)-K(13) + 两个变化牌
		// 构成A-K-A钢管或类似的循环连续结构
		actualNumbers := []int{}
		for _, num := range cardNumbers {
			if num != -1 {
				actualNumbers = append(actualNumbers, num)
			}
		}
		if len(actualNumbers) == 4 {
			uniqueActual := removeDuplicates(actualNumbers)
			if len(uniqueActual) == 2 &&
				((uniqueActual[0] == 1 && uniqueActual[1] == 13) || (uniqueActual[0] == 13 && uniqueActual[1] == 1)) {
				// A-K循环钢管：AA + KK + 两个变化牌
				return true, sortedCards
			}
		}

		// i, i, i+1, i+1, wild wild
		if matchesPattern(firstFour, TUBE_PATTERN_0011) {
			if sortedCards[3].Number < 14 { // i+1 smaller than Ace
				return true, sortedCards
			} else {
				// 重新排序：将后面的变化牌移到前面
				result := createResult()
				result = append(result, sortedCards[4:6]...)
				result = append(result, sortedCards[0:4]...)
				return true, result
			}
		}
		// i, i, i+2, i+2, wild, wild
		if matchesPattern(firstFour, TUBE_PATTERN_0022) {
			result := createResult()
			result = append(result, sortedCards[0:2]...)
			result = append(result, sortedCards[4:6]...)
			result = append(result, sortedCards[2:4]...)
			return true, result
		}
		// i i i+1 i+2 wild wild
		if matchesPattern(firstFour, TUBE_PATTERN_0012) {
			result := createResult()
			result = append(result, sortedCards[0:2]...)
			result = append(result, sortedCards[2], sortedCards[5])
			result = append(result, sortedCards[3], sortedCards[4])
			return true, result
		}
		// i i+1 i+1 i+2 wild wild
		if matchesPattern(firstFour, TUBE_PATTERN_0112_ALT) {
			result := createResult()
			result = append(result, sortedCards[0], sortedCards[5])
			result = append(result, sortedCards[1:3]...)
			result = append(result, sortedCards[3], sortedCards[4])
			return true, result
		}
		// i i+1 i+2 i+2 wild wild
		if matchesPattern(firstFour, TUBE_PATTERN_0122_ALT) {
			result := createResult()
			result = append(result, sortedCards[0], sortedCards[5])
			result = append(result, sortedCards[1], sortedCards[4])
			result = append(result, sortedCards[2:4]...)
			return true, result
		}
		return failWithSortedCards(cards)
	}

	return failWithSortedCards(cards)
}

func (t *Tube) GreaterThan(other CardComp) bool {
	if other.GetType() != TypeTube {
		return false
	}
	otherTube := other.(*Tube)

	// 使用规范化牌组进行比较
	myCards := t.NormalizedCards
	if myCards == nil {
		myCards = t.Cards
	}
	otherCards := otherTube.NormalizedCards
	if otherCards == nil {
		otherCards = otherTube.Cards
	}

	// 使用getTubeComparisonKey来正确比较钢管大小
	myKey := getTubeComparisonKey(myCards)
	otherKey := getTubeComparisonKey(otherCards)

	return myKey > otherKey
}

// getTubeComparisonKey 获取钢管的比较键值
func getTubeComparisonKey(cards []*Card) int {
	if len(cards) != TUBE_CARD_COUNT {
		return 0
	}

	// 钢管结构：三个连续对子
	// cards已经通过tubeSatisfy重构，按照连续对子的结构排列

	// 提取所有非变化牌的RawNumber，用于确定钢管的基础数值
	uniqueNumbers := make(map[int]bool)
	for _, card := range cards {
		if card != nil && !card.IsWildcard() {
			uniqueNumbers[card.RawNumber] = true
		}
	}

	// 将唯一数字转换为切片并排序
	numbers := make([]int, 0, len(uniqueNumbers))
	for num := range uniqueNumbers {
		numbers = append(numbers, num)
	}

	if len(numbers) == 0 {
		// 全是变化牌的情况，使用第一张牌的Number
		return cards[0].Number
	}

	// 检查A-2钢管的特殊情况（最小的钢管）
	if isA2Tube(numbers) {
		return 1 // A-2钢管是最小的钢管
	}

	// 检查K-A（A-K）钢管的特殊情况（最大的钢管）
	if isAKTube(numbers) {
		return 15 // K-A（A-K）钢管是最大的钢管
	}

	// 对于普通钢管，使用最小的RawNumber对应的Number作为比较键值
	minRawNumber := numbers[0]
	for _, num := range numbers {
		if num < minRawNumber {
			minRawNumber = num
		}
	}

	// 将RawNumber转换为Number用于比较
	for _, card := range cards {
		if !card.IsWildcard() && card.RawNumber == minRawNumber {
			return card.Number
		}
	}

	return 0
}

// isA2Tube 检查是否为A-2钢管
func isA2Tube(numbers []int) bool {
	// 对于钢管，检查是否包含A(1)、2、3，但不包含K(13)
	// A-2-3是最小的钢管
	hasA := false
	has2 := false
	has3 := false
	hasK := false

	for _, num := range numbers {
		if num == 1 {
			hasA = true
		}
		if num == 2 {
			has2 = true
		}
		if num == 3 {
			has3 = true
		}
		if num == 13 {
			hasK = true
		}
	}

	// A-2-3钢管：必须有A、2、3，且不能有K
	return hasA && has2 && has3 && !hasK
}

// isAKTube 检查是否为A-K钢管
func isAKTube(numbers []int) bool {
	// 对于钢管，检查是否包含A(1)、K(13)、Q(12)或2
	// Q-K-A 或 A-K-2 都是最大的钢管
	hasA := false
	hasK := false
	hasQ := false
	has2 := false

	for _, num := range numbers {
		if num == 1 {
			hasA = true
		}
		if num == 13 {
			hasK = true
		}
		if num == 12 {
			hasQ = true
		}
		if num == 2 {
			has2 = true
		}
	}

	// Q-K-A 或 A-K-2 钢管
	return hasA && hasK && (hasQ || has2)
}

func (t *Tube) IsBomb() bool {
	return false
}

// JokerBomb 王炸
type JokerBomb struct {
	BaseComp
}

func NewJokerBomb(cards []*Card) *JokerBomb {
	valid := false
	sortedCards := sortCards(cards)

	if len(cards) == 4 {
		numbers := make([]int, 0)
		for _, card := range sortedCards {
			numbers = append(numbers, card.Number)
		}
		sort.Ints(numbers)

		// 检查是否为两个小王和两个大王
		valid = len(numbers) == 4 && numbers[0] == 15 && numbers[1] == 15 && numbers[2] == 16 && numbers[3] == 16
	}

	return &JokerBomb{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeJokerBomb,
		},
	}
}

func (j *JokerBomb) GreaterThan(other CardComp) bool {
	// 王炸是最大的牌组
	return other.GetType() != TypeJokerBomb
}

func (j *JokerBomb) IsBomb() bool {
	return true
}

// NaiveBomb 炸弹
type NaiveBomb struct {
	BaseComp
}

func NewNaiveBomb(cards []*Card) *NaiveBomb {
	valid := false
	sortedCards := sortCards(cards)
	var normalizedCards []*Card

	if len(cards) >= 4 {
		normalCards := getNormalCards(sortedCards)

		// 检查是否所有正常牌都是同一数字
		if len(normalCards) > 0 {
			baseNumber := normalCards[0].Number
			allSame := true
			for _, card := range normalCards {
				if card.Number != baseNumber {
					allSame = false
					break
				}
			}
			valid = allSame
		} else {
			valid = countWildcards(sortedCards) >= 4
		}
		
		// 如果有效，创建规范化牌组
		if valid {
			normalizedCards = normalizeNaiveBomb(sortedCards)
		}
	}

	return &NaiveBomb{
		BaseComp: BaseComp{
			Cards:           sortedCards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypeNaiveBomb,
		},
	}
}

func (n *NaiveBomb) GreaterThan(other CardComp) bool {
	// 如果对方不是炸弹，炸弹总是更大
	if !other.IsBomb() {
		return true
	}

	// 如果对方是王炸，炸弹总是更小
	if other.GetType() == TypeJokerBomb {
		return false
	}

	// 如果对方是同花顺
	if other.GetType() == TypeStraightFlush {
		// 6张以上的炸弹 > 同花顺
		return len(n.Cards) >= 6
	}

	// 如果对方也是炸弹，按照Python逻辑比较张数然后比较数值
	if other.GetType() == TypeNaiveBomb {
		otherBomb := other.(*NaiveBomb)
		if len(n.Cards) > len(otherBomb.Cards) {
			return true
		} else if len(n.Cards) < len(otherBomb.Cards) {
			return false
		} else {
			// 张数相同，比较数值
			// 使用规范化牌组进行比较
			myCards := n.NormalizedCards
			if myCards == nil {
				myCards = n.Cards
			}
			otherCards := otherBomb.NormalizedCards
			if otherCards == nil {
				otherCards = otherBomb.Cards
			}
			return myCards[0].GreaterThan(otherCards[0])
		}
	}

	return false
}

func (n *NaiveBomb) IsBomb() bool {
	return true
}

// StraightFlush 同花顺
type StraightFlush struct {
	BaseComp
}

func NewStraightFlush(cards []*Card) *StraightFlush {
	valid := false
	sortedCards := make([]*Card, len(cards))
	var normalizedCards []*Card

	if len(cards) == 5 {
		// 首先检查是否为顺子
		straight := NewStraight(cards)
		if straight.IsValid() {
			sortedCards = straight.Cards
			wildcardCount := countWildcards(sortedCards)
			colors := make(map[string]int)

			for _, card := range sortedCards {
				if !card.IsWildcard() {
					colors[card.Color]++
				}
			}

			// 检查是否为同花（根据变化牌数量）
			valid = (wildcardCount == 0 && len(colors) == 1) ||
				(wildcardCount == 1 && len(colors) == 1) ||
				(wildcardCount == 2 && len(colors) == 1)
				
			// 如果有效，创建规范化牌组
			if valid {
				normalizedCards = normalizeStraightFlush(sortedCards)
			}
		}
	}

	return &StraightFlush{
		BaseComp: BaseComp{
			Cards:           sortedCards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypeStraightFlush,
		},
	}
}

func (s *StraightFlush) GreaterThan(other CardComp) bool {
	// 如果对方不是炸弹，同花顺总是更大
	if !other.IsBomb() {
		return true
	}

	// 如果对方是王炸，同花顺总是更小
	if other.GetType() == TypeJokerBomb {
		return false
	}

	// 如果对方是同花顺，比较数值
	if other.GetType() == TypeStraightFlush {
		otherStraightFlush := other.(*StraightFlush)
		// 使用规范化牌组进行比较
		myCards := s.NormalizedCards
		if myCards == nil {
			myCards = s.Cards
		}
		otherCards := otherStraightFlush.NormalizedCards
		if otherCards == nil {
			otherCards = otherStraightFlush.Cards
		}
		// 使用顺子的比较方式
		myKey := getStraightComparisonKey(myCards)
		otherKey := getStraightComparisonKey(otherCards)
		return myKey > otherKey
	}
	// 如果对方是炸弹，5张以下的炸弹 < 同花顺
	if other.GetType() == TypeNaiveBomb {
		return len(other.GetCards()) <= 5
	}

	return false
}

func (s *StraightFlush) IsBomb() bool {
	return true
}

// 万能牌替换相关工具函数

// cloneCard 克隆一张牌
func cloneCard(card *Card) *Card {
	return &Card{
		Number:    card.Number,
		RawNumber: card.RawNumber,
		Color:     card.Color,
		Level:     card.Level,
		Name:      card.Name,
	}
}

// cloneCards 克隆牌组
func cloneCards(cards []*Card) []*Card {
	result := make([]*Card, len(cards))
	for i, card := range cards {
		result[i] = cloneCard(card)
	}
	return result
}

// createReplacementCard 创建一张替换牌
func createReplacementCard(rawNumber int, color string, level int) *Card {
	// 使用NewCard来正确创建牌，确保Number和RawNumber都被正确设置
	number := rawNumber
	if rawNumber == 1 {
		number = 14 // Ace conversion
	}
	card, _ := NewCard(number, color, level)
	return card
}

// replaceWildcardInPlace 在克隆的牌组中原地替换万能牌
func replaceWildcardInPlace(cards []*Card, wildcardIndex int, rawNumber int, color string) {
	if wildcardIndex >= 0 && wildcardIndex < len(cards) && cards[wildcardIndex].IsWildcard() {
		cards[wildcardIndex] = createReplacementCard(rawNumber, color, cards[wildcardIndex].Level)
	}
}

// findWildcardIndices 找出所有万能牌的索引
func findWildcardIndices(cards []*Card) []int {
	indices := []int{}
	for i, card := range cards {
		if card.IsWildcard() {
			indices = append(indices, i)
		}
	}
	return indices
}

// getMostCommonColor 获取牌组中最常见的花色（用于同花顺）
func getMostCommonColor(cards []*Card) string {
	colorCount := make(map[string]int)
	for _, card := range cards {
		if !card.IsWildcard() && card.Color != "Joker" {
			colorCount[card.Color]++
		}
	}
	
	maxCount := 0
	mostCommon := "Spade" // 默认黑桃
	for color, count := range colorCount {
		if count > maxCount {
			maxCount = count
			mostCommon = color
		}
	}
	return mostCommon
}

// normalizeFullHouse 规范化葫芦牌组
// 葫芦按照 3+2 的顺序排列，万能牌替换为使三张部分最大
func normalizeFullHouse(cards []*Card) []*Card {
	if len(cards) != 5 {
		return cards
	}
	
	result := cloneCards(cards)
	// 葫芦已经按照 3+2 的顺序排列
	// 前3张是三张，后2张是对子
	
	// 处理三张部分的万能牌
	var tripleBase *Card
	for i := 0; i < 3; i++ {
		if !result[i].IsWildcard() {
			tripleBase = result[i]
			break
		}
	}
	if tripleBase != nil {
		for i := 0; i < 3; i++ {
			if result[i].IsWildcard() {
				result[i] = cloneCard(tripleBase)
			}
		}
	}
	
	// 处理对子部分的万能牌
	var pairBase *Card
	for i := 3; i < 5; i++ {
		if !result[i].IsWildcard() {
			pairBase = result[i]
			break
		}
	}
	if pairBase != nil {
		for i := 3; i < 5; i++ {
			if result[i].IsWildcard() {
				result[i] = cloneCard(pairBase)
			}
		}
	}
	
	return result
}

// normalizePlate 规范化钢板牌组（连续三张）
// 万能牌替换为使钢板最大的牌
func normalizePlate(cards []*Card) []*Card {
	if len(cards) != 6 {
		return cards
	}
	
	result := cloneCards(cards)
	
	// 钢板是按照三张为一组排列的：[AAA, BBB]
	// 处理第一组三张
	var firstGroupBase *Card
	for i := 0; i < 3; i++ {
		if !result[i].IsWildcard() {
			firstGroupBase = result[i]
			break
		}
	}
	if firstGroupBase != nil {
		for i := 0; i < 3; i++ {
			if result[i].IsWildcard() {
				result[i] = cloneCard(firstGroupBase)
			}
		}
	}
	
	// 处理第二组三张
	var secondGroupBase *Card
	for i := 3; i < 6; i++ {
		if !result[i].IsWildcard() {
			secondGroupBase = result[i]
			break
		}
	}
	if secondGroupBase != nil {
		for i := 3; i < 6; i++ {
			if result[i].IsWildcard() {
				result[i] = cloneCard(secondGroupBase)
			}
		}
	}
	
	// 如果某组全是万能牌，需要根据另一组推断
	if firstGroupBase == nil && secondGroupBase != nil {
		// 第一组全是万能牌，应该是第二组-1
		expectedRawNumber := secondGroupBase.RawNumber - 1
		for i := 0; i < 3; i++ {
			result[i] = createReplacementCard(expectedRawNumber, "Spade", result[i].Level)
		}
	} else if secondGroupBase == nil && firstGroupBase != nil {
		// 第二组全是万能牌，应该是第一组+1
		expectedRawNumber := firstGroupBase.RawNumber + 1
		for i := 3; i < 6; i++ {
			result[i] = createReplacementCard(expectedRawNumber, "Spade", result[i].Level)
		}
	}
	
	return result
}

// normalizeStraight 规范化顺子牌组
// 万能牌替换为使顺子最大的牌
func normalizeStraight(cards []*Card) []*Card {
	if len(cards) != 5 {
		return cards
	}
	
	result := cloneCards(cards)
	
	// 找出万能牌的位置
	wildcardIndices := []int{}
	for i, card := range result {
		if card.IsWildcard() {
			wildcardIndices = append(wildcardIndices, i)
		}
	}
	
	if len(wildcardIndices) == 0 {
		return result
	}
	
	// 获取非万能牌的数字，确定顺子的范围
	numbers := []int{}
	for _, card := range result {
		if !card.IsWildcard() {
			numbers = append(numbers, card.RawNumber)
		}
	}
	
	// 顺子已经是排好序的，万能牌需要填补空缺
	// 根据顺子的实际值来确定万能牌应该是什么
	for _, idx := range wildcardIndices {
		// 根据位置确定万能牌应该的值
		if idx == 0 {
			// 第一张，应该是第二张-1
			if result[1].IsWildcard() {
				// 如果第二张也是万能牌，看第三张
				result[idx] = createReplacementCard(result[2].RawNumber-2, "Spade", result[idx].Level)
			} else {
				result[idx] = createReplacementCard(result[1].RawNumber-1, "Spade", result[idx].Level)
			}
		} else if idx == 4 {
			// 最后一张，应该是倒数第二张+1
			if result[3].IsWildcard() {
				// 如果倒数第二张也是万能牌，看倒数第三张
				result[idx] = createReplacementCard(result[2].RawNumber+2, "Spade", result[idx].Level)
			} else {
				result[idx] = createReplacementCard(result[3].RawNumber+1, "Spade", result[idx].Level)
			}
		} else {
			// 中间的牌，根据前后推断
			expectedValue := 0
			if !result[idx-1].IsWildcard() {
				expectedValue = result[idx-1].RawNumber + 1
			} else if !result[idx+1].IsWildcard() {
				expectedValue = result[idx+1].RawNumber - 1
			}
			result[idx] = createReplacementCard(expectedValue, "Spade", result[idx].Level)
		}
	}
	
	// 特殊处理A高位顺子(10-J-Q-K-A)
	// 如果最后一张是A(RawNumber=1)，说明是A高位顺子
	if result[4].RawNumber == 1 {
		// A在高位顺子中相当于14
		for _, idx := range wildcardIndices {
			if idx == 4 {
				result[idx] = createReplacementCard(1, "Spade", result[idx].Level) // A保持为1
			} else {
				// 重新计算其他位置的值：10, 11, 12, 13, 1(A)
				positionValues := []int{10, 11, 12, 13, 1}
				result[idx] = createReplacementCard(positionValues[idx], "Spade", result[idx].Level)
			}
		}
	}
	
	return result
}


// normalizeNaiveBomb 规范化普通炸弹牌组
// 万能牌全部替换为与其他牌相同的牌
func normalizeNaiveBomb(cards []*Card) []*Card {
	result := cloneCards(cards)
	
	// 找到非万能牌作为基准
	var baseCard *Card
	for _, card := range result {
		if !card.IsWildcard() {
			baseCard = card
			break
		}
	}
	
	// 如果没有非万能牌（全是万能牌），则不需要替换
	if baseCard == nil {
		return result
	}
	
	// 将所有万能牌替换为基准牌
	for i, card := range result {
		if card.IsWildcard() {
			result[i] = createReplacementCard(
				baseCard.Number,
				baseCard.Color,
				card.Level,
			)
		}
	}
	
	return result
}

// normalizeStraightFlush 规范化同花顺牌组
// 万能牌替换时需要同时满足顺子和同花色要求
func normalizeStraightFlush(cards []*Card) []*Card {
	if len(cards) != 5 {
		return cards
	}
	
	// 先按照顺子规则规范化
	result := normalizeStraight(cards)
	
	// 然后调整花色，使所有牌同花色
	mostCommonColor := getMostCommonColor(cards)
	for i, card := range result {
		// 保留原始的RawNumber，只改变花色
		result[i] = &Card{
			Number:    card.Number,
			RawNumber: card.RawNumber,
			Color:     mostCommonColor,
			Level:     card.Level,
			Name:      card.Name,
		}
	}
	
	return result
}

// normalizeTube 规范化钢管牌组（连续对子）
// 万能牌替换为使钢管最大的牌
func normalizeTube(cards []*Card) []*Card {
	if len(cards) != 6 {
		return cards
	}
	
	result := cloneCards(cards)
	
	// 钢管是三个连续的对子，例如：[3,3,4,4,5,5]
	// 找出所有非万能牌的数字，确定钢管的基础值
	baseNumbers := []int{}
	wildcardCount := 0
	wildcardIndices := []int{}
	
	for i, card := range result {
		if card.IsWildcard() {
			wildcardCount++
			wildcardIndices = append(wildcardIndices, i)
		} else {
			baseNumbers = append(baseNumbers, card.RawNumber)
		}
	}
	
	if wildcardCount == 0 {
		return result
	}
	
	// 统计每个数字出现的次数
	numberCount := make(map[int]int)
	for _, num := range baseNumbers {
		numberCount[num]++
	}
	
	// 找出已有的对子和单张
	pairs := []int{}
	singles := []int{}
	for num, count := range numberCount {
		if count >= 2 {
			pairs = append(pairs, num)
		} else if count == 1 {
			singles = append(singles, num)
		}
	}
	
	// 排序以便处理
	sort.Ints(pairs)
	sort.Ints(singles)
	
	// 根据万能牌数量和现有牌的分布，确定如何替换万能牌
	if wildcardCount == 1 {
		// 一个万能牌：应该与某个单张组成对子
		if len(singles) > 0 {
			// 替换为最大的单张，使钢管尽可能大
			targetNum := singles[len(singles)-1]
			for _, idx := range wildcardIndices {
				result[idx] = createReplacementCard(targetNum, "Spade", result[idx].Level)
			}
		}
	} else if wildcardCount == 2 {
		// 两个万能牌
		if len(singles) == 2 {
			// 两个单张，每个万能牌与一个单张配对
			// 替换为较大的单张，使钢管更大
			targetNum := singles[1] // 第二大的单张
			for _, idx := range wildcardIndices {
				result[idx] = createReplacementCard(targetNum, "Spade", result[idx].Level)
			}
		} else if len(pairs) == 2 {
			// 已有两个对子，万能牌应该组成第三个对子
			// 确定缺失的数字
			minPair := pairs[0]
			maxPair := pairs[1]
			var targetNum int
			
			if maxPair - minPair == 2 {
				// 缺中间的对子
				targetNum = minPair + 1
			} else if maxPair - minPair == 1 {
				// 缺最大或最小的对子，选择最大的
				targetNum = maxPair + 1
				if targetNum > 13 {
					targetNum = minPair - 1
				}
			} else if minPair == 1 && maxPair == 13 {
				// A-K 特殊情况，缺中间的2
				targetNum = 2
			} else {
				// 其他情况，默认使用较大值+1
				targetNum = maxPair + 1
				if targetNum > 13 {
					targetNum = 1 // 循环到A
				}
			}
			
			for _, idx := range wildcardIndices {
				result[idx] = createReplacementCard(targetNum, "Spade", result[idx].Level)
			}
		}
	}
	
	// 处理更多万能牌的情况
	if wildcardCount > 2 {
		// 确定钢管的最大可能值
		// 从最大的非万能牌开始，构建连续对子
		maxNum := 0
		for _, card := range result {
			if !card.IsWildcard() && card.RawNumber > maxNum {
				maxNum = card.RawNumber
			}
		}
		
		// 根据已有的牌构建最优的钢管
		wildcardIdx := 0
		for i := range result {
			if result[i].IsWildcard() && wildcardIdx < len(wildcardIndices) {
				// 根据位置确定应该替换的数字
				// 这里需要更复杂的逻辑来确定最优替换
				result[i] = createReplacementCard(maxNum, "Spade", result[i].Level)
				wildcardIdx++
			}
		}
	}
	
	return result
}
