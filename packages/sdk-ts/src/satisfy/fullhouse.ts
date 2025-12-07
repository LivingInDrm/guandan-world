import { SdkCard } from '../card';
import { FULL_HOUSE_CARD_COUNT, sortCards } from '../compUtil';

interface FullHouseCombo {
  tripleNum: number;
  pairNum: number;
}

function classifyFullHouseCards(cards: SdkCard[]): {
  wildcards: SdkCard[];
  numberGroups: Map<number, SdkCard[]>;
} {
  const wildcards: SdkCard[] = [];
  const numberGroups = new Map<number, SdkCard[]>();

  for (const card of cards) {
    if (card.isWildcard()) {
      wildcards.push(card);
    } else {
      const arr = numberGroups.get(card.rank) ?? [];
      arr.push(card);
      numberGroups.set(card.rank, arr);
    }
  }

  return { wildcards, numberGroups };
}

function enumerateFullHouseCombos(
  numberGroups: Map<number, SdkCard[]>,
  distinctNumbers: number[],
  wildcardCount: number
): FullHouseCombo[] {
  const validCombos: FullHouseCombo[] = [];

  for (const tripleNum of distinctNumbers) {
    if (tripleNum === 15 || tripleNum === 16) continue;

    for (const pairNum of distinctNumbers) {
      if (tripleNum === pairNum) continue;

      const tripleNormalCount = numberGroups.get(tripleNum)?.length ?? 0;
      const pairNormalCount = numberGroups.get(pairNum)?.length ?? 0;

      const wildcardNeeded = 3 - tripleNormalCount + (2 - pairNormalCount);

      if (wildcardNeeded <= wildcardCount && tripleNormalCount <= 3 && pairNormalCount <= 2) {
        validCombos.push({ tripleNum, pairNum });
      }
    }
  }

  return validCombos;
}

function selectBestFullHouseCombo(
  combos: FullHouseCombo[],
  numberGroups: Map<number, SdkCard[]>
): FullHouseCombo | null {
  if (combos.length === 0) return null;

  let bestCombo = combos[0];
  let bestTripleCard = numberGroups.get(bestCombo.tripleNum)![0];

  for (let i = 1; i < combos.length; i++) {
    const combo = combos[i];
    const tripleCard = numberGroups.get(combo.tripleNum)![0];
    if (tripleCard.greaterThan(bestTripleCard)) {
      bestCombo = combo;
      bestTripleCard = tripleCard;
    }
  }

  return bestCombo;
}

function buildFullHouseNormalizedCards(
  tripleNum: number,
  pairNum: number,
  numberGroups: Map<number, SdkCard[]>,
  wildcards: SdkCard[]
): SdkCard[] {
  const result: SdkCard[] = [];

  const tripleNormalCards = numberGroups.get(tripleNum) ?? [];
  result.push(...tripleNormalCards);

  const tripleNormalCount = tripleNormalCards.length;
  const wildcardForTriple = 3 - tripleNormalCount;
  let wildcardUsed = 0;

  for (let i = 0; i < wildcardForTriple && i < wildcards.length; i++) {
    result.push(wildcards[i]);
    wildcardUsed++;
  }

  const pairNormalCards = numberGroups.get(pairNum) ?? [];
  result.push(...pairNormalCards);

  for (let i = wildcardUsed; i < wildcards.length; i++) {
    result.push(wildcards[i]);
  }

  return result;
}

export interface FullHouseSatisfyResult {
  isValid: boolean;
  normalizedCards: SdkCard[];
}

export function fullHouseSatisfy(cards: SdkCard[]): FullHouseSatisfyResult {
  if (!cards || cards.length !== FULL_HOUSE_CARD_COUNT) {
    return { isValid: false, normalizedCards: sortCards(cards ?? []) };
  }

  const sortedCards = sortCards(cards);
  const { wildcards, numberGroups } = classifyFullHouseCards(sortedCards);
  const wildcardCount = wildcards.length;

  const distinctNumbers = Array.from(numberGroups.keys());

  if (distinctNumbers.length < 2) {
    return { isValid: false, normalizedCards: sortedCards };
  }

  const smallJokerCount = numberGroups.get(15)?.length ?? 0;
  const bigJokerCount = numberGroups.get(16)?.length ?? 0;
  const totalJokers = smallJokerCount + bigJokerCount;

  if (totalJokers > 0) {
    if (smallJokerCount !== 2 && bigJokerCount !== 2) {
      return { isValid: false, normalizedCards: sortedCards };
    }
  }

  const validCombos = enumerateFullHouseCombos(numberGroups, distinctNumbers, wildcardCount);

  if (validCombos.length === 0) {
    return { isValid: false, normalizedCards: sortedCards };
  }

  const bestCombo = selectBestFullHouseCombo(validCombos, numberGroups);

  if (!bestCombo) {
    return { isValid: false, normalizedCards: sortedCards };
  }

  const normalizedCards = buildFullHouseNormalizedCards(
    bestCombo.tripleNum,
    bestCombo.pairNum,
    numberGroups,
    wildcards
  );

  return { isValid: true, normalizedCards };
}
