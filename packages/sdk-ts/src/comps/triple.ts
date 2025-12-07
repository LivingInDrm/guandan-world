import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';
import { sortCards, buildSameNumberComp } from '../compUtil';

export class Triple extends BaseComp {
  constructor(cards: SdkCard[]) {
    const sortedCards = sortCards(cards);
    const result = buildSameNumberComp(sortedCards, 3, 3);
    super(sortedCards, result.valid, CompType.Triple, result.normalizedCards);
  }

  greaterThan(other: CardComp): boolean {
    if (other.getType() !== CompType.Triple) {
      return false;
    }
    const otherTriple = other as Triple;
    const myCards = this.normalizedCards.length > 0 ? this.normalizedCards : this.cards;
    const otherCards =
      otherTriple.normalizedCards.length > 0 ? otherTriple.normalizedCards : otherTriple.cards;
    return myCards[0].greaterThan(otherCards[0]);
  }
}
