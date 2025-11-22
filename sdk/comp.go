package sdk

import (
	"fmt"
	"sort"
)

// 辅助函数
func failWithSortedCards(cards []*Card) (bool, []*Card) {
	return false, sortCards(cards)
}

func createResult() []*Card {
	return make([]*Card, 0)
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
	Cards           []*Card  `json:"cards"`            // 原始牌组（包含万能牌）
	NormalizedCards []*Card  `json:"-"`                // 规范化牌组（万能牌已替换为具体牌）- 不序列化
	Valid           bool     `json:"valid"`
	Type            CompType `json:"type"`
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
							card.DeckIndex,  // 复用原万能牌的索引
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
								card.DeckIndex,  // 复用原万能牌的索引
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
	ComparisonKey int // 比较键值（1-10），在构造时预计算
}

func NewStraight(cards []*Card) *Straight {
	valid := false
	var sortedCards []*Card
	var normalizedCards []*Card
	var comparisonKey int

	if len(cards) == 5 {
		// 使用新的枚举验证逻辑
		var ok bool
		ok, sortedCards, comparisonKey = straightSatisfyNew(cards)
		valid = ok
		
		// 如果有效，新逻辑已返回规范化结果
		if valid {
			normalizedCards = sortedCards
		}
	} else {
		sortedCards = sortCards(cards)
	}

	return &Straight{
		BaseComp: BaseComp{
			Cards:           cards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypeStraight,
		},
		ComparisonKey: comparisonKey,
	}
}

func (s *Straight) GreaterThan(other CardComp) bool {
	if other.GetType() != TypeStraight {
		return false
	}
	otherStraight := other.(*Straight)
	
	// 直接比较预计算的 ComparisonKey
	return s.ComparisonKey > otherStraight.ComparisonKey
}

func (s *Straight) IsBomb() bool {
	return false
}

// Plate 钢板（连续三张）
type Plate struct {
	BaseComp
	ComparisonKey int // 比较键值（1-13），在构造时预计算
}

func NewPlate(cards []*Card) *Plate {
	valid := false
	var sortedCards []*Card
	var normalizedCards []*Card
	var comparisonKey int

	if len(cards) == 6 {
		// 使用新的枚举验证逻辑
		var ok bool
		ok, sortedCards, comparisonKey = plateSatisfyNew(cards)
		valid = ok
		
		// 如果有效，新逻辑已返回规范化结果
		if valid {
			normalizedCards = sortedCards
		}
	} else {
		sortedCards = sortCards(cards)
	}

	return &Plate{
		BaseComp: BaseComp{
			Cards:           cards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypePlate,
		},
		ComparisonKey: comparisonKey,
	}
}

func (p *Plate) GreaterThan(other CardComp) bool {
	if other.GetType() != TypePlate {
		return false
	}
	otherPlate := other.(*Plate)
	
	// 直接比较预计算的 ComparisonKey
	return p.ComparisonKey > otherPlate.ComparisonKey
}

func (p *Plate) IsBomb() bool {
	return false
}

// Tube 钢管（连续对子）
type Tube struct {
	BaseComp
	ComparisonKey int // 比较键值（1-12），在构造时预计算
}

func NewTube(cards []*Card) *Tube {
	valid := false
	var sortedCards []*Card
	var normalizedCards []*Card
	var comparisonKey int

	if len(cards) == 6 {
		// 使用新的枚举验证逻辑
		var ok bool
		ok, sortedCards, comparisonKey = tubeSatisfyNew(cards)
		valid = ok
		
		// 如果有效，新逻辑已返回规范化结果
		if valid {
			normalizedCards = sortedCards
		}
	} else {
		sortedCards = sortCards(cards)
	}

	return &Tube{
		BaseComp: BaseComp{
			Cards:           cards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypeTube,
		},
		ComparisonKey: comparisonKey,
	}
}

func (t *Tube) GreaterThan(other CardComp) bool {
	if other.GetType() != TypeTube {
		return false
	}
	otherTube := other.(*Tube)
	
	// 直接比较预计算的 ComparisonKey
	return t.ComparisonKey > otherTube.ComparisonKey
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
	ComparisonKey int // 比较键值（1-10），在构造时预计算
}

func NewStraightFlush(cards []*Card) *StraightFlush {
	valid := false
	var sortedCards []*Card
	var normalizedCards []*Card
	var comparisonKey int

	if len(cards) == 5 {
		// 先用 straightSatisfyNew 检查顺子并获取 comparisonKey
		var isValidStraight bool
		isValidStraight, sortedCards, comparisonKey = straightSatisfyNew(cards)
		
		if isValidStraight {
			// 再检查花色
			colors := make(map[string]int)
			
			for _, card := range sortedCards {
				if !card.IsWildcard() {
					colors[card.Color]++
				}
			}
			
			// 同花条件：所有非万能牌同花色
			valid = (len(colors) == 1)
			
			if valid {
				normalizedCards = sortedCards
			}
		}
	} else {
		sortedCards = sortCards(cards)
	}

	return &StraightFlush{
		BaseComp: BaseComp{
			Cards:           cards,
			NormalizedCards: normalizedCards,
			Valid:           valid,
			Type:            TypeStraightFlush,
		},
		ComparisonKey: comparisonKey,
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

	// 如果对方是同花顺，直接比较预计算的 ComparisonKey
	if other.GetType() == TypeStraightFlush {
		otherStraightFlush := other.(*StraightFlush)
		return s.ComparisonKey > otherStraightFlush.ComparisonKey
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
		DeckIndex: card.DeckIndex,
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
// deckIndex: 使用原万能牌的 DeckIndex 保持唯一性
func createReplacementCard(rawNumber int, color string, level int, deckIndex int) *Card {
	// 使用NewCard来正确创建牌，确保Number和RawNumber都被正确设置
	number := rawNumber
	if rawNumber == 1 {
		number = 14 // Ace conversion
	}
	card, _ := NewCard(number, color, level)
	card.DeckIndex = deckIndex  // 复用原万能牌的索引
	return card
}

// replaceWildcardInPlace 在克隆的牌组中原地替换万能牌
func replaceWildcardInPlace(cards []*Card, wildcardIndex int, rawNumber int, color string) {
	if wildcardIndex >= 0 && wildcardIndex < len(cards) && cards[wildcardIndex].IsWildcard() {
		originalDeckIndex := cards[wildcardIndex].DeckIndex  // 保存原万能牌的索引
		cards[wildcardIndex] = createReplacementCard(rawNumber, color, cards[wildcardIndex].Level, originalDeckIndex)
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
		originalDeckIndex := result[idx].DeckIndex  // 保存原索引
		// 根据位置确定万能牌应该的值
		if idx == 0 {
			// 第一张，应该是第二张-1
			if result[1].IsWildcard() {
				// 如果第二张也是万能牌，看第三张
				result[idx] = createReplacementCard(result[2].RawNumber-2, "Spade", result[idx].Level, originalDeckIndex)
			} else {
				result[idx] = createReplacementCard(result[1].RawNumber-1, "Spade", result[idx].Level, originalDeckIndex)
			}
		} else if idx == 4 {
			// 最后一张，应该是倒数第二张+1
			if result[3].IsWildcard() {
				// 如果倒数第二张也是万能牌，看倒数第三张
				result[idx] = createReplacementCard(result[2].RawNumber+2, "Spade", result[idx].Level, originalDeckIndex)
			} else {
				result[idx] = createReplacementCard(result[3].RawNumber+1, "Spade", result[idx].Level, originalDeckIndex)
			}
		} else {
			// 中间的牌，根据前后推断
			expectedValue := 0
			if !result[idx-1].IsWildcard() {
				expectedValue = result[idx-1].RawNumber + 1
			} else if !result[idx+1].IsWildcard() {
				expectedValue = result[idx+1].RawNumber - 1
			}
			result[idx] = createReplacementCard(expectedValue, "Spade", result[idx].Level, originalDeckIndex)
		}
	}
	
	// 特殊处理A高位顺子(10-J-Q-K-A)
	// 如果最后一张是A(RawNumber=1)，说明是A高位顺子
	if result[4].RawNumber == 1 {
		// A在高位顺子中相当于14
		for _, idx := range wildcardIndices {
			originalDeckIndex := result[idx].DeckIndex  // 保存原索引
			if idx == 4 {
				result[idx] = createReplacementCard(1, "Spade", result[idx].Level, originalDeckIndex) // A保持为1
			} else {
				// 重新计算其他位置的值：10, 11, 12, 13, 1(A)
				positionValues := []int{10, 11, 12, 13, 1}
				result[idx] = createReplacementCard(positionValues[idx], "Spade", result[idx].Level, originalDeckIndex)
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
				card.DeckIndex,  // 复用原万能牌的索引
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


