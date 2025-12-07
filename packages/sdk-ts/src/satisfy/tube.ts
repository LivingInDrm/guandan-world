import { SdkCard } from '../card';
import {
  TUBE_CARD_COUNT,
  sortCards,
  sortCardsForConsecutive,
  hasJokers,
  countWildcards,
} from '../compUtil';

const TUBE_CARDS_PER_NUMBER = 2;

const allConsecutiveTriples: number[][] = [
  [2, 3, 4],
  [3, 4, 5],
  [4, 5, 6],
  [5, 6, 7],
  [6, 7, 8],
  [7, 8, 9],
  [8, 9, 10],
  [9, 10, 11],
  [10, 11, 12],
  [11, 12, 13],
  [1, 2, 3],
  [12, 13, 1],
];

function canFormTubeWithWildcards(
  triple: number[],
  cardCounts: Map<number, number>,
  wildcardCount: number
): boolean {
  let needed = 0;
  for (const num of triple) {
    const have = cardCounts.get(num) ?? 0;
    if (have < TUBE_CARDS_PER_NUMBER) {
      needed += TUBE_CARDS_PER_NUMBER - have;
    }
  }
  return needed === wildcardCount;
}

function findAllValidTubeTriples(
  cardCounts: Map<number, number>,
  wildcardCount: number
): number[][] {
  const valid: number[][] = [];
  for (const triple of allConsecutiveTriples) {
    if (canFormTubeWithWildcards(triple, cardCounts, wildcardCount)) {
      valid.push(triple);
    }
  }
  return valid;
}

function getTubeTripleComparisonKey(triple: number[]): number {
  return triple[0];
}

function selectBestTubeTriple(validTriples: number[][]): number[] | null {
  if (validTriples.length === 0) return null;

  let bestTriple = validTriples[0];
  let bestKey = getTubeTripleComparisonKey(bestTriple);

  for (let i = 1; i < validTriples.length; i++) {
    const currentKey = getTubeTripleComparisonKey(validTriples[i]);
    if (currentKey > bestKey) {
      bestKey = currentKey;
      bestTriple = validTriples[i];
    }
  }

  return bestTriple;
}

function normalizeTubeWithTriple(sortedCards: SdkCard[], bestTriple: number[]): SdkCard[] {
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

  for (const num of bestTriple) {
    const cardsOfNum = cardsByNum.get(num) ?? [];

    for (let i = 0; i < TUBE_CARDS_PER_NUMBER; i++) {
      if (i < cardsOfNum.length) {
        result.push(cardsOfNum[i]);
      } else if (wildcardPool.length > 0) {
        result.push(wildcardPool.shift()!);
      }
    }
  }

  if (result.length !== TUBE_CARD_COUNT) {
    return sortedCards;
  }

  return result;
}

export interface TubeSatisfyResult {
  isValid: boolean;
  normalizedCards: SdkCard[];
  comparisonKey: number;
}

export function tubeSatisfy(cards: SdkCard[]): TubeSatisfyResult {
  if (!cards || cards.length !== TUBE_CARD_COUNT) {
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
    if (count > TUBE_CARDS_PER_NUMBER) {
      return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
    }
  }

  const validTriples = findAllValidTubeTriples(cardCounts, wildcardCount);

  if (validTriples.length === 0) {
    return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
  }

  const bestTriple = selectBestTubeTriple(validTriples);

  if (!bestTriple) {
    return { isValid: false, normalizedCards: sortedCards, comparisonKey: 0 };
  }

  const comparisonKey = getTubeTripleComparisonKey(bestTriple);
  const normalizedCards = normalizeTubeWithTriple(sortedCards, bestTriple);

  return { isValid: true, normalizedCards, comparisonKey };
}
