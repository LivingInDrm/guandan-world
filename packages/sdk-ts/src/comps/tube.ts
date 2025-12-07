import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';
import { sortCards, TUBE_CARD_COUNT } from '../compUtil';
import { tubeSatisfy } from '../satisfy/tube';

export class Tube extends BaseComp {
  readonly comparisonKey: number;

  constructor(cards: SdkCard[]) {
    if (cards.length === TUBE_CARD_COUNT) {
      const result = tubeSatisfy(cards);
      super(cards, result.isValid, CompType.Tube, result.normalizedCards);
      this.comparisonKey = result.comparisonKey;
    } else {
      super(sortCards(cards), false, CompType.Tube);
      this.comparisonKey = 0;
    }
  }

  greaterThan(other: CardComp): boolean {
    if (other.getType() !== CompType.Tube) {
      return false;
    }
    const otherTube = other as Tube;
    return this.comparisonKey > otherTube.comparisonKey;
  }
}
