import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';
import { sortCards, FULL_HOUSE_CARD_COUNT } from '../compUtil';
import { fullHouseSatisfy } from '../satisfy/fullhouse';

export class FullHouse extends BaseComp {
  constructor(cards: SdkCard[]) {
    if (cards.length === FULL_HOUSE_CARD_COUNT) {
      const result = fullHouseSatisfy(cards);
      super(result.normalizedCards, result.isValid, CompType.FullHouse, result.normalizedCards);
    } else {
      super(sortCards(cards), false, CompType.FullHouse);
    }
  }

  greaterThan(other: CardComp): boolean {
    if (other.getType() !== CompType.FullHouse) {
      return false;
    }
    const otherFullHouse = other as FullHouse;
    const myCards = this.normalizedCards.length > 0 ? this.normalizedCards : this.cards;
    const otherCards =
      otherFullHouse.normalizedCards.length > 0
        ? otherFullHouse.normalizedCards
        : otherFullHouse.cards;
    return myCards[0].greaterThan(otherCards[0]);
  }
}
