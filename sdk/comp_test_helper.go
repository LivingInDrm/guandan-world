package sdk

// CreateCompByType 根据指定的类型创建牌组
// 用于测试中确保创建正确类型的牌组
func CreateCompByType(cards []*Card, compType string) CardComp {
	switch compType {
	case "Single":
		return NewSingle(cards)
	case "Pair":
		return NewPair(cards)
	case "Triple":
		return NewTriple(cards)
	case "FullHouse":
		return NewFullHouse(cards)
	case "Straight":
		return NewStraight(cards)
	case "Plate":
		return NewPlate(cards)
	case "Tube":
		return NewTube(cards)
	case "NaiveBomb":
		return NewNaiveBomb(cards)
	case "StraightFlush":
		return NewStraightFlush(cards)
	case "JokerBomb":
		return NewJokerBomb(cards)
	case "Fold":
		return &Fold{BaseComp: BaseComp{Cards: cards, Valid: true, Type: TypeFold}}
	case "IllegalComp":
		return &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}
	default:
		// 如果类型未知，使用 FromCardList
		return FromCardList(cards, nil)
	}
}

// NormalizeComp 对牌组进行万能牌替换
// 返回一个新的牌组，其中万能牌已被替换为具体的牌
func NormalizeComp(comp CardComp) CardComp {
	// 如果牌组为nil，返回nil
	if comp == nil {
		return nil
	}
	
	// 如果牌组无效，直接返回原牌组
	if !comp.IsValid() {
		return comp
	}

	// 根据不同类型进行处理
	switch c := comp.(type) {
	case *Pair:
		if c.NormalizedCards != nil {
			// 使用构造函数创建，确保正确的初始化
			return NewPair(c.NormalizedCards)
		}
	case *Triple:
		if c.NormalizedCards != nil {
			return NewTriple(c.NormalizedCards)
		}
	case *FullHouse:
		if c.NormalizedCards != nil {
			return NewFullHouse(c.NormalizedCards)
		}
	case *Straight:
		if c.NormalizedCards != nil {
			return NewStraight(c.NormalizedCards)
		}
	case *Plate:
		if c.NormalizedCards != nil {
			return NewPlate(c.NormalizedCards)
		}
	case *Tube:
		if c.NormalizedCards != nil {
			return NewTube(c.NormalizedCards)
		}
	case *NaiveBomb:
		if c.NormalizedCards != nil {
			return NewNaiveBomb(c.NormalizedCards)
		}
	case *StraightFlush:
		if c.NormalizedCards != nil {
			// 使用 NewStraightFlush 创建，确保正确的初始化
			return NewStraightFlush(c.NormalizedCards)
		}
	}

	// 对于其他类型或没有NormalizedCards的情况，返回原牌组
	return comp
}

// FormatCardsSimple 将卡牌列表转换为简单字符串格式
// 格式: "12S,12C,13H,13D,1S,1C"
// RawNumber + Color首字母，逗号分隔
func FormatCardsSimple(cards []*Card) string {
	if len(cards) == 0 {
		return ""
	}

	var cardStrs []string
	for _, card := range cards {
		if card.Color == "Joker" {
			if card.Number == 15 {
				cardStrs = append(cardStrs, "SJ")
			} else if card.Number == 16 {
				cardStrs = append(cardStrs, "BJ")
			} else {
				cardStrs = append(cardStrs, card.Name)
			}
		} else {
			var suitStr string
			switch card.Color {
			case "Spade":
				suitStr = "S"
			case "Heart":
				suitStr = "H"
			case "Club":
				suitStr = "C"
			case "Diamond":
				suitStr = "D"
			default:
				suitStr = "?"
			}
			cardStrs = append(cardStrs, formatInt(card.RawNumber)+suitStr)
		}
	}

	return joinStrings(cardStrs, ",")
}

func formatInt(n int) string {
	if n < 0 {
		return "-" + formatInt(-n)
	}
	if n < 10 {
		return string(byte('0' + n))
	}
	return formatInt(n/10) + string(byte('0'+n%10))
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}