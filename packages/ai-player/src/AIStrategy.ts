import {
  SdkCard,
  findMinPlay,
  FirstPlayRecommender,
  TributeRecommender
} from '@guandan/sdk-ts';
import type { CardComp } from '@guandan/sdk-ts';

export class AIStrategy {
  private level: number;

  constructor(level: number = 1) {
    this.level = level;
  }

  selectCardsToPlay(hand: SdkCard[], isLeader: boolean, prevComp?: CardComp): SdkCard[] | null {
    if (isLeader) {
      const recommender = new FirstPlayRecommender(hand);
      const next = recommender.next();
      return next || (hand.length > 0 ? [hand[0]] : null);
    } else {
      if (!prevComp) {
        return null;
      }
      const minPlay = findMinPlay(hand, prevComp);
      return minPlay?.getCards() || null;
    }
  }

  selectTributeCard(poolCards: SdkCard[]): SdkCard | null {
    if (poolCards.length === 0) return null;
    const sorted = [...poolCards].sort((a, b) => {
      if (b.greaterThan(a)) return 1;
      if (a.greaterThan(b)) return -1;
      return 0;
    });
    return sorted[0] || null;
  }

  selectReturnTributeCard(hand: SdkCard[]): SdkCard | null {
    if (hand.length === 0) return null;
    const recommender = new TributeRecommender(hand);
    const card = recommender.next();
    if (card) {
      return card;
    }
    const sorted = [...hand].sort((a, b) => {
      if (a.greaterThan(b)) return 1;
      if (b.greaterThan(a)) return -1;
      return 0;
    });
    return sorted[0];
  }

  getLevel(): number {
    return this.level;
  }
}
