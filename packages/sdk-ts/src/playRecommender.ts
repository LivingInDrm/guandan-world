import { SdkCard } from './card';
import type { CardComp } from './compInterface';
import { findAllBombs, findMinPlay } from './findMinPlay';

function groupCardsByRank(cards: SdkCard[]): SdkCard[][] {
  const byRank = new Map<number, SdkCard[]>();

  for (const card of cards) {
    const list = byRank.get(card.rank) ?? [];
    list.push(card);
    byRank.set(card.rank, list);
  }

  const groups = Array.from(byRank.entries());
  groups.sort((a, b) => {
    const cardA = a[1][0];
    const cardB = b[1][0];
    if (cardA.greaterThan(cardB)) return 1;
    if (cardB.greaterThan(cardA)) return -1;
    return 0;
  });

  return groups.map(([_rank, cards]) => cards);
}

export class FirstPlayRecommender {
  private groups: SdkCard[][];
  private currentIndex: number = -1;

  constructor(cards: SdkCard[]) {
    this.groups = groupCardsByRank(cards);
  }

  next(): SdkCard[] | null {
    if (this.groups.length === 0) return null;
    this.currentIndex = (this.currentIndex + 1) % this.groups.length;
    return this.groups[this.currentIndex];
  }
}

export class NextPlayRecommender {
  private prevComp: CardComp;
  private cards: SdkCard[];
  private currentComp: CardComp;

  constructor(cards: SdkCard[], prevComp: CardComp) {
    this.cards = cards;
    this.prevComp = prevComp;
    this.currentComp = prevComp;
  }

  next(): CardComp | null {
    const result = findMinPlay(this.cards, this.currentComp);
    if (result) {
      this.currentComp = result;
      return result;
    } else {
      this.currentComp = this.prevComp;
      return null;
    }
  }
}

export class TributeRecommender {
  private candidates: SdkCard[];
  private currentIndex: number = -1;

  constructor(cards: SdkCard[]) {
    const bombs = findAllBombs(cards);

    const bombCardKeys = new Set<string>();
    for (const bomb of bombs) {
      for (const card of bomb.getCards()) {
        bombCardKeys.add(`${card.deckIndex}-${card.suit}-${card.rank}`);
      }
    }

    const sortedCards = [...cards].sort((a, b) => {
      if (a.greaterThan(b)) return 1;
      if (b.greaterThan(a)) return -1;
      return 0;
    });

    const nonBombCards = sortedCards.filter((card) => {
      const key = `${card.deckIndex}-${card.suit}-${card.rank}`;
      return !bombCardKeys.has(key);
    });

    this.candidates = nonBombCards.length > 0 ? nonBombCards : sortedCards;
  }

  next(): SdkCard | null {
    if (this.candidates.length === 0) return null;
    this.currentIndex = (this.currentIndex + 1) % this.candidates.length;
    return this.candidates[this.currentIndex];
  }
}
