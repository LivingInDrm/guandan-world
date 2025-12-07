import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';
import { sortCards, buildSameNumberComp } from '../compUtil';

export class Pair extends BaseComp {
  constructor(cards: SdkCard[]) {
    const sortedCards = sortCards(cards);
    const result = buildSameNumberComp(sortedCards, 2, 2);
    super(sortedCards, result.valid, CompType.Pair, result.normalizedCards);
  }

  greaterThan(other: CardComp): boolean {
    if (other.getType() !== CompType.Pair) {
      return false;
    }
    const otherPair = other as Pair;
    const myCards = this.normalizedCards.length > 0 ? this.normalizedCards : this.cards;
    const otherCards =
      otherPair.normalizedCards.length > 0 ? otherPair.normalizedCards : otherPair.cards;
    return myCards[0].greaterThan(otherCards[0]);
  }
}
