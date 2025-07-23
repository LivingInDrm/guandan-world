package sdk

import (
	"testing"
)

func TestNewSingle(t *testing.T) {
	// 测试单张
	card, _ := NewCard(5, "Spade", 2)
	single := NewSingle([]*Card{card})

	if !single.IsValid() {
		t.Error("Single should be valid")
	}

	if single.GetType() != TypeSingle {
		t.Error("Type should be Single")
	}

	if single.IsBomb() {
		t.Error("Single should not be bomb")
	}

	// 测试单张比较
	card2, _ := NewCard(7, "Heart", 2)
	single2 := NewSingle([]*Card{card2})

	if !single2.GreaterThan(single) {
		t.Error("7 should be greater than 5")
	}

	if single.GreaterThan(single2) {
		t.Error("5 should not be greater than 7")
	}
}

func TestNewPair(t *testing.T) {
	// 测试正常对子
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	pair := NewPair([]*Card{card1, card2})

	if !pair.IsValid() {
		t.Error("Pair should be valid")
	}

	if pair.GetType() != TypePair {
		t.Error("Type should be Pair")
	}

	// 测试变化牌对子
	wildcard, _ := NewCard(2, "Heart", 2) // 变化牌
	normalCard, _ := NewCard(5, "Spade", 2)
	wildPair := NewPair([]*Card{wildcard, normalCard})

	if !wildPair.IsValid() {
		t.Error("Wildcard pair should be valid")
	}

	// 测试对子比较
	card3, _ := NewCard(7, "Spade", 2)
	card4, _ := NewCard(7, "Heart", 2)
	pair2 := NewPair([]*Card{card3, card4})

	if !pair2.GreaterThan(pair) {
		t.Error("7 pair should be greater than 5 pair")
	}
}

func TestNewTriple(t *testing.T) {
	// 测试正常三张
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	card3, _ := NewCard(5, "Diamond", 2)
	triple := NewTriple([]*Card{card1, card2, card3})

	if !triple.IsValid() {
		t.Error("Triple should be valid")
	}

	if triple.GetType() != TypeTriple {
		t.Error("Type should be Triple")
	}

	// 测试包含王的三张（应该无效）
	joker, _ := NewCard(15, "Joker", 2)
	invalidTriple := NewTriple([]*Card{card1, card2, joker})

	if invalidTriple.IsValid() {
		t.Error("Triple with joker should be invalid")
	}
}

func TestNewFullHouse(t *testing.T) {
	// 测试正常葫芦（三带二）
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	card3, _ := NewCard(5, "Diamond", 2)
	card4, _ := NewCard(7, "Spade", 2)
	card5, _ := NewCard(7, "Heart", 2)

	fullhouse := NewFullHouse([]*Card{card1, card2, card3, card4, card5})

	if !fullhouse.IsValid() {
		t.Error("FullHouse should be valid")
	}

	if fullhouse.GetType() != TypeFullHouse {
		t.Error("Type should be FullHouse")
	}

	if fullhouse.IsBomb() {
		t.Error("FullHouse should not be bomb")
	}
}

func TestNewStraight(t *testing.T) {
	// 测试正常顺子
	card1, _ := NewCard(3, "Spade", 2)
	card2, _ := NewCard(4, "Heart", 2)
	card3, _ := NewCard(5, "Diamond", 2)
	card4, _ := NewCard(6, "Spade", 2)
	card5, _ := NewCard(7, "Heart", 2)

	straight := NewStraight([]*Card{card1, card2, card3, card4, card5})

	if !straight.IsValid() {
		t.Error("Straight should be valid")
	}

	if straight.GetType() != TypeStraight {
		t.Error("Type should be Straight")
	}

	if straight.IsBomb() {
		t.Error("Straight should not be bomb")
	}
}

func TestNewJokerBomb(t *testing.T) {
	// 测试王炸
	joker1, _ := NewCard(15, "Joker", 2) // 小王
	joker2, _ := NewCard(15, "Joker", 2) // 小王
	joker3, _ := NewCard(16, "Joker", 2) // 大王
	joker4, _ := NewCard(16, "Joker", 2) // 大王

	jokerBomb := NewJokerBomb([]*Card{joker1, joker2, joker3, joker4})

	if !jokerBomb.IsValid() {
		t.Error("JokerBomb should be valid")
	}

	if jokerBomb.GetType() != TypeJokerBomb {
		t.Error("Type should be JokerBomb")
	}

	if !jokerBomb.IsBomb() {
		t.Error("JokerBomb should be bomb")
	}
}

func TestNewNaiveBomb(t *testing.T) {
	// 测试4张炸弹
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	card3, _ := NewCard(5, "Diamond", 2)
	card4, _ := NewCard(5, "Club", 2)

	bomb := NewNaiveBomb([]*Card{card1, card2, card3, card4})

	if !bomb.IsValid() {
		t.Error("NaiveBomb should be valid")
	}

	if bomb.GetType() != TypeNaiveBomb {
		t.Error("Type should be NaiveBomb")
	}

	if !bomb.IsBomb() {
		t.Error("NaiveBomb should be bomb")
	}

	// 测试5张炸弹
	card5, _ := NewCard(7, "Spade", 2)
	bomb5 := NewNaiveBomb([]*Card{card5, card5, card5, card5, card5})

	if !bomb5.IsValid() {
		t.Error("5-card NaiveBomb should be valid")
	}
}

func TestNewStraightFlush(t *testing.T) {
	// 测试同花顺
	card1, _ := NewCard(3, "Spade", 2)
	card2, _ := NewCard(4, "Spade", 2)
	card3, _ := NewCard(5, "Spade", 2)
	card4, _ := NewCard(6, "Spade", 2)
	card5, _ := NewCard(7, "Spade", 2)

	straightFlush := NewStraightFlush([]*Card{card1, card2, card3, card4, card5})

	if !straightFlush.IsValid() {
		t.Error("StraightFlush should be valid")
	}

	if straightFlush.GetType() != TypeStraightFlush {
		t.Error("Type should be StraightFlush")
	}

	if !straightFlush.IsBomb() {
		t.Error("StraightFlush should be bomb")
	}
}

func TestNewPlate(t *testing.T) {
	// 测试钢板（连续三张）
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	card3, _ := NewCard(5, "Diamond", 2)
	card4, _ := NewCard(6, "Spade", 2)
	card5, _ := NewCard(6, "Heart", 2)
	card6, _ := NewCard(6, "Diamond", 2)

	plate := NewPlate([]*Card{card1, card2, card3, card4, card5, card6})

	if !plate.IsValid() {
		t.Error("Plate should be valid")
	}

	if plate.GetType() != TypePlate {
		t.Error("Type should be Plate")
	}

	if plate.IsBomb() {
		t.Error("Plate should not be bomb")
	}
}

func TestNewTube(t *testing.T) {
	// 测试钢管（连续对子）
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	card3, _ := NewCard(6, "Spade", 2)
	card4, _ := NewCard(6, "Heart", 2)
	card5, _ := NewCard(7, "Spade", 2)
	card6, _ := NewCard(7, "Heart", 2)

	tube := NewTube([]*Card{card1, card2, card3, card4, card5, card6})

	if !tube.IsValid() {
		t.Error("Tube should be valid")
	}

	if tube.GetType() != TypeTube {
		t.Error("Type should be Tube")
	}

	if tube.IsBomb() {
		t.Error("Tube should not be bomb")
	}
}

func TestFromCardList(t *testing.T) {
	// 测试从牌列表生成牌组

	// 测试单张
	card1, _ := NewCard(5, "Spade", 2)
	comp1 := FromCardList([]*Card{card1}, nil)
	if comp1.GetType() != TypeSingle {
		t.Error("Should create Single")
	}

	// 测试对子
	card2, _ := NewCard(5, "Heart", 2)
	comp2 := FromCardList([]*Card{card1, card2}, nil)
	if comp2.GetType() != TypePair {
		t.Error("Should create Pair")
	}

	// 测试王炸
	joker1, _ := NewCard(15, "Joker", 2)
	joker2, _ := NewCard(15, "Joker", 2)
	joker3, _ := NewCard(16, "Joker", 2)
	joker4, _ := NewCard(16, "Joker", 2)
	comp3 := FromCardList([]*Card{joker1, joker2, joker3, joker4}, nil)
	if comp3.GetType() != TypeJokerBomb {
		t.Error("Should create JokerBomb")
	}

	// 测试弃牌
	comp4 := FromCardList([]*Card{}, nil)
	if comp4.GetType() != TypeFold {
		t.Error("Should create Fold")
	}

	// 测试非法牌组
	invalidCards := []*Card{card1, card2, joker1} // 3张不同的牌
	comp5 := FromCardList(invalidCards, nil)
	if comp5.GetType() != TypeIllegal {
		t.Error("Should create IllegalComp")
	}
}

func TestCompareComps(t *testing.T) {
	// 测试牌组比较

	// 创建不同类型的牌组
	card1, _ := NewCard(5, "Spade", 2)
	single := NewSingle([]*Card{card1})

	card2, _ := NewCard(7, "Heart", 2)
	single2 := NewSingle([]*Card{card2})

	// 4张炸弹
	card3, _ := NewCard(5, "Spade", 2)
	card4, _ := NewCard(5, "Heart", 2)
	card5, _ := NewCard(5, "Diamond", 2)
	card6, _ := NewCard(5, "Club", 2)
	bomb := NewNaiveBomb([]*Card{card3, card4, card5, card6})

	// 王炸
	joker1, _ := NewCard(15, "Joker", 2)
	joker2, _ := NewCard(15, "Joker", 2)
	joker3, _ := NewCard(16, "Joker", 2)
	joker4, _ := NewCard(16, "Joker", 2)
	jokerBomb := NewJokerBomb([]*Card{joker1, joker2, joker3, joker4})

	// 测试单张比较
	if !single2.GreaterThan(single) {
		t.Error("7 should be greater than 5")
	}

	// 测试炸弹 > 单张
	if !bomb.GreaterThan(single) {
		t.Error("Bomb should be greater than single")
	}

	// 测试王炸 > 炸弹
	if !jokerBomb.GreaterThan(bomb) {
		t.Error("JokerBomb should be greater than bomb")
	}

	// 测试王炸 vs 王炸
	if jokerBomb.GreaterThan(jokerBomb) {
		t.Error("JokerBomb should not be greater than itself")
	}
}

func TestBombLogic(t *testing.T) {
	// 测试炸弹逻辑

	// 4张炸弹
	card1, _ := NewCard(5, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	card3, _ := NewCard(5, "Diamond", 2)
	card4, _ := NewCard(5, "Club", 2)
	bomb4 := NewNaiveBomb([]*Card{card1, card2, card3, card4})

	// 5张炸弹
	card5, _ := NewCard(7, "Spade", 2)
	card6, _ := NewCard(7, "Heart", 2)
	card7, _ := NewCard(7, "Diamond", 2)
	card8, _ := NewCard(7, "Club", 2)
	card9, _ := NewCard(7, "Spade", 2) // 复制一张，实际应该是变化牌
	bomb5 := NewNaiveBomb([]*Card{card5, card6, card7, card8, card9})

	// 测试5张炸弹 > 4张炸弹
	if !bomb5.GreaterThan(bomb4) {
		t.Error("5-card bomb should be greater than 4-card bomb")
	}

	// 测试同花顺
	sCard1, _ := NewCard(3, "Spade", 2)
	sCard2, _ := NewCard(4, "Spade", 2)
	sCard3, _ := NewCard(5, "Spade", 2)
	sCard4, _ := NewCard(6, "Spade", 2)
	sCard5, _ := NewCard(7, "Spade", 2)
	straightFlush := NewStraightFlush([]*Card{sCard1, sCard2, sCard3, sCard4, sCard5})

	// 测试同花顺 > 5张以下炸弹
	if !straightFlush.GreaterThan(bomb4) {
		t.Error("StraightFlush should be greater than 4-card bomb")
	}

	// 测试6张炸弹 > 同花顺
	card10, _ := NewCard(9, "Spade", 2)
	card11, _ := NewCard(9, "Heart", 2)
	card12, _ := NewCard(9, "Diamond", 2)
	card13, _ := NewCard(9, "Club", 2)
	card14, _ := NewCard(9, "Spade", 2)
	card15, _ := NewCard(9, "Heart", 2)
	bomb6 := NewNaiveBomb([]*Card{card10, card11, card12, card13, card14, card15})

	if !bomb6.GreaterThan(straightFlush) {
		t.Error("6-card bomb should be greater than StraightFlush")
	}
}

func TestFoldAndIllegal(t *testing.T) {
	// 测试弃牌
	fold := &Fold{BaseComp: BaseComp{Cards: []*Card{}, Valid: true, Type: TypeFold}}

	if fold.GreaterThan(fold) {
		t.Error("Fold should not be greater than anything")
	}

	if fold.IsBomb() {
		t.Error("Fold should not be bomb")
	}

	// 测试非法牌组
	illegal := &IllegalComp{BaseComp: BaseComp{Cards: []*Card{}, Valid: false, Type: TypeIllegal}}

	if illegal.GreaterThan(illegal) {
		t.Error("IllegalComp should not be greater than anything")
	}

	if illegal.IsBomb() {
		t.Error("IllegalComp should not be bomb")
	}
}
