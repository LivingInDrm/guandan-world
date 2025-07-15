#!/usr/bin/env python3

# 颜色定义
colors = ['Spade', 'Heart', 'Club', 'Diamond', 'Joker']
name_map = {1: 'A', 11: 'J', 12: 'Q', 13: 'K', 14: 'A', 15: 'Red Joker', 16: 'Black Joker'}

class Card:
    def __init__(self, number, color, level):
        assert color in colors
        self.number = number
        self.color = color
        self.level = level
        self.raw_number = number if number != level else 100
        
        if 2 <= number <= 10:
            self.name = str(number)
        else:
            self.name = name_map[self.number]

    def clone(self):
        return Card(self.number, self.color, self.level)

    def is_wildcard(self):
        return self.number == self.level and self.color == 'Heart'

    def greater_than(self, card):
        # level card: this level's special number, as the greatest number other than the jokers
        assert isinstance(card, Card)
        assert isinstance(self.level, int)
        if card.number == self.level:
            if self.number >= 15:
                return True
            else:
                return False
        else:
            if self.number == self.level:
                if card.number <= 14:
                    return True
                else:
                    return False
            else:
                return self.number > card.number

    def consecutive_greater_than(self, card):
        return self.raw_number > card.raw_number

    def __lt__(self, card):
        if card.greater_than(self):
            return True
        elif self.equals(card) and card.color == 'Heart' and self.color != 'Heart':
            return True
        return False

    def __gt__(self, card):
        if self.greater_than(card):
            return True
        elif self.equals(card) and self.color == 'Heart' and card.color != 'Heart':
            return True
        return False

    def equals(self, card):
        return self.number == card.number

    def __str__(self):
        if self.color != 'Joker':
            return self.name + ' of ' + self.color
        else:
            return self.name

    def __repr__(self):
        return str(self)

class CardComp:
    def __init__(self, cards):
        assert isinstance(cards, list)

    def greater_than(self, card_comp):
        pass

    def __str__(self):
        return type(self).__name__ + ': ' + str(self.cards)

    def __repr__(self):
        return str(self)

    @staticmethod
    def is_bomb():
        pass

class NaiveBomb(CardComp):
    '''
    At least 4 cards with the same number
    '''

    def __init__(self, cards):
        super().__init__(cards)
        whether, sorted_cards = self.satisfy(cards)
        self.valid = whether
        self.cards = sorted_cards

    def greater_than(self, card_comp):
        # if the other is not bomb
        if not card_comp.is_bomb():
            return True

        # if the other is actually bomb
        # JokerBomb > everything
        if type(card_comp).__name__ == 'JokerBomb':
            return False
        # 6 Bomb > StraightFlush > 5 Bomb
        elif type(card_comp).__name__ == 'StraightFlush':
            if len(self.cards) >= 6:
                return True
            else:
                return False
        # Also NaiveBomb, compare number of cards first
        # then actual value
        else:
            if len(self.cards) > len(card_comp.cards):
                return True
            else:
                return self.cards[0].greater_than(card_comp.cards[0])

    @staticmethod
    def satisfy(cards):
        sorted_cards = sorted(cards)
        card_numbers = [card.number for card in sorted_cards]
        num_wildcards = sum([card.is_wildcard() for card in sorted_cards])
        if len(cards) < 4:
            return False, sorted(cards)
        elif num_wildcards == 0 and len(set(card_numbers)) == 1:
            return True, sorted_cards
        elif num_wildcards == 1 and len(set(card_numbers[0:-1])) == 1:
            return True, sorted_cards
        elif num_wildcards == 2 and len(set(card_numbers[0:-2])) == 1:
            return True, sorted_cards
        else:
            return False, sorted(cards)

    @staticmethod
    def is_bomb():
        return True

def test_naive_bomb_case():
    """测试异常的NaiveBomb比较case"""
    
    # 设置level为5
    level = 5
    
    # Comp1: [1♠, 5♥, 1♣, 1♦, 5♥] - 5张A炸弹（含2张变化牌）
    comp1_cards = [
        Card(1, "Spade", level),      # A♠
        Card(5, "Heart", level),      # 5♥ (变化牌)
        Card(1, "Club", level),       # A♣
        Card(1, "Diamond", level),    # A♦
        Card(5, "Heart", level),      # 5♥ (变化牌)
    ]
    
    # Comp2: [3♠, 3♠, 3♣, 3♦, 3♦, 3♥, 5♥, 5♥] - 8张3炸弹（含2张变化牌）
    comp2_cards = [
        Card(3, "Spade", level),      # 3♠
        Card(3, "Spade", level),      # 3♠
        Card(3, "Club", level),       # 3♣
        Card(3, "Diamond", level),    # 3♦
        Card(3, "Diamond", level),    # 3♦
        Card(3, "Heart", level),      # 3♥
        Card(5, "Heart", level),      # 5♥ (变化牌)
        Card(5, "Heart", level),      # 5♥ (变化牌)
    ]
    
    # 创建NaiveBomb对象
    bomb1 = NaiveBomb(comp1_cards)
    bomb2 = NaiveBomb(comp2_cards)
    
    print("=== 测试case详情 ===")
    print(f"Level: {level}")
    print(f"Comp1 (5张A炸弹): {bomb1}")
    print(f"  - 有效性: {bomb1.valid}")
    print(f"  - 张数: {len(bomb1.cards)}")
    print(f"  - 第一张牌: {bomb1.cards[0]} (number={bomb1.cards[0].number})")
    print()
    
    print(f"Comp2 (8张3炸弹): {bomb2}")
    print(f"  - 有效性: {bomb2.valid}")
    print(f"  - 张数: {len(bomb2.cards)}")
    print(f"  - 第一张牌: {bomb2.cards[0]} (number={bomb2.cards[0].number})")
    print()
    
    # 进行比较
    comp1_gt_comp2 = bomb1.greater_than(bomb2)
    comp2_gt_comp1 = bomb2.greater_than(bomb1)
    
    print("=== 实际比较结果 ===")
    print(f"comp1 > comp2: {comp1_gt_comp2}")
    print(f"comp2 > comp1: {comp2_gt_comp1}")
    print()
    
    print("=== 期望结果（来自测试数据）===")
    print("comp1 > comp2: True")
    print("comp2 > comp1: True")
    print()
    
    print("=== 分析 ===")
    if comp1_gt_comp2 and comp2_gt_comp1:
        print("❌ 异常：两个牌型互相都比对方大！这违反了比较的传递性。")
    elif comp1_gt_comp2:
        print("✅ 正常：5张A炸弹 > 8张3炸弹")
    elif comp2_gt_comp1:
        print("✅ 正常：8张3炸弹 > 5张A炸弹")
    else:
        print("⚠️  异常：两个牌型相等")
    
    # 详细分析比较过程
    print("\n=== 详细比较过程 ===")
    print(f"1. 张数比较: {len(bomb1.cards)} vs {len(bomb2.cards)}")
    if len(bomb1.cards) != len(bomb2.cards):
        longer_wins = len(bomb1.cards) > len(bomb2.cards)
        print(f"   张数不同，较多张数的获胜: bomb1胜利={longer_wins}")
    else:
        print("   张数相同，比较牌面价值")
        first_card_comparison = bomb1.cards[0].greater_than(bomb2.cards[0])
        print(f"   第一张牌比较: {bomb1.cards[0]} > {bomb2.cards[0]} = {first_card_comparison}")
    
    # 检查变化牌
    print(f"\n=== 变化牌检查 ===")
    comp1_wildcards = sum(1 for card in bomb1.cards if card.is_wildcard())
    comp2_wildcards = sum(1 for card in bomb2.cards if card.is_wildcard())
    print(f"Comp1 变化牌数量: {comp1_wildcards}")
    print(f"Comp2 变化牌数量: {comp2_wildcards}")
    
    # 检查排序后的卡片
    print(f"\n=== 排序后的卡片 ===")
    print("Comp1 排序后:")
    for i, card in enumerate(bomb1.cards):
        wildcard_status = "(变化牌)" if card.is_wildcard() else ""
        print(f"  [{i}]: {card} - number={card.number} {wildcard_status}")
    
    print("Comp2 排序后:")
    for i, card in enumerate(bomb2.cards):
        wildcard_status = "(变化牌)" if card.is_wildcard() else ""
        print(f"  [{i}]: {card} - number={card.number} {wildcard_status}")

if __name__ == "__main__":
    test_naive_bomb_case() 