import { SdkCard } from './card';
import type { CardComp } from './compInterface';
import { findMinPlay } from './findMinPlay';

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
