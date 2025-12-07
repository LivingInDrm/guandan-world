import { SdkCard } from '../card';
import {
  PLATE_CARD_COUNT,
  sortCards,
  sortCardsForConsecutive,
  hasJokers,
  countWildcards,
} from '../compUtil';

const PLATE_CARDS_PER_NUMBER = 3;

const allConsecutivePairs: number[][] = [
  [1, 2],
  [2, 3],
  [3, 4],
  [4, 5],
  [5, 6],
  [6, 7],
  [7, 8],
  [8, 9],
  [9, 10],
  [10, 11],
  [11, 12],
  [12, 13],
  [13, 1],
];

function canFormPlateWithWildcards(
  pair: number[],
  cardCounts: Map<number, number>,
  wildcardCount: number
): boolean {
  let needed = 0;
  for (const num of pair) {
    const have = cardCounts.get(num) ?? 0;
    if (have < PLATE_CARDS_PER_NUMBER) {
      needed += PLATE_CARDS_PER_NUMBER - have;
    }
  }
  return needed === wildcardCount;
}

function findAllValidPlatePairs(
  cardCounts: Map<number, number>,
  wildcardCount: number
): number[][] {
  const valid: number[][] = [];
  for (const pair of allConsecutivePairs) {
    if (canFormPlateWithWildcards(pair, cardCounts, wildcardCount)) {
      valid.push(pair);
    }
  }
  return valid;
}

function getPlatePairComparisonKey(pair: number[]): number {
  return pair[0];
}

function selectBestPlatePair(validPairs: number[][]): number[] | null {
  if (validPairs.length === 0) return null;

  let bestPair = validPairs[0];
  let bestKey = getPlatePairComparisonKey(bestPair);

  for (let i = 1; i < validPairs.length; i++) {
    const currentKey = getPlatePairComparisonKey(validPairs[i]);
    if (currentKey > bestKey) {
      bestKey = currentKey;
      bestPair = validPairs[i];
    }
  }

  return bestPair;
}

function normalizePlateWithPair(sortedCards: SdkCard[], bestPair: number[]): SdkCard[] {
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

  for (const num of bestPair) {
    const cardsOfNum = cardsByNum.get(num) ?? [];

    for (let i = 0; i < PLATE_CARDS_PER_NUMBER; i++) {
      if (i < cardsOfNum.length) {
        result.push(cardsOfNum[i]);
      } else if (wildcardPool.length > 0) {
        result.push(wildcardPool.shift()!);
      }
    }
  }

  if (result.length !== PLATE_CARD_COUNT) {
    return sortedCards;
  }

  return result;
}

export interface PlateSatisfyResult {
  isValid: boolean;
  normalizedCards: SdkCard[];
  comparisonKey: number;
}

export function plateSatisfy(cards: SdkCard[]): PlateSatisfyResult {
  if (!cards || cards.length !== PLATE_CARD_COUNT) {
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
    if (count > PLATE_CARDS_PER_NUMBER) {
      return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
    }
  }

  const validPairs = findAllValidPlatePairs(cardCounts, wildcardCount);

  if (validPairs.length === 0) {
    return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
  }

  const bestPair = selectBestPlatePair(validPairs);

  if (!bestPair) {
    return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
  }

  const comparisonKey = getPlatePairComparisonKey(bestPair);
  const normalizedCards = normalizePlateWithPair(sortedCards, bestPair);

  return { isValid: true, normalizedCards, comparisonKey };
}
