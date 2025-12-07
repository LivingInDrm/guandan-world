import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';

export class Single extends BaseComp {
  constructor(cards: SdkCard[]) {
    const valid = cards.length === 1;
    super(cards, valid, CompType.Single);
  }

  greaterThan(other: CardComp): boolean {
    if (other.getType() !== CompType.Single) {
      return false;
    }
    const otherSingle = other as Single;
    return this.cards[0].greaterThan(otherSingle.cards[0]);
  }
}
