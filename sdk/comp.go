package sdk

import (
	"fmt"
	"sort"
)

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
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].LessThan(sortedCards[j])
	})

	// 统计变化牌数量
	numWildcards := 0
	for _, card := range sortedCards {
		if card.IsWildcard() {
			numWildcards++
		}
	}

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
	normalCards := []*Card{}

	for _, card := range sortedCards {
		if card.IsWildcard() {
			wildcards = append(wildcards, card)
		} else {
			normalCards = append(normalCards, card)
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

		// 如果前一个牌组不是钢板，优先尝试钢管
		if prev == nil || prev.GetType() != TypePlate {
			if comp := NewTube(cards); comp.IsValid() {
				return comp
			}
		}

		if comp := NewPlate(cards); comp.IsValid() {
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

	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].LessThan(sortedCards[j])
	})

	return &Pair{
		BaseComp: BaseComp{
			Cards: sortedCards,
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
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].LessThan(sortedCards[j])
	})

	if len(cards) == 3 {
		// 如果有王，则非法
		if sortedCards[len(sortedCards)-1].Color == "Joker" {
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
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].LessThan(sortedCards[j])
	})

	if len(cards) == 5 {
		// 复杂的葫芦判断逻辑
		// 简化实现：检查是否有3张相同和2张相同
		cardCounts := make(map[int]int)
		for _, card := range sortedCards {
			if card.IsWildcard() {
				continue // 变化牌稍后处理
			}
			cardCounts[card.Number]++
		}

		// 统计变化牌数量
		wildcardCount := 0
		for _, card := range sortedCards {
			if card.IsWildcard() {
				wildcardCount++
			}
		}

		// 检查是否能构成葫芦
		counts := make([]int, 0)
		for _, count := range cardCounts {
			counts = append(counts, count)
		}
		sort.Ints(counts)

		// 根据变化牌数量判断是否能构成葫芦
		if wildcardCount == 0 {
			valid = len(counts) == 2 && counts[0] == 2 && counts[1] == 3
		} else if wildcardCount == 1 {
			valid = (len(counts) == 2 && ((counts[0] == 1 && counts[1] == 3) || (counts[0] == 2 && counts[1] == 2))) ||
				(len(counts) == 1 && counts[0] == 4)
		} else if wildcardCount == 2 {
			valid = (len(counts) == 1 && counts[0] == 3) ||
				(len(counts) == 2 && counts[0] == 1 && counts[1] == 2)
		}
	}

	return &FullHouse{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeFullHouse,
		},
	}
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
	sortedCards := SortNoLevel(cards)

	if len(cards) == 5 {
		// 统计变化牌数量
		wildcardCount := 0
		for _, card := range sortedCards {
			if card.IsWildcard() {
				wildcardCount++
			}
		}

		// 获取非变化牌的数字
		normalCards := []*Card{}
		for _, card := range sortedCards {
			if !card.IsWildcard() {
				normalCards = append(normalCards, card)
			}
		}

		// 检查最大牌不超过A
		maxNumber := 0
		for _, card := range cards {
			if card.Number > maxNumber {
				maxNumber = card.Number
			}
		}

		if maxNumber <= 14 {
			valid = checkStraightValid(normalCards, wildcardCount)
		}
	}

	return &Straight{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeStraight,
		},
	}
}

func checkStraightValid(normalCards []*Card, wildcardCount int) bool {
	if len(normalCards) == 0 {
		return wildcardCount == 5
	}

	// 按RawNumber排序用于顺子检查
	sort.Slice(normalCards, func(i, j int) bool {
		return normalCards[i].RawNumber < normalCards[j].RawNumber
	})

	// 检查是否能构成顺子
	needed := 5 - len(normalCards)
	if needed != wildcardCount {
		return false
	}

	// 检查正常牌之间的间隔
	gaps := 0
	for i := 1; i < len(normalCards); i++ {
		gap := normalCards[i].RawNumber - normalCards[i-1].RawNumber - 1
		if gap < 0 {
			return false // 有重复
		}
		gaps += gap
	}

	// 检查首尾需要的牌数
	if len(normalCards) > 1 {
		totalSpan := normalCards[len(normalCards)-1].RawNumber - normalCards[0].RawNumber + 1
		if totalSpan > 5 {
			return false
		}
		return gaps <= wildcardCount
	}

	return true
}

func (s *Straight) GreaterThan(other CardComp) bool {
	if other.GetType() != TypeStraight {
		return false
	}
	otherStraight := other.(*Straight)
	return s.Cards[0].ConsecutiveGreaterThan(otherStraight.Cards[0])
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
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].LessThan(sortedCards[j])
	})

	if len(cards) == 6 {
		// 检查是否没有王
		hasJoker := false
		for _, card := range sortedCards {
			if card.Color == "Joker" {
				hasJoker = true
				break
			}
		}

		if !hasJoker {
			// 简化的钢板判断
			valid = checkPlateValid(sortedCards)
		}
	}

	return &Plate{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypePlate,
		},
	}
}

func checkPlateValid(cards []*Card) bool {
	// 检查是否为两个连续的三张
	// 简化实现
	cardCounts := make(map[int]int)
	for _, card := range cards {
		if card.IsWildcard() {
			continue
		}
		cardCounts[card.Number]++
	}

	// 统计变化牌
	wildcardCount := 0
	for _, card := range cards {
		if card.IsWildcard() {
			wildcardCount++
		}
	}

	numbers := make([]int, 0)
	for num := range cardCounts {
		numbers = append(numbers, num)
	}
	sort.Ints(numbers)

	// 检查是否为连续的两个三张
	if len(numbers) == 2 && numbers[1] == numbers[0]+1 {
		count1 := cardCounts[numbers[0]]
		count2 := cardCounts[numbers[1]]
		return count1+count2+wildcardCount == 6 && count1 >= 1 && count2 >= 1
	}

	return false
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
	sortedCards := SortNoLevel(cards)

	if len(cards) == 6 {
		valid = checkTubeValid(sortedCards)
	}

	return &Tube{
		BaseComp: BaseComp{
			Cards: sortedCards,
			Valid: valid,
			Type:  TypeTube,
		},
	}
}

func checkTubeValid(cards []*Card) bool {
	// 检查是否为三个连续的对子
	// 简化实现
	cardCounts := make(map[int]int)
	for _, card := range cards {
		if card.IsWildcard() {
			continue
		}
		cardCounts[card.Number]++
	}

	wildcardCount := 0
	for _, card := range cards {
		if card.IsWildcard() {
			wildcardCount++
		}
	}

	numbers := make([]int, 0)
	for num := range cardCounts {
		numbers = append(numbers, num)
	}
	sort.Ints(numbers)

	// 检查是否为连续的三个对子
	if len(numbers) == 3 {
		if numbers[1] == numbers[0]+1 && numbers[2] == numbers[1]+1 {
			totalCount := 0
			for _, count := range cardCounts {
				totalCount += count
			}
			return totalCount+wildcardCount == 6
		}
	}

	return false
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
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].LessThan(sortedCards[j])
	})

	if len(cards) == 4 {
		numbers := make([]int, 0)
		for _, card := range sortedCards {
			numbers = append(numbers, card.Number)
		}
		sort.Ints(numbers)

		// 检查是否为两个小王和两个大王
		if len(numbers) == 4 && numbers[0] == 15 && numbers[1] == 15 && numbers[2] == 16 && numbers[3] == 16 {
			valid = true
		}
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
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].LessThan(sortedCards[j])
	})

	if len(cards) >= 4 {
		// 统计变化牌数量
		wildcardCount := 0
		for _, card := range sortedCards {
			if card.IsWildcard() {
				wildcardCount++
			}
		}

		// 获取非变化牌的数字
		normalCards := []*Card{}
		for _, card := range sortedCards {
			if !card.IsWildcard() {
				normalCards = append(normalCards, card)
			}
		}

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
			valid = wildcardCount >= 4
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

	// 如果对方也是炸弹，比较张数然后比较数值
	if other.GetType() == TypeNaiveBomb {
		otherBomb := other.(*NaiveBomb)
		if len(n.Cards) != len(otherBomb.Cards) {
			return len(n.Cards) > len(otherBomb.Cards)
		}
		return n.Cards[0].GreaterThan(otherBomb.Cards[0])
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
	copy(sortedCards, cards)

	if len(cards) == 5 {
		// 首先检查是否为顺子
		straight := NewStraight(cards)
		if straight.IsValid() {
			sortedCards = straight.Cards

			// 然后检查是否为同花
			wildcardCount := 0
			colors := make(map[string]int)

			for _, card := range sortedCards {
				if card.IsWildcard() {
					wildcardCount++
				} else {
					colors[card.Color]++
				}
			}

			// 检查是否为同花
			if wildcardCount == 0 {
				valid = len(colors) == 1
			} else if wildcardCount == 1 {
				valid = len(colors) == 1
			} else if wildcardCount == 2 {
				valid = len(colors) == 1
			}
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
