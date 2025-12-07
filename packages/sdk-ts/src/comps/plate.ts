import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';
import { sortCards, PLATE_CARD_COUNT } from '../compUtil';
import { plateSatisfy } from '../satisfy/plate';

export class Plate extends BaseComp {
  readonly comparisonKey: number;

  constructor(cards: SdkCard[]) {
    if (cards.length === PLATE_CARD_COUNT) {
      const result = plateSatisfy(cards);
      super(cards, result.isValid, CompType.Plate, result.normalizedCards);
      this.comparisonKey = result.comparisonKey;
    } else {
      super(sortCards(cards), false, CompType.Plate);
      this.comparisonKey = 0;
    }
  }

  greaterThan(other: CardComp): boolean {
    if (other.getType() !== CompType.Plate) {
      return false;
    }
    const otherPlate = other as Plate;
    return this.comparisonKey > otherPlate.comparisonKey;
  }
}
