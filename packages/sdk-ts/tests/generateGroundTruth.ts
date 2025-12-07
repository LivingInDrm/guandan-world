import { SdkCard } from '../src/card';
import { Suit } from '../src/types';
import { findMinPlay } from '../src/findMinPlay';
import { fromCardList } from '../src/fromCardList';
import { cardsToShortStrings, createStandardHandCards } from './testUtils';

interface TestCase {
  name: string;
  handCards: SdkCard[];
  prevCards: SdkCard[];
}

const level = 5;

const testCases: TestCase[] = [
  {
    name: 'findMinPlay: should return null when no valid play exists',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 0, level, 1),
      new SdkCard(6, 0, level, 2),
    ],
    prevCards: [new SdkCard(7, 0, level, 10)],
  },
  {
    name: 'findMinBomb: should find 4-card NaiveBomb to beat a Single',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(3, 3, level, 3),
    ],
    prevCards: [new SdkCard(14, 0, level, 10)],
  },
  {
    name: 'findMinBomb: should find the smallest NaiveBomb among multiple options',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(3, 3, level, 3),
      new SdkCard(7, 0, level, 4),
      new SdkCard(7, 1, level, 5),
      new SdkCard(7, 2, level, 6),
      new SdkCard(7, 3, level, 7),
    ],
    prevCards: [new SdkCard(14, 0, level, 10)],
  },
  {
    name: 'findMinBomb: should find bigger NaiveBomb to beat smaller NaiveBomb',
    handCards: [
      new SdkCard(7, 0, level, 0),
      new SdkCard(7, 1, level, 1),
      new SdkCard(7, 2, level, 2),
      new SdkCard(7, 3, level, 3),
    ],
    prevCards: [
      new SdkCard(3, 0, level, 10),
      new SdkCard(3, 1, level, 11),
      new SdkCard(3, 2, level, 12),
      new SdkCard(3, 3, level, 13),
    ],
  },
  {
    name: 'findMinBomb: should return null when no bomb can beat prev bomb',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(3, 3, level, 3),
    ],
    prevCards: [
      new SdkCard(7, 0, level, 10),
      new SdkCard(7, 1, level, 11),
      new SdkCard(7, 2, level, 12),
      new SdkCard(7, 3, level, 13),
    ],
  },
  {
    name: 'findMinBomb: should find StraightFlush to beat a smaller bomb',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 0, level, 1),
      new SdkCard(5, 0, level, 2),
      new SdkCard(6, 0, level, 3),
      new SdkCard(7, 0, level, 4),
    ],
    prevCards: [
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(2, 2, level, 12),
      new SdkCard(2, 3, level, 13),
    ],
  },
  {
    name: 'findMinBomb: should find JokerBomb to beat a big bomb',
    handCards: [
      new SdkCard(15, Suit.Joker, level, 0),
      new SdkCard(15, Suit.Joker, level, 1),
      new SdkCard(16, Suit.Joker, level, 2),
      new SdkCard(16, Suit.Joker, level, 3),
    ],
    prevCards: [
      new SdkCard(14, 0, level, 0),
      new SdkCard(14, 1, level, 0),
      new SdkCard(14, 2, level, 0),
      new SdkCard(14, 3, level, 0),
      new SdkCard(14, 0, level, 1),
      new SdkCard(14, 1, level, 1),
      new SdkCard(14, 2, level, 1),
      new SdkCard(14, 3, level, 1),
    ],
  },
  {
    name: 'findMinBomb: should find NaiveBomb with wildcard to beat a smaller bomb',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(level, Suit.Heart, level, 3),
    ],
    prevCards: [
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(2, 2, level, 12),
      new SdkCard(2, 3, level, 13),
    ],
  },
  {
    name: 'findMinBomb: should prefer smaller bomb over larger bomb when beating a bomb',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(3, 3, level, 3),
      new SdkCard(7, 0, level, 4),
      new SdkCard(7, 1, level, 5),
      new SdkCard(7, 2, level, 6),
      new SdkCard(7, 3, level, 7),
      new SdkCard(level, Suit.Heart, level, 8),
    ],
    prevCards: [
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(2, 2, level, 12),
      new SdkCard(2, 3, level, 13),
    ],
  },
  {
    name: 'findMinBomb 27cards: 场景1 prev = 4个9',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(9, 0, level, 0),
      new SdkCard(9, Suit.Heart, level, 0),
      new SdkCard(9, Suit.Club, level, 0),
      new SdkCard(9, Suit.Diamond, level, 0),
    ],
  },
  {
    name: 'findMinBomb 27cards: 场景2 prev = 5个10',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(10, 0, level, 0),
      new SdkCard(10, 0, level, 1),
      new SdkCard(10, Suit.Heart, level, 0),
      new SdkCard(10, Suit.Club, level, 0),
      new SdkCard(10, Suit.Diamond, level, 0),
    ],
  },
  {
    name: 'findMinBomb 27cards: 场景3 prev = 6个5',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(5, 0, level, 0),
      new SdkCard(5, 0, level, 1),
      new SdkCard(5, Suit.Club, level, 0),
      new SdkCard(5, Suit.Club, level, 1),
      new SdkCard(5, Suit.Diamond, level, 0),
      new SdkCard(5, Suit.Diamond, level, 1),
    ],
  },
  {
    name: 'findMinBomb 27cards: 场景4 prev = 7个6',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(6, 0, level, 0),
      new SdkCard(6, 0, level, 1),
      new SdkCard(6, Suit.Heart, level, 0),
      new SdkCard(6, Suit.Heart, level, 1),
      new SdkCard(6, Suit.Club, level, 0),
      new SdkCard(6, Suit.Club, level, 1),
      new SdkCard(6, Suit.Diamond, level, 0),
    ],
  },
  {
    name: 'findMinBomb 27cards: 场景5 prev = 8个2',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(2, 0, level, 0),
      new SdkCard(2, 0, level, 1),
      new SdkCard(2, Suit.Heart, level, 0),
      new SdkCard(2, Suit.Heart, level, 1),
      new SdkCard(2, Suit.Club, level, 0),
      new SdkCard(2, Suit.Club, level, 1),
      new SdkCard(2, Suit.Diamond, level, 0),
      new SdkCard(2, Suit.Diamond, level, 1),
    ],
  },
  {
    name: 'findMinBomb 27cards: 场景6 prev = Heart 3,4,5,6,7',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(3, Suit.Heart, level, 0),
      new SdkCard(4, Suit.Heart, level, 0),
      new SdkCard(5, Suit.Heart, level, 0),
      new SdkCard(6, Suit.Heart, level, 0),
      new SdkCard(7, Suit.Heart, level, 0),
    ],
  },
  {
    name: 'findMinBomb 27cards: 场景7 prev = Club 5,6,7,8,9',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(5, Suit.Club, level, 0),
      new SdkCard(6, Suit.Club, level, 0),
      new SdkCard(7, Suit.Club, level, 0),
      new SdkCard(8, Suit.Club, level, 0),
      new SdkCard(9, Suit.Club, level, 0),
    ],
  },
  {
    name: 'Single: should find smallest single to beat prev single',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(7, 0, level, 1),
      new SdkCard(10, 0, level, 2),
    ],
    prevCards: [new SdkCard(6, 0, level, 10)],
  },
  {
    name: 'Single: should return null when no single can beat prev',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 0, level, 1),
    ],
    prevCards: [new SdkCard(10, 0, level, 10)],
  },
  {
    name: 'Single: should use wildcard as single',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(level, Suit.Heart, level, 1),
    ],
    prevCards: [new SdkCard(10, 0, level, 10)],
  },
  {
    name: 'Pair: should find smallest pair to beat prev pair',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(7, 0, level, 2),
      new SdkCard(7, 1, level, 3),
      new SdkCard(10, 0, level, 4),
      new SdkCard(10, 1, level, 5),
    ],
    prevCards: [
      new SdkCard(6, 0, level, 10),
      new SdkCard(6, 1, level, 11),
    ],
  },
  {
    name: 'Pair: should find pair with wildcard',
    handCards: [
      new SdkCard(7, 0, level, 0),
      new SdkCard(level, Suit.Heart, level, 1),
    ],
    prevCards: [
      new SdkCard(6, 0, level, 10),
      new SdkCard(6, 1, level, 11),
    ],
  },
  {
    name: 'Pair: should return null when no pair can beat prev',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
    ],
    prevCards: [
      new SdkCard(10, 0, level, 10),
      new SdkCard(10, 1, level, 11),
    ],
  },
  {
    name: 'Triple: should find smallest triple to beat prev triple',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(3, 2, level, 2),
      new SdkCard(8, 0, level, 3),
      new SdkCard(8, 1, level, 4),
      new SdkCard(8, 2, level, 5),
    ],
    prevCards: [
      new SdkCard(6, 0, level, 10),
      new SdkCard(6, 1, level, 11),
      new SdkCard(6, 2, level, 12),
    ],
  },
  {
    name: 'Triple: should find triple with wildcard',
    handCards: [
      new SdkCard(8, 0, level, 0),
      new SdkCard(8, 1, level, 1),
      new SdkCard(level, Suit.Heart, level, 2),
    ],
    prevCards: [
      new SdkCard(6, 0, level, 10),
      new SdkCard(6, 1, level, 11),
      new SdkCard(6, 2, level, 12),
    ],
  },
  {
    name: 'Triple: should use pure wildcards as triple',
    handCards: [
      new SdkCard(level, Suit.Heart, level, 0),
      new SdkCard(level, Suit.Heart, level, 1),
      new SdkCard(level, Suit.Heart, level, 2),
    ],
    prevCards: [
      new SdkCard(10, 0, level, 10),
      new SdkCard(10, 1, level, 11),
      new SdkCard(10, 2, level, 12),
    ],
  },
  {
    name: 'Single 27cards: 场景1 prev = 单张黑桃3',
    handCards: createStandardHandCards(level),
    prevCards: [new SdkCard(3, 0, level, 0)],
  },
  {
    name: 'Single 27cards: 场景2 prev = 单张黑桃K',
    handCards: createStandardHandCards(level),
    prevCards: [new SdkCard(13, 0, level, 0)],
  },
  {
    name: 'Single 27cards: 场景3 prev = 单张A',
    handCards: createStandardHandCards(level),
    prevCards: [new SdkCard(14, 0, level, 0)],
  },
  {
    name: 'Single 27cards: 场景4 prev = 单张小王',
    handCards: createStandardHandCards(level),
    prevCards: [new SdkCard(15, Suit.Joker, level, 0)],
  },
  {
    name: 'Single 27cards: 场景5 prev = 单张大王',
    handCards: createStandardHandCards(level),
    prevCards: [new SdkCard(16, Suit.Joker, level, 0)],
  },
  {
    name: 'Single 27cards: 场景6 prev = 单张5 (级牌)',
    handCards: createStandardHandCards(level),
    prevCards: [new SdkCard(5, 0, level, 0)],
  },
  {
    name: 'Pair 27cards: 场景1 prev = 对3',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
    ],
  },
  {
    name: 'Pair 27cards: 场景2 prev = 对8',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(8, 0, level, 0),
      new SdkCard(8, 1, level, 0),
    ],
  },
  {
    name: 'Pair 27cards: 场景3 prev = 对K',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(13, 0, level, 0),
      new SdkCard(13, 1, level, 0),
    ],
  },
  {
    name: 'Pair 27cards: 场景4 prev = 对小王',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(15, Suit.Joker, level, 0),
      new SdkCard(15, Suit.Joker, level, 1),
    ],
  },
  {
    name: 'Pair 27cards: 场景5 prev = 对大王',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(16, Suit.Joker, level, 0),
      new SdkCard(16, Suit.Joker, level, 1),
    ],
  },
  {
    name: 'Pair 27cards: 场景6 prev = 对5 (级牌)',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(5, 0, level, 0),
      new SdkCard(5, Suit.Club, level, 0),
    ],
  },
  {
    name: 'Triple 27cards: 场景1 prev = 三张3',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
      new SdkCard(3, 2, level, 0),
    ],
  },
  {
    name: 'Triple 27cards: 场景2 prev = 三张6',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(6, 0, level, 0),
      new SdkCard(6, 1, level, 0),
      new SdkCard(6, 2, level, 0),
    ],
  },
  {
    name: 'Triple 27cards: 场景3 prev = 三张K',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(13, 0, level, 0),
      new SdkCard(13, 1, level, 0),
      new SdkCard(13, 2, level, 0),
    ],
  },
  {
    name: 'Triple 27cards: 场景4 prev = 三张8',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(8, 0, level, 0),
      new SdkCard(8, 1, level, 0),
      new SdkCard(8, 2, level, 0),
    ],
  },
  {
    name: 'Straight: should find smallest straight to beat prev straight',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 0, level, 1),
      new SdkCard(5, 0, level, 2),
      new SdkCard(6, 0, level, 3),
      new SdkCard(7, 0, level, 4),
      new SdkCard(8, 0, level, 5),
      new SdkCard(9, 0, level, 6),
      new SdkCard(10, 0, level, 7),
      new SdkCard(11, 0, level, 8),
    ],
    prevCards: [
      new SdkCard(2, 0, level, 10),
      new SdkCard(3, 1, level, 11),
      new SdkCard(4, 2, level, 12),
      new SdkCard(5, 3, level, 13),
      new SdkCard(6, 0, level, 14),
    ],
  },
  {
    name: 'Straight: should find straight with wildcard',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 0, level, 1),
      new SdkCard(6, 0, level, 3),
      new SdkCard(7, 0, level, 4),
      new SdkCard(level, Suit.Heart, level, 5),
    ],
    prevCards: [
      new SdkCard(2, 0, level, 10),
      new SdkCard(3, 1, level, 11),
      new SdkCard(4, 2, level, 12),
      new SdkCard(5, 3, level, 13),
      new SdkCard(6, 0, level, 14),
    ],
  },
  {
    name: 'Straight: should return null when no straight can beat prev',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 1, level, 1),
      new SdkCard(5, 2, level, 2),
      new SdkCard(6, 3, level, 3),
      new SdkCard(7, 0, level, 4),
    ],
    prevCards: [
      new SdkCard(6, 0, level, 10),
      new SdkCard(7, 1, level, 11),
      new SdkCard(8, 2, level, 12),
      new SdkCard(9, 3, level, 13),
      new SdkCard(10, 0, level, 14),
    ],
  },
  {
    name: 'Straight: should find straight 10-A',
    handCards: [
      new SdkCard(10, 0, level, 0),
      new SdkCard(11, 0, level, 1),
      new SdkCard(12, 0, level, 2),
      new SdkCard(13, 0, level, 3),
      new SdkCard(14, 0, level, 4),
    ],
    prevCards: [
      new SdkCard(8, 0, level, 10),
      new SdkCard(9, 1, level, 11),
      new SdkCard(10, 2, level, 12),
      new SdkCard(11, 3, level, 13),
      new SdkCard(12, 0, level, 14),
    ],
  },
  {
    name: 'Straight 27cards: 场景1 prev = 顺子 A-2-3-4-5',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(14, 0, level, 0),
      new SdkCard(2, 1, level, 0),
      new SdkCard(3, 2, level, 0),
      new SdkCard(4, 3, level, 0),
      new SdkCard(5, 0, level, 0),
    ],
  },
  {
    name: 'Straight 27cards: 场景2 prev = 顺子 3-4-5-6-7',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(4, 1, level, 0),
      new SdkCard(5, 2, level, 0),
      new SdkCard(6, 3, level, 0),
      new SdkCard(7, 0, level, 0),
    ],
  },
  {
    name: 'Straight 27cards: 场景3 prev = 顺子 5-6-7-8-9',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(5, 0, level, 0),
      new SdkCard(6, 1, level, 0),
      new SdkCard(7, 2, level, 0),
      new SdkCard(8, 3, level, 0),
      new SdkCard(9, 0, level, 0),
    ],
  },
  {
    name: 'Straight 27cards: 场景4 prev = 顺子 6-7-8-9-10',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(6, 0, level, 0),
      new SdkCard(7, 1, level, 0),
      new SdkCard(8, 2, level, 0),
      new SdkCard(9, 3, level, 0),
      new SdkCard(10, 0, level, 0),
    ],
  },
  {
    name: 'Straight 27cards: 场景5 prev = 顺子 10-J-Q-K-A',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(10, 0, level, 0),
      new SdkCard(11, 1, level, 0),
      new SdkCard(12, 2, level, 0),
      new SdkCard(13, 3, level, 0),
      new SdkCard(14, 0, level, 0),
    ],
  },
  {
    name: 'Tube: should find smallest tube to beat prev tube',
    handCards: [
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
    ],
    prevCards: [
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(3, 2, level, 12),
      new SdkCard(3, 3, level, 13),
      new SdkCard(4, 0, level, 14),
      new SdkCard(4, 1, level, 15),
    ],
  },
  {
    name: 'Tube: should find tube with wildcard',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(4, 0, level, 2),
      new SdkCard(4, 1, level, 3),
      new SdkCard(level, Suit.Heart, level, 4),
      new SdkCard(level, Suit.Heart, level, 5),
    ],
    prevCards: [
      new SdkCard(2, 0, level, 10),
      new SdkCard(2, 1, level, 11),
      new SdkCard(3, 2, level, 12),
      new SdkCard(3, 3, level, 13),
      new SdkCard(4, 0, level, 14),
      new SdkCard(4, 1, level, 15),
    ],
  },
  {
    name: 'Tube: should return null when no tube can beat prev',
    handCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 1),
      new SdkCard(4, 0, level, 2),
      new SdkCard(4, 1, level, 3),
      new SdkCard(5, 2, level, 4),
      new SdkCard(5, 3, level, 5),
    ],
    prevCards: [
      new SdkCard(5, 0, level, 10),
      new SdkCard(5, 1, level, 11),
      new SdkCard(6, 2, level, 12),
      new SdkCard(6, 3, level, 13),
      new SdkCard(7, 0, level, 14),
      new SdkCard(7, 1, level, 15),
    ],
  },
  {
    name: 'Tube: should find tube QKA',
    handCards: [
      new SdkCard(12, 0, level, 0),
      new SdkCard(12, 1, level, 1),
      new SdkCard(13, 0, level, 2),
      new SdkCard(13, 1, level, 3),
      new SdkCard(14, 0, level, 4),
      new SdkCard(14, 1, level, 5),
    ],
    prevCards: [
      new SdkCard(10, 0, level, 10),
      new SdkCard(10, 1, level, 11),
      new SdkCard(11, 2, level, 12),
      new SdkCard(11, 3, level, 13),
      new SdkCard(12, 0, level, 14),
      new SdkCard(12, 1, level, 15),
    ],
  },
  {
    name: 'FullHouse 27cards: 场景1 prev = 333+44',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
      new SdkCard(3, 2, level, 0),
      new SdkCard(4, 0, level, 0),
      new SdkCard(4, 1, level, 0),
    ],
  },
  {
    name: 'FullHouse 27cards: 场景2 prev = 666+33',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(6, 0, level, 0),
      new SdkCard(6, 1, level, 0),
      new SdkCard(6, 2, level, 0),
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
    ],
  },
  {
    name: 'FullHouse 27cards: 场景3 prev = 777+33',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(7, 0, level, 0),
      new SdkCard(7, 1, level, 0),
      new SdkCard(7, 2, level, 0),
      new SdkCard(3, 0, level, 0),
      new SdkCard(3, 1, level, 0),
    ],
  },
  {
    name: 'FullHouse 27cards: 场景4 prev = 888+44',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(8, 0, level, 0),
      new SdkCard(8, 1, level, 0),
      new SdkCard(8, 2, level, 0),
      new SdkCard(4, 0, level, 0),
      new SdkCard(4, 1, level, 0),
    ],
  },
  {
    name: 'FullHouse 27cards: 场景5 prev = 777+AA',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(7, 0, level, 0),
      new SdkCard(7, 1, level, 0),
      new SdkCard(7, 2, level, 0),
      new SdkCard(14, 0, level, 0),
      new SdkCard(14, 1, level, 0),
    ],
  },
  {
    name: 'FullHouse 27cards: 场景6 prev = 888+大王对',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(8, 0, level, 0),
      new SdkCard(8, 1, level, 0),
      new SdkCard(8, 2, level, 0),
      new SdkCard(16, Suit.Joker, level, 0),
      new SdkCard(16, Suit.Joker, level, 1),
    ],
  },
  {
    name: 'FullHouse 27cards: 场景7 prev = KKK+大王对',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(13, 0, level, 0),
      new SdkCard(13, 1, level, 0),
      new SdkCard(13, 2, level, 0),
      new SdkCard(16, Suit.Joker, level, 0),
      new SdkCard(16, Suit.Joker, level, 1),
    ],
  },
  {
    name: 'FullHouse 27cards: 场景8 prev = 999+1010',
    handCards: createStandardHandCards(level),
    prevCards: [
      new SdkCard(9, 0, level, 0),
      new SdkCard(9, 1, level, 0),
      new SdkCard(9, 2, level, 0),
      new SdkCard(10, 0, level, 0),
      new SdkCard(10, 1, level, 0),
    ],
  },
];

console.log('// Ground Truth for findMinPlay tests');
console.log('// Generated at:', new Date().toISOString());
console.log('export const groundTruth: Record<string, string[] | null> = {');

for (const tc of testCases) {
  const prev = fromCardList(tc.prevCards);
  const result = findMinPlay(tc.handCards, prev);
  const output = cardsToShortStrings(result);
  const jsonOutput = output === null ? 'null' : JSON.stringify(output);
  console.log(`  ${JSON.stringify(tc.name)}: ${jsonOutput},`);
}

console.log('};');
