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

// BaseComp 基础牌组结构
type BaseComp struct {
	Cards []*Card
	Valid bool
	Type  CompType
}

// GetCards 获取牌组中的牌
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

// SortNoLevel 排序时将级别牌放在合适位置
func SortNoLevel(cards []*Card) []*Card {
	if len(cards) == 0 {
		return cards
	}

	// 首先按照正常规则排序
	sortedCards := sortCards(cards)

	// 统计变化牌数量
	numWildcards := countWildcards(sortedCards)

	if numWildcards == 0 {
		// 检查最后一张牌是否为级别牌
		lastCard := sortedCards[len(sortedCards)-1]
		if lastCard.Level != lastCard.Number {
			return sortedCards
		}
		// 如果最后一张是级别牌，则需要重新排序
	}

	// 将非变化牌重新排序
	wildcards := []*Card{}
	normalCards := getNormalCards(sortedCards)

	for _, card := range sortedCards {
		if card.IsWildcard() {
			wildcards = append(wildcards, card)
		}
	}

	// 按数字值排序正常牌
	sort.Slice(normalCards, func(i, j int) bool {
		return normalCards[i].Number < normalCards[j].Number
	})

	// 合并结果
	result := make([]*Card, 0, len(cards))
	result = append(result, normalCards...)
	result = append(result, wildcards...)

	return result
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
		if comp := NewFullHouse(cards); comp.IsValid() {
			return comp
		}
		if comp := NewStraight(cards); comp.IsValid() {
			return comp
		}
		return &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}

	case 6:
		// 炸弹 > 钢板/钢管
		if comp := NewNaiveBomb(cards); comp.IsValid() {
			return comp
		}

		// 优先级处理：根据Python实现，如果前一个牌组不是钢板，优先尝试钢管
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
	if len(cards) == 2 {
		levelCond0 := cards[0].IsWildcard() && cards[1].Color != "Joker"
		levelCond1 := cards[1].IsWildcard() && cards[0].Color != "Joker"
		valid = cards[0].Equals(cards[1]) || levelCond0 || levelCond1
	}

	return &Pair{
		BaseComp: BaseComp{
			Cards: sortCards(cards),
			Valid: valid,
			Type:  TypePair,
		},
	}
}

func (p *Pair) GreaterThan(other CardComp) bool {
	if other.GetType() != TypePair {
		return false
	}
	otherPair := other.(*Pair)
	return p.Cards[0].GreaterThan(otherPair.Cards[0])
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

			// 如果有变化牌，调整其数值
			if valid {
				for i := 1; i < len(sortedCards); i++ {
					if sortedCards[i].IsWildcard() {
						sortedCards[i].Number = baseCard.Number
					}
				}
			}
		}
	}

	return &Triple{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeTriple,
		},
	}
}

func (t *Triple) GreaterThan(other CardComp) bool {
	if other.GetType() != TypeTriple {
		return false
	}
	otherTriple := other.(*Triple)
	return t.Cards[0].GreaterThan(otherTriple.Cards[0])
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

	if len(cards) == 5 {
		// 使用Python的satisfy逻辑
		var ok bool
		ok, sortedCards = fullHouseSatisfy(cards)
		valid = ok
	} else {
		sortedCards = sortCards(cards)
	}

	return &FullHouse{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeFullHouse,
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
	return f.Cards[0].GreaterThan(otherFullHouse.Cards[0])
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

	if len(cards) == 5 {
		// 使用Python的satisfy逻辑
		var ok bool
		ok, sortedCards = straightSatisfy(cards)
		valid = ok
	} else {
		sortedCards = sortCards(cards)
	}

	return &Straight{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeStraight,
		},
	}
}

// straightSatisfy 实现Python的Straight.satisfy逻辑
func straightSatisfy(cards []*Card) (bool, []*Card) {
	if len(cards) != STRAIGHT_CARD_COUNT {
		return failWithSortedCards(cards)
	}

	// 排序卡片，将level卡片放在适当位置
	sortedCards := SortNoLevel(cards)

	// 统计变化牌数量
	numWildcards := countWildcards(sortedCards)

	// 获取卡片数字
	cardNumbers := make([]int, len(sortedCards))
	for i, card := range sortedCards {
		cardNumbers[i] = card.Number
	}

	// 最大牌不能超过A
	if getMaxCardNumber(sortedCards) > 14 {
		return failWithSortedCards(cards)
	}

	// 没有变化牌
	if numWildcards == 0 {
		if cardNumbers[0]+4 == cardNumbers[4] && len(removeDuplicates(cardNumbers)) == 5 {
			return true, sortedCards
		}
		return failWithSortedCards(cards)
	}

	// 一个变化牌
	if numWildcards == 1 {
		firstFour := computeRelativeDiffs(cardNumbers, 4)

		// i, i+1, i+2, i+3 wild
		if matchesPattern(firstFour, STRAIGHT_PATTERN_0123) {
			if cardNumbers[3] <= 13 {
				return true, sortedCards
			}
			if cardNumbers[3] == 14 {
				// A的特殊处理：将变化牌放在前面
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				newOrder[0] = sortedCards[4] // 变化牌
				copy(newOrder[1:], sortedCards[0:3])
				newOrder[4] = sortedCards[3] // A
				return true, newOrder
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
		firstThree := computeRelativeDiffs(cardNumbers, 3)

		// i, i+1, i+2, wild, wild
		if matchesPattern(firstThree, STRAIGHT_PATTERN_012) {
			if cardNumbers[2] <= 12 {
				return true, sortedCards
			}
			if cardNumbers[2] == 13 {
				// K的特殊处理
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				newOrder[0] = sortedCards[4]
				copy(newOrder[1:4], sortedCards[1:4])
				newOrder[4] = sortedCards[3]
				return true, newOrder
			}
			if cardNumbers[2] == 14 {
				// A的特殊处理
				newOrder := createNewOrder(STRAIGHT_CARD_COUNT)
				copy(newOrder[0:2], sortedCards[3:5])
				copy(newOrder[2:], sortedCards[1:4])
				return true, newOrder
			}
		}

		// 处理其他二变化牌的情况
		if anyPatternMatches(firstThree, STRAIGHT_PATTERN_023, STRAIGHT_PATTERN_013, STRAIGHT_PATTERN_024, STRAIGHT_PATTERN_034, STRAIGHT_PATTERN_014) {
			return true, sortedCards
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

	// 使用getStraightComparisonKey来正确比较顺子大小
	myKey := getStraightComparisonKey(s.Cards)
	otherKey := getStraightComparisonKey(otherStraight.Cards)

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

	if len(cards) == 6 {
		// 使用Python的satisfy逻辑
		var ok bool
		ok, sortedCards = plateSatisfy(cards)
		valid = ok
	} else {
		sortedCards = sortCards(cards)
	}

	return &Plate{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypePlate,
		},
	}
}

// plateSatisfy 实现Python的Plate.satisfy逻辑
func plateSatisfy(cards []*Card) (bool, []*Card) {
	if len(cards) != PLATE_CARD_COUNT {
		return failWithSortedCards(cards)
	}

	// 排序卡片
	sortedCards := sortCards(cards)

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
			card1Num := triple1.Cards[0].Number
			card2Num := triple2.Cards[0].Number

			// 普通连续情况：如3-4
			if card1Num+1 == card2Num {
				return true, append(triple1.Cards, triple2.Cards...)
			}

			// A的特殊情况：2-A，需要重新排序为A-2
			if card1Num == 2 && card2Num == 14 {
				return true, append(triple2.Cards, triple1.Cards...)
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

			tripleNum := triple[0].Number
			pairNum := pair[0].Number

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

			// A的特殊情况
			if tripleNum == 14 && pairNum == 2 {
				result := createResult()
				result = append(result, triple...)
				result = append(result, pair...)
				result = append(result, sortedCards[5]) // 变化牌
				return true, result
			}

			if tripleNum == 2 && pairNum == 14 {
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
	return p.Cards[0].ConsecutiveGreaterThan(otherPlate.Cards[0])
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

	if len(cards) == 6 {
		// 使用Python的satisfy逻辑
		var ok bool
		ok, sortedCards = tubeSatisfy(cards)
		valid = ok
	} else {
		sortedCards = sortCards(cards)
	}

	return &Tube{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeTube,
		},
	}
}

// tubeSatisfy 实现Python的Tube.satisfy逻辑
func tubeSatisfy(cards []*Card) (bool, []*Card) {
	if len(cards) != TUBE_CARD_COUNT {
		return failWithSortedCards(cards)
	}

	// 使用sort_no_level排序
	sortedCards := SortNoLevel(cards)
	wildcardCount := countWildcards(sortedCards)

	// 获取牌的数字
	cardNumbers := make([]int, len(sortedCards))
	for i, card := range sortedCards {
		cardNumbers[i] = card.Number
	}

	// 检查最大牌数不超过A
	if getMaxCardNumber(sortedCards) > 14 {
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
	return t.Cards[0].ConsecutiveGreaterThan(otherTube.Cards[0])
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
	}

	return &NaiveBomb{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeNaiveBomb,
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
			return n.Cards[0].GreaterThan(otherBomb.Cards[0])
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
		}
	}

	return &StraightFlush{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeStraightFlush,
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
		return s.Cards[0].GreaterThan(otherStraightFlush.Cards[0])
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
