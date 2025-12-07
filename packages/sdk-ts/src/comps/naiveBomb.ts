import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';
import { sortCards, buildSameNumberComp } from '../compUtil';

export class NaiveBomb extends BaseComp {
  constructor(cards: SdkCard[]) {
    const sortedCards = sortCards(cards);
    const result = buildSameNumberComp(sortedCards, 4, 8);
    super(sortedCards, result.valid, CompType.NaiveBomb, result.normalizedCards);
  }

  greaterThan(other: CardComp): boolean {
    if (!other.isBomb()) {
      return true;
    }

    if (other.getType() === CompType.JokerBomb) {
      return false;
    }

    if (other.getType() === CompType.StraightFlush) {
      return this.cards.length >= 6;
    }

    if (other.getType() === CompType.NaiveBomb) {
      const otherBomb = other as NaiveBomb;
      if (this.cards.length > otherBomb.cards.length) {
        return true;
      } else if (this.cards.length < otherBomb.cards.length) {
        return false;
      } else {
        const myCards = this.normalizedCards.length > 0 ? this.normalizedCards : this.cards;
        const otherCards =
          otherBomb.normalizedCards.length > 0 ? otherBomb.normalizedCards : otherBomb.cards;
        return myCards[0].greaterThan(otherCards[0]);
      }
    }

    return false;
  }
}
