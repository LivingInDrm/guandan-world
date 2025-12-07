import { SdkCard } from '../card';
import {
  STRAIGHT_CARD_COUNT,
  sortCards,
  sortCardsForConsecutive,
  hasJokers,
  countWildcards,
} from '../compUtil';

const STRAIGHT_CARDS_PER_NUMBER = 1;

const allConsecutiveFives: number[][] = [
  [1, 2, 3, 4, 5],
  [2, 3, 4, 5, 6],
  [3, 4, 5, 6, 7],
  [4, 5, 6, 7, 8],
  [5, 6, 7, 8, 9],
  [6, 7, 8, 9, 10],
  [7, 8, 9, 10, 11],
  [8, 9, 10, 11, 12],
  [9, 10, 11, 12, 13],
  [10, 11, 12, 13, 1],
];

function canFormStraightWithWildcards(
  five: number[],
  cardCounts: Map<number, number>,
  wildcardCount: number
): boolean {
  let needed = 0;
  for (const num of five) {
    const have = cardCounts.get(num) ?? 0;
    if (have < STRAIGHT_CARDS_PER_NUMBER) {
      needed += STRAIGHT_CARDS_PER_NUMBER - have;
    }
  }
  return needed === wildcardCount;
}

function findAllValidStraightFives(
  cardCounts: Map<number, number>,
  wildcardCount: number
): number[][] {
  const valid: number[][] = [];
  for (const five of allConsecutiveFives) {
    if (canFormStraightWithWildcards(five, cardCounts, wildcardCount)) {
      valid.push(five);
    }
  }
  return valid;
}

function getStraightFiveComparisonKey(five: number[]): number {
  return five[0];
}

function selectBestStraightFive(validFives: number[][]): number[] | null {
  if (validFives.length === 0) return null;

  let bestFive = validFives[0];
  let bestKey = getStraightFiveComparisonKey(bestFive);

  for (let i = 1; i < validFives.length; i++) {
    const currentKey = getStraightFiveComparisonKey(validFives[i]);
    if (currentKey > bestKey) {
      bestKey = currentKey;
      bestFive = validFives[i];
    }
  }

  return bestFive;
}

function normalizeStraightWithFive(sortedCards: SdkCard[], bestFive: number[]): SdkCard[] {
  const result: SdkCard[] = [];
  const wildcardPool: SdkCard[] = [];

  const cardsByNum = new Map<number, SdkCard[]>();
  for (const card of sortedCards) {
    if (card.isWildcard()) {
      wildcardPool.push(card);
    } else {
      const arr = cardsByNum.get(card.rawRank) ?? [];
      arr.push(card);
      cardsByNum.set(card.rawRank, arr);
    }
  }

  for (const num of bestFive) {
    const cardsOfNum = cardsByNum.get(num) ?? [];
    if (cardsOfNum.length > 0) {
      result.push(cardsOfNum[0]);
    } else if (wildcardPool.length > 0) {
      result.push(wildcardPool.shift()!);
    }
  }

  if (result.length !== STRAIGHT_CARD_COUNT) {
    return sortedCards;
  }

  return result;
}

export interface StraightSatisfyResult {
  isValid: boolean;
  normalizedCards: SdkCard[];
  comparisonKey: number;
}

export function straightSatisfy(cards: SdkCard[]): StraightSatisfyResult {
  if (!cards || cards.length !== STRAIGHT_CARD_COUNT) {
    return { isValid: false, normalizedCards: sortCards(cards ?? []), comparisonKey: 0 };
  }

  const sortedCards = sortCardsForConsecutive(cards);

  if (hasJokers(sortedCards)) {
    return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
  }

  const wildcardCount = countWildcards(sortedCards);

  const cardCounts = new Map<number, number>();
  for (const card of sortedCards) {
    if (!card.isWildcard()) {
      cardCounts.set(card.rawRank, (cardCounts.get(card.rawRank) ?? 0) + 1);
    }
  }

  for (const count of cardCounts.values()) {
    if (count > STRAIGHT_CARDS_PER_NUMBER) {
      return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
    }
  }

  const validFives = findAllValidStraightFives(cardCounts, wildcardCount);

  if (validFives.length === 0) {
    return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
  }

  const bestFive = selectBestStraightFive(validFives);

  if (!bestFive) {
    return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
  }

  const comparisonKey = getStraightFiveComparisonKey(bestFive);
  const normalizedCards = normalizeStraightWithFive(sortedCards, bestFive);

  return { isValid: true, normalizedCards, comparisonKey };
}
