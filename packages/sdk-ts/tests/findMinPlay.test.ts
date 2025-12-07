import { describe, it, expect } from 'vitest';
import { findMinPlay } from '../src/findMinPlay';
import { SdkCard } from '../src/card';
import { fromCardList } from '../src/fromCardList';
import { CompType } from '../src/compType';
import { Suit } from '../src/types';

describe('findMinPlay', () => {
  const level = 5;

  it('should return null when no valid play exists (framework placeholder)', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 0, level, 1),
      new SdkCard(6, 0, level, 2),
    ];

    const prevCards = [new SdkCard(7, 0, level, 10)];
    const prev = fromCardList(prevCards);

    const result = findMinPlay(handCards, prev);
    expect(result).toBeNull();
  });

  it('should be a function that accepts cards and prev', () => {
    expect(typeof findMinPlay).toBe('function');
  });
});

describe('findMinBomb', () => {
  const level = 5;

  it('should find 4-card NaiveBomb to beat a Single', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(3, 3, level, 3),
    ];

    const prevCards = [new SdkCard(14, 0, level, 10)];
    const prev = fromCardList(prevCards);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
    expect(result!.getCards().length).toBe(4);
  });

  it('should find the smallest NaiveBomb among multiple options', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(3, 3, level, 3),
      new SdkCard(7, 0, level, 4),
      new SdkCard(7, 1, level, 5),
      new SdkCard(7, 2, level, 6),
      new SdkCard(7, 3, level, 7),
    ];

    const prevCards = [new SdkCard(14, 0, level, 10)];
    const prev = fromCardList(prevCards);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
    expect(result!.getCards()[0].rank).toBe(3);
  });

  it('should find bigger NaiveBomb to beat smaller NaiveBomb', () => {
    const handCards = [
      new SdkCard(7, 0, level, 0),
      new SdkCard(7, 1, level, 1),
      new SdkCard(7, 2, level, 2),
      new SdkCard(7, 3, level, 3),
    ];

    const prevCards = [
      new SdkCard(3, 0, level, 10),
      new SdkCard(3, 1, level, 11),
      new SdkCard(3, 2, level, 12),
      new SdkCard(3, 3, level, 13),
    ];
    const prev = fromCardList(prevCards);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
    expect(result!.getCards()[0].rank).toBe(7);
  });

  it('should return null when no bomb can beat prev bomb', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(3, 3, level, 3),
    ];

    const prevCards = [
      new SdkCard(7, 0, level, 10),
      new SdkCard(7, 1, level, 11),
      new SdkCard(7, 2, level, 12),
      new SdkCard(7, 3, level, 13),
    ];
    const prev = fromCardList(prevCards);

    const result = findMinPlay(handCards, prev);
    expect(result).toBeNull();
  });

  it('should find StraightFlush to beat a smaller bomb', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 0, level, 1),
      new SdkCard(5, 0, level, 2),
      new SdkCard(6, 0, level, 3),
      new SdkCard(7, 0, level, 4),
    ];

    const prevCards = [
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(2, 2, level, 12),
      new SdkCard(2, 3, level, 13),
    ];
    const prev = fromCardList(prevCards);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.StraightFlush);
  });

  it('should find JokerBomb to beat a big bomb', () => {
    const handCards = [
      new SdkCard(15, Suit.Joker, level, 0),
      new SdkCard(15, Suit.Joker, level, 1),
      new SdkCard(16, Suit.Joker, level, 2),
      new SdkCard(16, Suit.Joker, level, 3),
    ];

    const prevCards = [
      new SdkCard(14, 0, level, 0),
      new SdkCard(14, 1, level, 0),
      new SdkCard(14, 2, level, 0),
      new SdkCard(14, 3, level, 0),
      new SdkCard(14, 0, level, 1),
      new SdkCard(14, 1, level, 1),
      new SdkCard(14, 2, level, 1),
      new SdkCard(14, 3, level, 1),
    ];
    const prev = fromCardList(prevCards);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.JokerBomb);
  });

  it('should find NaiveBomb with wildcard to beat a smaller bomb', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(level, Suit.Heart, level, 3),
    ];

    const prevCards = [
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(2, 2, level, 12),
      new SdkCard(2, 3, level, 13),
    ];
    const prev = fromCardList(prevCards);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
    expect(result!.getCards().length).toBe(4);
  });

  it('should prefer smaller bomb over larger bomb when beating a bomb', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(3, 3, level, 3),
      new SdkCard(7, 0, level, 4),
      new SdkCard(7, 1, level, 5),
      new SdkCard(7, 2, level, 6),
      new SdkCard(7, 3, level, 7),
      new SdkCard(level, Suit.Heart, level, 8),
    ];

    const prevCards = [
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(2, 2, level, 12),
      new SdkCard(2, 3, level, 13),
    ];
    const prev = fromCardList(prevCards);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
    expect(result!.getCards().length).toBe(4);
  });
});

describe('findMinBomb with level=5 and 2 wildcards (27 cards)', () => {
  const level = 5;

  const handCards = [
    new SdkCard(5, Suit.Heart, level, 0),
    new SdkCard(5, Suit.Heart, level, 1),
    new SdkCard(7, 0, level, 0),
    new SdkCard(7, 0, level, 1),
    new SdkCard(7, Suit.Heart, level, 0),
    new SdkCard(7, Suit.Club, level, 0),
    new SdkCard(13, 0, level, 0),
    new SdkCard(13, Suit.Heart, level, 0),
    new SdkCard(13, Suit.Club, level, 0),
    new SdkCard(13, Suit.Club, level, 1),
    new SdkCard(8, 0, level, 0),
    new SdkCard(8, Suit.Heart, level, 0),
    new SdkCard(8, Suit.Club, level, 0),
    new SdkCard(3, 0, level, 0),
    new SdkCard(3, 0, level, 1),
    new SdkCard(4, 0, level, 0),
    new SdkCard(4, 0, level, 1),
    new SdkCard(6, 0, level, 0),
    new SdkCard(6, 0, level, 1),
    new SdkCard(9, Suit.Heart, level, 0),
    new SdkCard(9, Suit.Club, level, 0),
    new SdkCard(10, Suit.Heart, level, 0),
    new SdkCard(10, Suit.Club, level, 0),
    new SdkCard(15, Suit.Joker, level, 0),
    new SdkCard(15, Suit.Joker, level, 1),
    new SdkCard(16, Suit.Joker, level, 0),
    new SdkCard(16, Suit.Joker, level, 1),
  ];

  it('场景1: prev = 4个9', () => {
    const prev = fromCardList([
      new SdkCard(9, 0, level, 0),
      new SdkCard(9, Suit.Heart, level, 0),
      new SdkCard(9, Suit.Club, level, 0),
      new SdkCard(9, Suit.Diamond, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
    expect(result!.getCards().length).toBe(4);
    expect(result!.getCards()[0].rank).toBe(10);
  });

  it('场景2: prev = 5个10', () => {
    const prev = fromCardList([
      new SdkCard(10, 0, level, 0),
      new SdkCard(10, 0, level, 1),
      new SdkCard(10, Suit.Heart, level, 0),
      new SdkCard(10, Suit.Club, level, 0),
      new SdkCard(10, Suit.Diamond, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
    expect(result!.getCards().length).toBe(5);
    expect(result!.getCards()[0].rank).toBe(13);
  });

  it('场景3: prev = 6个5', () => {
    const prev = fromCardList([
      new SdkCard(5, 0, level, 0),
      new SdkCard(5, 0, level, 1),
      new SdkCard(5, Suit.Club, level, 0),
      new SdkCard(5, Suit.Club, level, 1),
      new SdkCard(5, Suit.Diamond, level, 0),
      new SdkCard(5, Suit.Diamond, level, 1),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.JokerBomb);
  });

  it('场景4: prev = 7个6', () => {
    const prev = fromCardList([
      new SdkCard(6, 0, level, 0),
      new SdkCard(6, 0, level, 1),
      new SdkCard(6, Suit.Heart, level, 0),
      new SdkCard(6, Suit.Heart, level, 1),
      new SdkCard(6, Suit.Club, level, 0),
      new SdkCard(6, Suit.Club, level, 1),
      new SdkCard(6, Suit.Diamond, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.JokerBomb);
  });

  it('场景5: prev = 8个2', () => {
    const prev = fromCardList([
      new SdkCard(2, 0, level, 0),
      new SdkCard(2, 0, level, 1),
      new SdkCard(2, Suit.Heart, level, 0),
      new SdkCard(2, Suit.Heart, level, 1),
      new SdkCard(2, Suit.Club, level, 0),
      new SdkCard(2, Suit.Club, level, 1),
      new SdkCard(2, Suit.Diamond, level, 0),
      new SdkCard(2, Suit.Diamond, level, 1),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.JokerBomb);
  });

  it('场景6: prev = Heart 3,4,5,6,7 (红桃同花顺)', () => {
    const prev = fromCardList([
      new SdkCard(3, Suit.Heart, level, 0),
      new SdkCard(4, Suit.Heart, level, 0),
      new SdkCard(5, Suit.Heart, level, 0),
      new SdkCard(6, Suit.Heart, level, 0),
      new SdkCard(7, Suit.Heart, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.StraightFlush);
  });

  it('场景7: prev = Club 5,6,7,8,9 (梅花同花顺)', () => {
    const prev = fromCardList([
      new SdkCard(5, Suit.Club, level, 0),
      new SdkCard(6, Suit.Club, level, 0),
      new SdkCard(7, Suit.Club, level, 0),
      new SdkCard(8, Suit.Club, level, 0),
      new SdkCard(9, Suit.Club, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.StraightFlush);
  });
});

describe('findMinSameNumber', () => {
  const level = 5;

  describe('Single', () => {
    it('should find smallest single to beat prev single', () => {
      const handCards = [
        new SdkCard(3, 0, level, 0),
        new SdkCard(7, 0, level, 1),
        new SdkCard(10, 0, level, 2),
      ];
      const prev = fromCardList([new SdkCard(6, 0, level, 10)]);

      const result = findMinPlay(handCards, prev);
      expect(result).not.toBeNull();
      expect(result!.getType()).toBe(CompType.Single);
      expect(result!.getCards()[0].rank).toBe(7);
    });

    it('should return null when no single can beat prev', () => {
      const handCards = [
        new SdkCard(3, 0, level, 0),
        new SdkCard(4, 0, level, 1),
      ];
      const prev = fromCardList([new SdkCard(10, 0, level, 10)]);

      const result = findMinPlay(handCards, prev);
      expect(result).toBeNull();
    });

    it('should use wildcard as single', () => {
      const handCards = [
        new SdkCard(3, 0, level, 0),
        new SdkCard(level, Suit.Heart, level, 1),
      ];
      const prev = fromCardList([new SdkCard(10, 0, level, 10)]);

      const result = findMinPlay(handCards, prev);
      expect(result).not.toBeNull();
      expect(result!.getType()).toBe(CompType.Single);
      expect(result!.getCards()[0].isWildcard()).toBe(true);
    });
  });

  describe('Pair', () => {
    it('should find smallest pair to beat prev pair', () => {
      const handCards = [
        new SdkCard(3, 0, level, 0),
        new SdkCard(3, 1, level, 1),
        new SdkCard(7, 0, level, 2),
        new SdkCard(7, 1, level, 3),
        new SdkCard(10, 0, level, 4),
        new SdkCard(10, 1, level, 5),
      ];
      const prev = fromCardList([
        new SdkCard(6, 0, level, 10),
        new SdkCard(6, 1, level, 11),
      ]);

      const result = findMinPlay(handCards, prev);
      expect(result).not.toBeNull();
      expect(result!.getType()).toBe(CompType.Pair);
      expect(result!.getCards()[0].rank).toBe(7);
    });

    it('should find pair with wildcard', () => {
      const handCards = [
        new SdkCard(7, 0, level, 0),
        new SdkCard(level, Suit.Heart, level, 1),
      ];
      const prev = fromCardList([
        new SdkCard(6, 0, level, 10),
        new SdkCard(6, 1, level, 11),
      ]);

      const result = findMinPlay(handCards, prev);
      expect(result).not.toBeNull();
      expect(result!.getType()).toBe(CompType.Pair);
    });

    it('should return null when no pair can beat prev', () => {
      const handCards = [
        new SdkCard(3, 0, level, 0),
        new SdkCard(3, 1, level, 1),
      ];
      const prev = fromCardList([
        new SdkCard(10, 0, level, 10),
        new SdkCard(10, 1, level, 11),
      ]);

      const result = findMinPlay(handCards, prev);
      expect(result).toBeNull();
    });
  });

  describe('Triple', () => {
    it('should find smallest triple to beat prev triple', () => {
      const handCards = [
        new SdkCard(3, 0, level, 0),
        new SdkCard(3, 1, level, 1),
        new SdkCard(3, 2, level, 2),
        new SdkCard(8, 0, level, 3),
        new SdkCard(8, 1, level, 4),
        new SdkCard(8, 2, level, 5),
      ];
      const prev = fromCardList([
        new SdkCard(6, 0, level, 10),
        new SdkCard(6, 1, level, 11),
        new SdkCard(6, 2, level, 12),
      ]);

      const result = findMinPlay(handCards, prev);
      expect(result).not.toBeNull();
      expect(result!.getType()).toBe(CompType.Triple);
      expect(result!.getCards()[0].rank).toBe(8);
    });

    it('should find triple with wildcard', () => {
      const handCards = [
        new SdkCard(8, 0, level, 0),
        new SdkCard(8, 1, level, 1),
        new SdkCard(level, Suit.Heart, level, 2),
      ];
      const prev = fromCardList([
        new SdkCard(6, 0, level, 10),
        new SdkCard(6, 1, level, 11),
        new SdkCard(6, 2, level, 12),
      ]);

      const result = findMinPlay(handCards, prev);
      expect(result).not.toBeNull();
      expect(result!.getType()).toBe(CompType.Triple);
    });

    it('should use pure wildcards as triple', () => {
      const handCards = [
        new SdkCard(level, Suit.Heart, level, 0),
        new SdkCard(level, Suit.Heart, level, 1),
        new SdkCard(level, Suit.Heart, level, 2),
      ];
      const prev = fromCardList([
        new SdkCard(10, 0, level, 10),
        new SdkCard(10, 1, level, 11),
        new SdkCard(10, 2, level, 12),
      ]);

      const result = findMinPlay(handCards, prev);
      expect(result).not.toBeNull();
      expect(result!.getType()).toBe(CompType.Triple);
    });
  });
});

describe('findMinSameType - Single (27 cards)', () => {
  const level = 5;

  const handCards = [
    new SdkCard(5, Suit.Heart, level, 0),
    new SdkCard(5, Suit.Heart, level, 1),
    new SdkCard(7, 0, level, 0),
    new SdkCard(7, 0, level, 1),
    new SdkCard(7, Suit.Heart, level, 0),
    new SdkCard(7, Suit.Club, level, 0),
    new SdkCard(13, 0, level, 0),
    new SdkCard(13, Suit.Heart, level, 0),
    new SdkCard(13, Suit.Club, level, 0),
    new SdkCard(13, Suit.Club, level, 1),
    new SdkCard(8, 0, level, 0),
    new SdkCard(8, Suit.Heart, level, 0),
    new SdkCard(8, Suit.Club, level, 0),
    new SdkCard(3, 0, level, 0),
    new SdkCard(3, 0, level, 1),
    new SdkCard(4, 0, level, 0),
    new SdkCard(4, 0, level, 1),
    new SdkCard(6, 0, level, 0),
    new SdkCard(6, 0, level, 1),
    new SdkCard(9, Suit.Heart, level, 0),
    new SdkCard(9, Suit.Club, level, 0),
    new SdkCard(10, Suit.Heart, level, 0),
    new SdkCard(10, Suit.Club, level, 0),
    new SdkCard(15, Suit.Joker, level, 0),
    new SdkCard(15, Suit.Joker, level, 1),
    new SdkCard(16, Suit.Joker, level, 0),
    new SdkCard(16, Suit.Joker, level, 1),
  ];

  it('场景1: prev = 单张黑桃3', () => {
    const prev = fromCardList([new SdkCard(3, 0, level, 0)]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Single);
    expect(result!.getCards()[0].rank).toBe(4);
  });

  it('场景2: prev = 单张黑桃K', () => {
    const prev = fromCardList([new SdkCard(13, 0, level, 0)]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Single);
    expect(result!.getCards()[0].isWildcard()).toBe(true);
  });

  it('场景3: prev = 单张A', () => {
    const prev = fromCardList([new SdkCard(14, 0, level, 0)]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Single);
    expect(result!.getCards()[0].isWildcard()).toBe(true);
  });

  it('场景4: prev = 单张小王', () => {
    const prev = fromCardList([new SdkCard(15, Suit.Joker, level, 0)]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Single);
    expect(result!.getCards()[0].isBigJoker()).toBe(true);
  });

  it('场景5: prev = 单张大王', () => {
    const prev = fromCardList([new SdkCard(16, Suit.Joker, level, 0)]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
  });

  it('场景6: prev = 单张5 (级牌)', () => {
    const prev = fromCardList([new SdkCard(5, 0, level, 0)]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Single);
    expect(result!.getCards()[0].isSmallJoker()).toBe(true);
  });
});

describe('findMinSameType - Pair (27 cards)', () => {
  const level = 5;

  const handCards = [
    new SdkCard(5, Suit.Heart, level, 0),
    new SdkCard(5, Suit.Heart, level, 1),
    new SdkCard(7, 0, level, 0),
    new SdkCard(7, 0, level, 1),
    new SdkCard(7, Suit.Heart, level, 0),
    new SdkCard(7, Suit.Club, level, 0),
    new SdkCard(13, 0, level, 0),
    new SdkCard(13, Suit.Heart, level, 0),
    new SdkCard(13, Suit.Club, level, 0),
    new SdkCard(13, Suit.Club, level, 1),
    new SdkCard(8, 0, level, 0),
    new SdkCard(8, Suit.Heart, level, 0),
    new SdkCard(8, Suit.Club, level, 0),
    new SdkCard(3, 0, level, 0),
    new SdkCard(3, 0, level, 1),
    new SdkCard(4, 0, level, 0),
    new SdkCard(4, 0, level, 1),
    new SdkCard(6, 0, level, 0),
    new SdkCard(6, 0, level, 1),
    new SdkCard(9, Suit.Heart, level, 0),
    new SdkCard(9, Suit.Club, level, 0),
    new SdkCard(10, Suit.Heart, level, 0),
    new SdkCard(10, Suit.Club, level, 0),
    new SdkCard(15, Suit.Joker, level, 0),
    new SdkCard(15, Suit.Joker, level, 1),
    new SdkCard(16, Suit.Joker, level, 0),
    new SdkCard(16, Suit.Joker, level, 1),
  ];

  it('场景1: prev = 对3', () => {
    const prev = fromCardList([
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Pair);
    expect(result!.getCards()[0].rank).toBe(4);
  });

  it('场景2: prev = 对8', () => {
    const prev = fromCardList([
      new SdkCard(8, 0, level, 0),
      new SdkCard(8, 1, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Pair);
    expect(result!.getCards()[0].rank).toBe(9);
  });

  it('场景3: prev = 对K', () => {
    const prev = fromCardList([
      new SdkCard(13, 0, level, 0),
      new SdkCard(13, 1, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Pair);
    expect(result!.getCards().every(c => c.isWildcard())).toBe(true);
  });

  it('场景4: prev = 对小王', () => {
    const prev = fromCardList([
      new SdkCard(15, Suit.Joker, level, 0),
      new SdkCard(15, Suit.Joker, level, 1),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Pair);
    expect(result!.getCards().every(c => c.isBigJoker())).toBe(true);
  });

  it('场景5: prev = 对大王', () => {
    const prev = fromCardList([
      new SdkCard(16, Suit.Joker, level, 0),
      new SdkCard(16, Suit.Joker, level, 1),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
  });

  it('场景6: prev = 对5 (级牌)', () => {
    const prev = fromCardList([
      new SdkCard(5, 0, level, 0),
      new SdkCard(5, Suit.Club, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Pair);
    expect(result!.getCards().every(c => c.isSmallJoker())).toBe(true);
  });
});

describe('findMinSameType - Triple (27 cards)', () => {
  const level = 5;

  const handCards = [
    new SdkCard(5, Suit.Heart, level, 0),
    new SdkCard(5, Suit.Heart, level, 1),
    new SdkCard(7, 0, level, 0),
    new SdkCard(7, 0, level, 1),
    new SdkCard(7, Suit.Heart, level, 0),
    new SdkCard(7, Suit.Club, level, 0),
    new SdkCard(13, 0, level, 0),
    new SdkCard(13, Suit.Heart, level, 0),
    new SdkCard(13, Suit.Club, level, 0),
    new SdkCard(13, Suit.Club, level, 1),
    new SdkCard(8, 0, level, 0),
    new SdkCard(8, Suit.Heart, level, 0),
    new SdkCard(8, Suit.Club, level, 0),
    new SdkCard(3, 0, level, 0),
    new SdkCard(3, 0, level, 1),
    new SdkCard(4, 0, level, 0),
    new SdkCard(4, 0, level, 1),
    new SdkCard(6, 0, level, 0),
    new SdkCard(6, 0, level, 1),
    new SdkCard(9, Suit.Heart, level, 0),
    new SdkCard(9, Suit.Club, level, 0),
    new SdkCard(10, Suit.Heart, level, 0),
    new SdkCard(10, Suit.Club, level, 0),
    new SdkCard(15, Suit.Joker, level, 0),
    new SdkCard(15, Suit.Joker, level, 1),
    new SdkCard(16, Suit.Joker, level, 0),
    new SdkCard(16, Suit.Joker, level, 1),
  ];

  it('场景1: prev = 三张3', () => {
    const prev = fromCardList([
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
      new SdkCard(3, 2, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Triple);
    expect(result!.getCards().some(c => c.rank === 4)).toBe(true);
  });

  it('场景2: prev = 三张6', () => {
    const prev = fromCardList([
      new SdkCard(6, 0, level, 0),
      new SdkCard(6, 1, level, 0),
      new SdkCard(6, 2, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Triple);
    expect(result!.getCards().filter(c => c.rank === 7).length).toBeGreaterThanOrEqual(2);
  });

  it('场景3: prev = 三张K', () => {
    const prev = fromCardList([
      new SdkCard(13, 0, level, 0),
      new SdkCard(13, 1, level, 0),
      new SdkCard(13, 2, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
  });

  it('场景4: prev = 三张8 (测试wildcard凑三张)', () => {
    const prev = fromCardList([
      new SdkCard(8, 0, level, 0),
      new SdkCard(8, 1, level, 0),
      new SdkCard(8, 2, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Triple);
    expect(result!.getCards().some(c => c.rank === 9)).toBe(true);
  });
});

describe('findMinStraight', () => {
  const level = 5;

  it('should find smallest straight to beat prev straight', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 0, level, 1),
      new SdkCard(5, 0, level, 2),
      new SdkCard(6, 0, level, 3),
      new SdkCard(7, 0, level, 4),
      new SdkCard(8, 0, level, 5),
      new SdkCard(9, 0, level, 6),
      new SdkCard(10, 0, level, 7),
      new SdkCard(11, 0, level, 8),
    ];

    const prev = fromCardList([
      new SdkCard(2, 0, level, 10),
      new SdkCard(3, 1, level, 11),
      new SdkCard(4, 2, level, 12),
      new SdkCard(5, 3, level, 13),
      new SdkCard(6, 0, level, 14),
    ]);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Straight);
  });

  it('should find straight with wildcard', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 0, level, 1),
      new SdkCard(6, 0, level, 3),
      new SdkCard(7, 0, level, 4),
      new SdkCard(level, Suit.Heart, level, 5),
    ];

    const prev = fromCardList([
      new SdkCard(2, 0, level, 10),
      new SdkCard(3, 1, level, 11),
      new SdkCard(4, 2, level, 12),
      new SdkCard(5, 3, level, 13),
      new SdkCard(6, 0, level, 14),
    ]);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Straight);
  });

  it('should return null when no straight can beat prev', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 1, level, 1),
      new SdkCard(5, 2, level, 2),
      new SdkCard(6, 3, level, 3),
      new SdkCard(7, 0, level, 4),
    ];

    const prev = fromCardList([
      new SdkCard(6, 0, level, 10),
      new SdkCard(7, 1, level, 11),
      new SdkCard(8, 2, level, 12),
      new SdkCard(9, 3, level, 13),
      new SdkCard(10, 0, level, 14),
    ]);

    const result = findMinPlay(handCards, prev);
    expect(result).toBeNull();
  });

  it('should find straight 10-A', () => {
    const handCards = [
      new SdkCard(10, 0, level, 0),
      new SdkCard(11, 0, level, 1),
      new SdkCard(12, 0, level, 2),
      new SdkCard(13, 0, level, 3),
      new SdkCard(14, 0, level, 4),
    ];

    const prev = fromCardList([
      new SdkCard(8, 0, level, 10),
      new SdkCard(9, 1, level, 11),
      new SdkCard(10, 2, level, 12),
      new SdkCard(11, 3, level, 13),
      new SdkCard(12, 0, level, 14),
    ]);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Straight);
  });
});

describe('findMinSameType - Straight (27 cards)', () => {
  const level = 5;

  const handCards = [
    new SdkCard(5, Suit.Heart, level, 0),
    new SdkCard(5, Suit.Heart, level, 1),
    new SdkCard(7, 0, level, 0),
    new SdkCard(7, 0, level, 1),
    new SdkCard(7, Suit.Heart, level, 0),
    new SdkCard(7, Suit.Club, level, 0),
    new SdkCard(13, 0, level, 0),
    new SdkCard(13, Suit.Heart, level, 0),
    new SdkCard(13, Suit.Club, level, 0),
    new SdkCard(13, Suit.Club, level, 1),
    new SdkCard(8, 0, level, 0),
    new SdkCard(8, Suit.Heart, level, 0),
    new SdkCard(8, Suit.Club, level, 0),
    new SdkCard(3, 0, level, 0),
    new SdkCard(3, 0, level, 1),
    new SdkCard(4, 0, level, 0),
    new SdkCard(4, 0, level, 1),
    new SdkCard(6, 0, level, 0),
    new SdkCard(6, 0, level, 1),
    new SdkCard(9, Suit.Heart, level, 0),
    new SdkCard(9, Suit.Club, level, 0),
    new SdkCard(10, Suit.Heart, level, 0),
    new SdkCard(10, Suit.Club, level, 0),
    new SdkCard(15, Suit.Joker, level, 0),
    new SdkCard(15, Suit.Joker, level, 1),
    new SdkCard(16, Suit.Joker, level, 0),
    new SdkCard(16, Suit.Joker, level, 1),
  ];

  it('场景1: prev = 顺子 A-2-3-4-5', () => {
    const prev = fromCardList([
      new SdkCard(14, 0, level, 0),
      new SdkCard(2, 1, level, 0),
      new SdkCard(3, 2, level, 0),
      new SdkCard(4, 3, level, 0),
      new SdkCard(5, 0, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Straight);
    const rawRanks = result!.getCards().map(c => c.rawRank).sort((a, b) => a - b);
    expect(rawRanks).toEqual([3, 4, 5, 6, 7]);
    expect(result!.getCards().some(c => c.isWildcard())).toBe(true);
  });

  it('场景2: prev = 顺子 3-4-5-6-7', () => {
    const prev = fromCardList([
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 1, level, 0),
      new SdkCard(5, 2, level, 0),
      new SdkCard(6, 3, level, 0),
      new SdkCard(7, 0, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Straight);
    const rawRanks = result!.getCards().map(c => c.rawRank).sort((a, b) => a - b);
    expect(rawRanks).toEqual([4, 5, 6, 7, 8]);
    expect(result!.getCards().some(c => c.isWildcard())).toBe(true);
  });

  it('场景3: prev = 顺子 5-6-7-8-9', () => {
    const prev = fromCardList([
      new SdkCard(5, 0, level, 0),
      new SdkCard(6, 1, level, 0),
      new SdkCard(7, 2, level, 0),
      new SdkCard(8, 3, level, 0),
      new SdkCard(9, 0, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Straight);
    const rawRanks = result!.getCards().map(c => c.rawRank).sort((a, b) => a - b);
    expect(rawRanks).toEqual([6, 7, 8, 9, 10]);
    expect(result!.getCards().some(c => c.isWildcard())).toBe(false);
  });

  it('场景4: prev = 顺子 6-7-8-9-10', () => {
    const prev = fromCardList([
      new SdkCard(6, 0, level, 0),
      new SdkCard(7, 1, level, 0),
      new SdkCard(8, 2, level, 0),
      new SdkCard(9, 3, level, 0),
      new SdkCard(10, 0, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Straight);
    const rawRanks = result!.getCards().map(c => c.rawRank).sort((a, b) => a - b);
    expect(rawRanks).toEqual([5, 7, 8, 9, 10]);
    expect(result!.getCards().some(c => c.isWildcard())).toBe(true);
  });

  it('场景5: prev = 顺子 10-J-Q-K-A', () => {
    const prev = fromCardList([
      new SdkCard(10, 0, level, 0),
      new SdkCard(11, 1, level, 0),
      new SdkCard(12, 2, level, 0),
      new SdkCard(13, 3, level, 0),
      new SdkCard(14, 0, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
  });
});

describe('findMinTube', () => {
  const level = 5;

  it('should find smallest tube to beat prev tube', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(4, 0, level, 2),
      new SdkCard(4, 1, level, 3),
      new SdkCard(5, 0, level, 4),
      new SdkCard(5, 1, level, 5),
      new SdkCard(6, 0, level, 6),
      new SdkCard(6, 1, level, 7),
      new SdkCard(7, 0, level, 8),
      new SdkCard(7, 1, level, 9),
    ];

    const prev = fromCardList([
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(3, 2, level, 12),
      new SdkCard(3, 3, level, 13),
      new SdkCard(4, 0, level, 14),
      new SdkCard(4, 1, level, 15),
    ]);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Tube);
  });

  it('should find tube with wildcard', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(4, 0, level, 2),
      new SdkCard(4, 1, level, 3),
      new SdkCard(level, Suit.Heart, level, 4),
      new SdkCard(level, Suit.Heart, level, 5),
    ];

    const prev = fromCardList([
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(3, 2, level, 12),
      new SdkCard(3, 3, level, 13),
      new SdkCard(4, 0, level, 14),
      new SdkCard(4, 1, level, 15),
    ]);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Tube);
  });

  it('should return null when no tube can beat prev', () => {
    const handCards = [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(4, 0, level, 2),
      new SdkCard(4, 1, level, 3),
      new SdkCard(5, 2, level, 4),
      new SdkCard(5, 3, level, 5),
    ];

    const prev = fromCardList([
      new SdkCard(5, 0, level, 10),
      new SdkCard(5, 1, level, 11),
      new SdkCard(6, 2, level, 12),
      new SdkCard(6, 3, level, 13),
      new SdkCard(7, 0, level, 14),
      new SdkCard(7, 1, level, 15),
    ]);

    const result = findMinPlay(handCards, prev);
    expect(result).toBeNull();
  });

  it('should find tube QKA', () => {
    const handCards = [
      new SdkCard(12, 0, level, 0),
      new SdkCard(12, 1, level, 1),
      new SdkCard(13, 0, level, 2),
      new SdkCard(13, 1, level, 3),
      new SdkCard(14, 0, level, 4),
      new SdkCard(14, 1, level, 5),
    ];

    const prev = fromCardList([
      new SdkCard(10, 0, level, 10),
      new SdkCard(10, 1, level, 11),
      new SdkCard(11, 2, level, 12),
      new SdkCard(11, 3, level, 13),
      new SdkCard(12, 0, level, 14),
      new SdkCard(12, 1, level, 15),
    ]);

    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.Tube);
  });
});

describe('findMinSameType - FullHouse (27 cards)', () => {
  const level = 5;

  const handCards = [
    new SdkCard(5, Suit.Heart, level, 0),
    new SdkCard(5, Suit.Heart, level, 1),
    new SdkCard(7, 0, level, 0),
    new SdkCard(7, 0, level, 1),
    new SdkCard(7, Suit.Heart, level, 0),
    new SdkCard(7, Suit.Club, level, 0),
    new SdkCard(13, 0, level, 0),
    new SdkCard(13, Suit.Heart, level, 0),
    new SdkCard(13, Suit.Club, level, 0),
    new SdkCard(13, Suit.Club, level, 1),
    new SdkCard(8, 0, level, 0),
    new SdkCard(8, Suit.Heart, level, 0),
    new SdkCard(8, Suit.Club, level, 0),
    new SdkCard(3, 0, level, 0),
    new SdkCard(3, 0, level, 1),
    new SdkCard(4, 0, level, 0),
    new SdkCard(4, 0, level, 1),
    new SdkCard(6, 0, level, 0),
    new SdkCard(6, 0, level, 1),
    new SdkCard(9, Suit.Heart, level, 0),
    new SdkCard(9, Suit.Club, level, 0),
    new SdkCard(10, Suit.Heart, level, 0),
    new SdkCard(10, Suit.Club, level, 0),
    new SdkCard(15, Suit.Joker, level, 0),
    new SdkCard(15, Suit.Joker, level, 1),
    new SdkCard(16, Suit.Joker, level, 0),
    new SdkCard(16, Suit.Joker, level, 1),
  ];

  it('场景1: prev = 333+44 (基础三带二)', () => {
    const prev = fromCardList([
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
      new SdkCard(3, 2, level, 0),
      new SdkCard(4, 0, level, 0),
      new SdkCard(4, 1, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.FullHouse);
    expect(result!.greaterThan(prev)).toBe(true);
  });

  it('场景2: prev = 666+33 (测试选最小)', () => {
    const prev = fromCardList([
      new SdkCard(6, 0, level, 0),
      new SdkCard(6, 1, level, 0),
      new SdkCard(6, 2, level, 0),
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.FullHouse);
    expect(result!.greaterThan(prev)).toBe(true);
  });

  it('场景3: prev = 777+33 (需要用888)', () => {
    const prev = fromCardList([
      new SdkCard(7, 0, level, 0),
      new SdkCard(7, 1, level, 0),
      new SdkCard(7, 2, level, 0),
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.FullHouse);
    expect(result!.greaterThan(prev)).toBe(true);
  });

  it('场景4: prev = 888+44 (需要百搭凑三张)', () => {
    const prev = fromCardList([
      new SdkCard(8, 0, level, 0),
      new SdkCard(8, 1, level, 0),
      new SdkCard(8, 2, level, 0),
      new SdkCard(4, 0, level, 0),
      new SdkCard(4, 1, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.FullHouse);
    expect(result!.greaterThan(prev)).toBe(true);
  });

  it('场景5: prev = 777+AA (测试用百搭或王对)', () => {
    const prev = fromCardList([
      new SdkCard(7, 0, level, 0),
      new SdkCard(7, 1, level, 0),
      new SdkCard(7, 2, level, 0),
      new SdkCard(14, 0, level, 0),
      new SdkCard(14, 1, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.FullHouse);
    expect(result!.greaterThan(prev)).toBe(true);
  });

  it('场景6: prev = 888+大王对 (测试更大三张)', () => {
    const prev = fromCardList([
      new SdkCard(8, 0, level, 0),
      new SdkCard(8, 1, level, 0),
      new SdkCard(8, 2, level, 0),
      new SdkCard(16, Suit.Joker, level, 0),
      new SdkCard(16, Suit.Joker, level, 1),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.FullHouse);
    expect(result!.greaterThan(prev)).toBe(true);
  });

  it('场景7: prev = KKK+大王对 (无法压过用炸弹)', () => {
    const prev = fromCardList([
      new SdkCard(13, 0, level, 0),
      new SdkCard(13, 1, level, 0),
      new SdkCard(13, 2, level, 0),
      new SdkCard(16, Suit.Joker, level, 0),
      new SdkCard(16, Suit.Joker, level, 1),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.NaiveBomb);
  });

  it('场景8: prev = 999+1010 (边界测试)', () => {
    const prev = fromCardList([
      new SdkCard(9, 0, level, 0),
      new SdkCard(9, 1, level, 0),
      new SdkCard(9, 2, level, 0),
      new SdkCard(10, 0, level, 0),
      new SdkCard(10, 1, level, 0),
    ]);
    const result = findMinPlay(handCards, prev);
    expect(result).not.toBeNull();
    expect(result!.getType()).toBe(CompType.FullHouse);
  });
});
