import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';
import { sortCards, STRAIGHT_CARD_COUNT } from '../compUtil';
import { straightSatisfy } from '../satisfy/straight';

export class Straight extends BaseComp {
  readonly comparisonKey: number;

  constructor(cards: SdkCard[]) {
    if (cards.length === STRAIGHT_CARD_COUNT) {
      const result = straightSatisfy(cards);
      super(cards, result.isValid, CompType.Straight, result.normalizedCards);
      this.comparisonKey = result.comparisonKey;
    } else {
      super(sortCards(cards), false, CompType.Straight);
      this.comparisonKey = 0;
    }
  }

  greaterThan(other: CardComp): boolean {
    if (other.getType() !== CompType.Straight) {
      return false;
    }
    const otherStraight = other as Straight;
    return this.comparisonKey > otherStraight.comparisonKey;
  }
}
