import { SdkCard } from './card';
import { Suit } from './types';

export const STRAIGHT_CARD_COUNT = 5;
export const PLATE_CARD_COUNT = 6;
export const TUBE_CARD_COUNT = 6;
export const FULL_HOUSE_CARD_COUNT = 5;

export function sortCards(cards: SdkCard[]): SdkCard[] {
  const sorted = [...cards];
  sorted.sort((a, b) => (a.lessThan(b) ? -1 : a.equals(b) ? 0 : 1));
  return sorted;
}

export function sortCardsForConsecutive(cards: SdkCard[]): SdkCard[] {
  const sorted = [...cards];
  sorted.sort((a, b) => {
    if (!a.isWildcard() && b.isWildcard()) return -1;
    if (a.isWildcard() && !b.isWildcard()) return 1;
    if (a.isWildcard() && b.isWildcard()) return 0;
    return a.rawRank - b.rawRank;
  });
  return sorted;
}

export function countWildcards(cards: SdkCard[]): number {
  return cards.filter((c) => c.isWildcard()).length;
}

export function hasJokers(cards: SdkCard[]): boolean {
  return cards.some((c) => c.suit === Suit.Joker);
}

export function getNormalCards(cards: SdkCard[]): SdkCard[] {
  return cards.filter((c) => !c.isWildcard());
}

export function sortNormalFirst(cards: SdkCard[]): SdkCard[] {
  const sorted = [...cards];
  sorted.sort((a, b) => {
    if (!a.isWildcard() && b.isWildcard()) return -1;
    if (a.isWildcard() && !b.isWildcard()) return 1;
    return a.lessThan(b) ? -1 : a.equals(b) ? 0 : 1;
  });
  return sorted;
}

export interface BuildSameNumberResult {
  valid: boolean;
  normalizedCards: SdkCard[];
}

export function buildSameNumberComp(
  cards: SdkCard[],
  minLen: number,
  maxLen: number
): BuildSameNumberResult {
  if (cards.length < minLen || cards.length > maxLen) {
    return { valid: false, normalizedCards: [] };
  }

  const normalCards = getNormalCards(cards);
  const wildcardCount = countWildcards(cards);

  if (normalCards.length > 0) {
    const baseNumber = normalCards[0].rank;
    for (const card of normalCards) {
      if (card.rank !== baseNumber) {
        return { valid: false, normalizedCards: [] };
      }
    }

    if (normalCards[0].suit === Suit.Joker && wildcardCount > 0) {
      return { valid: false, normalizedCards: [] };
    }
  }

  return { valid: true, normalizedCards: sortNormalFirst(cards) };
}

export function extractDeckIndexes(cards: SdkCard[]): number[] {
  return cards.map((c) => c.deckIndex);
}

export function findCardsByDeckIndexes(source: SdkCard[], deckIndexes: number[]): SdkCard[] | null {
  const result: SdkCard[] = [];
  for (const idx of deckIndexes) {
    const found = source.find((c) => c.deckIndex === idx);
    if (!found) return null;
    result.push(found);
  }
  return result;
}

export function findCardByDeckIndex(source: SdkCard[], deckIndex: number): SdkCard | null {
  return source.find((c) => c.deckIndex === deckIndex) ?? null;
}
