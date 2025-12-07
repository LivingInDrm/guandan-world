import { SdkCard } from './card';
import { CompType } from './compType';
import type { CardComp } from './compInterface';
import { JokerBomb } from './comps/jokerBomb';
import { NaiveBomb } from './comps/naiveBomb';
import { StraightFlush } from './comps/straightFlush';
import { Single } from './comps/single';
import { Pair } from './comps/pair';
import { Triple } from './comps/triple';
import { Straight } from './comps/straight';
import { Tube } from './comps/tube';
import { Plate } from './comps/plate';
import { FullHouse } from './comps/fullhouse';

interface HandAnalysis {
  byRank: Map<number, SdkCard[]>;
  bySuitRaw: Map<number, Map<number, SdkCard[]>>;
  byRawRank: Map<number, SdkCard[]>;
  rawRankCounts: Map<number, number>;
  wildcards: SdkCard[];
  smallJokers: SdkCard[];
  bigJokers: SdkCard[];
}

function analyzeHand(cards: SdkCard[]): HandAnalysis {
  const byRank = new Map<number, SdkCard[]>();
  const bySuitRaw = new Map<number, Map<number, SdkCard[]>>();
  const byRawRank = new Map<number, SdkCard[]>();
  const rawRankCounts = new Map<number, number>();
  const wildcards: SdkCard[] = [];
  const smallJokers: SdkCard[] = [];
  const bigJokers: SdkCard[] = [];

  for (const card of cards) {
    if (card.isSmallJoker()) {
      smallJokers.push(card);
      continue;
    }
    if (card.isBigJoker()) {
      bigJokers.push(card);
      continue;
    }
    if (card.isWildcard()) {
      wildcards.push(card);
      continue;
    }

    const rankList = byRank.get(card.rank) ?? [];
    rankList.push(card);
    byRank.set(card.rank, rankList);

    const rawRankList = byRawRank.get(card.rawRank) ?? [];
    rawRankList.push(card);
    byRawRank.set(card.rawRank, rawRankList);
    rawRankCounts.set(card.rawRank, (rawRankCounts.get(card.rawRank) ?? 0) + 1);

    if (!bySuitRaw.has(card.suit)) {
      bySuitRaw.set(card.suit, new Map());
    }
    const suitMap = bySuitRaw.get(card.suit)!;
    const rawList = suitMap.get(card.rawRank) ?? [];
    rawList.push(card);
    suitMap.set(card.rawRank, rawList);
  }

  return { byRank, bySuitRaw, byRawRank, rawRankCounts, wildcards, smallJokers, bigJokers };
}

const straightStartPositions = [
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

const tubeStartPositions = [
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

const plateStartPositions = [
  [1, 2], [2, 3], [3, 4], [4, 5], [5, 6], [6, 7], [7, 8],
  [8, 9], [9, 10], [10, 11], [11, 12], [12, 13], [13, 1],
];

function findAllBombs(cards: SdkCard[]): CardComp[] {
  const analysis = analyzeHand(cards);
  const bombs: CardComp[] = [];

  if (analysis.smallJokers.length >= 2 && analysis.bigJokers.length >= 2) {
    const jokerCards = [
      analysis.smallJokers[0],
      analysis.smallJokers[1],
      analysis.bigJokers[0],
      analysis.bigJokers[1],
    ];
    const jokerBomb = new JokerBomb(jokerCards);
    if (jokerBomb.isValid()) {
      bombs.push(jokerBomb);
    }
  }

  for (const [_rank, rankCards] of analysis.byRank) {
    const availableWildcards = analysis.wildcards.length;
    const totalAvailable = rankCards.length + availableWildcards;

    for (let size = 4; size <= Math.min(8, totalAvailable); size++) {
      const normalNeeded = Math.min(size, rankCards.length);
      const wildcardNeeded = size - normalNeeded;

      if (wildcardNeeded <= availableWildcards) {
        const bombCards = [
          ...rankCards.slice(0, normalNeeded),
          ...analysis.wildcards.slice(0, wildcardNeeded),
        ];
        const bomb = new NaiveBomb(bombCards);
        if (bomb.isValid()) {
          bombs.push(bomb);
        }
      }
    }
  }

  for (let suit = 0; suit <= 3; suit++) {
    const suitMap = analysis.bySuitRaw.get(suit);
    if (!suitMap) continue;

    for (const positions of straightStartPositions) {
      const selectedCards: SdkCard[] = [];
      let missing = 0;

      for (const rawRank of positions) {
        const cardsAtRank = suitMap.get(rawRank);
        if (cardsAtRank && cardsAtRank.length > 0) {
          selectedCards.push(cardsAtRank[0]);
        } else {
          missing++;
        }
      }

      if (missing <= analysis.wildcards.length) {
        const flushCards = [
          ...selectedCards,
          ...analysis.wildcards.slice(0, missing),
        ];
        const straightFlush = new StraightFlush(flushCards);
        if (straightFlush.isValid()) {
          bombs.push(straightFlush);
        }
      }
    }
  }

  return bombs;
}

function compareComps(a: CardComp, b: CardComp): number {
  if (a.greaterThan(b)) return 1;
  if (b.greaterThan(a)) return -1;
  const aWildcards = a.getCards().filter(c => c.isWildcard()).length;
  const bWildcards = b.getCards().filter(c => c.isWildcard()).length;
  return aWildcards - bWildcards;
}

function findFirstGreater(candidates: CardComp[], prev: CardComp): CardComp | null {
  candidates.sort(compareComps);
  return candidates.find(comp => comp.greaterThan(prev)) ?? null;
}

function fillWithWildcards(
  normalCards: SdkCard[],
  wildcards: SdkCard[],
  targetCount: number,
  wildcardOffset: number = 0
): { cards: SdkCard[]; wildcardUsed: number } {
  const result = normalCards.slice(0, targetCount);
  const needed = targetCount - result.length;
  result.push(...wildcards.slice(wildcardOffset, wildcardOffset + needed));
  return { cards: result, wildcardUsed: needed };
}

function createCompByCount(cards: SdkCard[], count: 1 | 2 | 3): CardComp {
  switch (count) {
    case 1:
      return new Single(cards);
    case 2:
      return new Pair(cards);
    case 3:
      return new Triple(cards);
  }
}

function findMinSameNumber(
  cards: SdkCard[],
  prev: CardComp,
  count: 1 | 2 | 3
): CardComp | null {
  const analysis = analyzeHand(cards);
  const candidates: SdkCard[][] = [];

  for (const [_rank, rankCards] of analysis.byRank) {
    if (rankCards.length >= count) {
      candidates.push(rankCards.slice(0, count));
    } else if (rankCards.length + analysis.wildcards.length >= count) {
      const needed = count - rankCards.length;
      candidates.push([...rankCards, ...analysis.wildcards.slice(0, needed)]);
    }
  }

  if (analysis.wildcards.length >= count) {
    candidates.push(analysis.wildcards.slice(0, count));
  }

  if (analysis.smallJokers.length >= count) {
    candidates.push(analysis.smallJokers.slice(0, count));
  }

  if (analysis.bigJokers.length >= count) {
    candidates.push(analysis.bigJokers.slice(0, count));
  }

  const comps = candidates
    .map((c) => createCompByCount(c, count))
    .filter((c) => c.isValid());
  return findFirstGreater(comps, prev);
}

function findMinStraight(cards: SdkCard[], prev: CardComp): CardComp | null {
  const analysis = analyzeHand(cards);
  const candidates: CardComp[] = [];

  for (const positions of straightStartPositions) {
    const selectedCards: SdkCard[] = [];
    let missing = 0;

    for (const rawRank of positions) {
      const cardsAtRank = analysis.byRawRank.get(rawRank) ?? [];
      if (cardsAtRank.length > 0) {
        selectedCards.push(cardsAtRank[0]);
      } else {
        missing++;
      }
    }

    if (missing <= analysis.wildcards.length) {
      const straightCards = [...selectedCards, ...analysis.wildcards.slice(0, missing)];
      const straight = new Straight(straightCards);
      if (straight.isValid()) {
        candidates.push(straight);
      }
    }
  }

  return findFirstGreater(candidates, prev);
}

function findMinTube(cards: SdkCard[], prev: CardComp): CardComp | null {
  const analysis = analyzeHand(cards);
  const candidates: CardComp[] = [];

  for (const triple of tubeStartPositions) {
    let needed = 0;
    for (const rawRank of triple) {
      const have = analysis.rawRankCounts.get(rawRank) ?? 0;
      if (have < 2) needed += 2 - have;
    }

    if (needed <= analysis.wildcards.length) {
      const selectedCards: SdkCard[] = [];
      let wildcardOffset = 0;
      for (const rawRank of triple) {
        const cardsAtRank = analysis.byRawRank.get(rawRank) ?? [];
        const filled = fillWithWildcards(cardsAtRank, analysis.wildcards, 2, wildcardOffset);
        selectedCards.push(...filled.cards);
        wildcardOffset += filled.wildcardUsed;
      }
      const tube = new Tube(selectedCards);
      if (tube.isValid()) {
        candidates.push(tube);
      }
    }
  }

  return findFirstGreater(candidates, prev);
}

function findMinPlate(cards: SdkCard[], prev: CardComp): CardComp | null {
  const analysis = analyzeHand(cards);
  const candidates: CardComp[] = [];

  for (const pair of plateStartPositions) {
    let needed = 0;
    for (const rawRank of pair) {
      const have = analysis.rawRankCounts.get(rawRank) ?? 0;
      if (have < 3) needed += 3 - have;
    }

    if (needed <= analysis.wildcards.length) {
      const selectedCards: SdkCard[] = [];
      let wildcardOffset = 0;
      for (const rawRank of pair) {
        const cardsAtRank = analysis.byRawRank.get(rawRank) ?? [];
        const filled = fillWithWildcards(cardsAtRank, analysis.wildcards, 3, wildcardOffset);
        selectedCards.push(...filled.cards);
        wildcardOffset += filled.wildcardUsed;
      }
      const plate = new Plate(selectedCards);
      if (plate.isValid()) {
        candidates.push(plate);
      }
    }
  }

  return findFirstGreater(candidates, prev);
}

function findMinFullHouse(cards: SdkCard[], prev: CardComp): CardComp | null {
  const analysis = analyzeHand(cards);

  const cardsByRank = new Map<number, SdkCard[]>();
  const cardCounts = new Map<number, number>();
  for (const [rank, rankCards] of analysis.byRank) {
    cardsByRank.set(rank, rankCards);
    cardCounts.set(rank, rankCards.length);
  }

  const candidates: CardComp[] = [];

  for (let tripleRank = 1; tripleRank <= 13; tripleRank++) {
    const tripleHave = cardCounts.get(tripleRank) ?? 0;
    const tripleNeed = Math.max(0, 3 - tripleHave);

    if (tripleNeed > analysis.wildcards.length) continue;

    const wildcardForTriple = tripleNeed;
    const remainingWildcards = analysis.wildcards.length - wildcardForTriple;

    for (let pairRank = 1; pairRank <= 13; pairRank++) {
      if (pairRank === tripleRank) continue;

      const pairHave = cardCounts.get(pairRank) ?? 0;
      const pairNeed = Math.max(0, 2 - pairHave);

      if (pairNeed > remainingWildcards) continue;

      const tripleCards = cardsByRank.get(tripleRank) ?? [];
      const tripleFilled = fillWithWildcards(tripleCards, analysis.wildcards, 3, 0);

      const pairCards = cardsByRank.get(pairRank) ?? [];
      const pairFilled = fillWithWildcards(pairCards, analysis.wildcards, 2, tripleFilled.wildcardUsed);

      const selectedCards = [...tripleFilled.cards, ...pairFilled.cards];
      const fullHouse = new FullHouse(selectedCards);
      if (fullHouse.isValid()) {
        candidates.push(fullHouse);
      }
    }

    if (analysis.smallJokers.length >= 2) {
      const tripleCards = cardsByRank.get(tripleRank) ?? [];
      const tripleFilled = fillWithWildcards(tripleCards, analysis.wildcards, 3, 0);
      const selectedCards = [...tripleFilled.cards, analysis.smallJokers[0], analysis.smallJokers[1]];

      const fullHouse = new FullHouse(selectedCards);
      if (fullHouse.isValid()) {
        candidates.push(fullHouse);
      }
    }

    if (analysis.bigJokers.length >= 2) {
      const tripleCards = cardsByRank.get(tripleRank) ?? [];
      const tripleFilled = fillWithWildcards(tripleCards, analysis.wildcards, 3, 0);
      const selectedCards = [...tripleFilled.cards, analysis.bigJokers[0], analysis.bigJokers[1]];

      const fullHouse = new FullHouse(selectedCards);
      if (fullHouse.isValid()) {
        candidates.push(fullHouse);
      }
    }
  }

  return findFirstGreater(candidates, prev);
}

export function findMinPlay(cards: SdkCard[], prev: CardComp): CardComp | null {
  const prevType = prev.getType();

  if (prev.isBomb()) {
    return findMinBomb(cards, prev);
  }

  const sameTypeResult = findMinSameType(cards, prev, prevType);
  if (sameTypeResult) return sameTypeResult;

  return findMinBomb(cards, prev);
}

function findMinSameType(
  cards: SdkCard[],
  prev: CardComp,
  type: CompType
): CardComp | null {
  switch (type) {
    case CompType.Single:
      return findMinSameNumber(cards, prev, 1);
    case CompType.Pair:
      return findMinSameNumber(cards, prev, 2);
    case CompType.Triple:
      return findMinSameNumber(cards, prev, 3);
    case CompType.Straight:
      return findMinStraight(cards, prev);
    case CompType.Tube:
      return findMinTube(cards, prev);
    case CompType.Plate:
      return findMinPlate(cards, prev);
    case CompType.FullHouse:
      return findMinFullHouse(cards, prev);
    default:
      return null;
  }
}

function findMinBomb(cards: SdkCard[], prev: CardComp): CardComp | null {
  const allBombs = findAllBombs(cards);
  return findFirstGreater(allBombs, prev);
}
