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


// sortNormalFirst 排序：普通牌在前，wildcard 在后
// 普通牌按 LessThan 排序，wildcard 按 LessThan 排序
func sortNormalFirst(cards []*Card) []*Card {
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		// 普通牌在前，wildcard 在后
		if !sortedCards[i].IsWildcard() && sortedCards[j].IsWildcard() {
			return true
		}
		if sortedCards[i].IsWildcard() && !sortedCards[j].IsWildcard() {
			return false
		}
		// 同类内部按 LessThan 排序
		return sortedCards[i].LessThan(sortedCards[j])
	})
	return sortedCards
}

// buildSameNumberComp 验证并规范化"同数字"牌型
// 规则：
// 1. 非 wildcard 必须同一 number
// 2. wildcard 不能代替 Joker
// 返回: (valid, normalizedCards)
func buildSameNumberComp(cards []*Card, minLen, maxLen int) (bool, []*Card) {
	if len(cards) < minLen || len(cards) > maxLen {
		return false, nil
	}

	normalCards := getNormalCards(cards)
	wildcardCount := countWildcards(cards)

	// 检查非 wildcard 是否同一 number
	if len(normalCards) > 0 {
		baseNumber := normalCards[0].Number
		for _, card := range normalCards {
			if card.Number != baseNumber {
				return false, nil
			}
		}

		// wildcard 不能代替 Joker
		if normalCards[0].Color == "Joker" && wildcardCount > 0 {
			return false, nil
		}
	}

	return true, sortNormalFirst(cards)
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

// IsBomb 判断是否为炸弹类型（王炸、炸弹、同花顺）
func (b *BaseComp) IsBomb() bool {
	return b.Type == TypeJokerBomb ||
		b.Type == TypeNaiveBomb ||
		b.Type == TypeStraightFlush
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

// IllegalComp 非法牌组
type IllegalComp struct {
	BaseComp
}

func (i *IllegalComp) GreaterThan(other CardComp) bool {
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

// Pair 对子
type Pair struct {
	BaseComp
}

func NewPair(cards []*Card) *Pair {
	sortedCards := sortCards(cards)
	valid, normalizedCards := buildSameNumberComp(sortedCards, 2, 2)

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

// Triple 三张
type Triple struct {
	BaseComp
}

func NewTriple(cards []*Card) *Triple {
	sortedCards := sortCards(cards)
	valid, normalizedCards := buildSameNumberComp(sortedCards, 3, 3)

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

// FullHouse 葫芦（三带二）
type FullHouse struct {
	BaseComp
}

func NewFullHouse(cards []*Card) *FullHouse {
	valid := false
	var sortedCards []*Card
	var normalizedCards []*Card

	if len(cards) == 5 {
		var ok bool
		ok, sortedCards = FullHouseSatisfyNew(cards)
		valid = ok
		
		if valid {
			normalizedCards = sortedCards
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

// NaiveBomb 炸弹
type NaiveBomb struct {
	BaseComp
}

func NewNaiveBomb(cards []*Card) *NaiveBomb {
	sortedCards := sortCards(cards)
	valid, normalizedCards := buildSameNumberComp(sortedCards, 4, 8)

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



