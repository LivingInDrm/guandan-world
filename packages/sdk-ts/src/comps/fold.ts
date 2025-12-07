import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';

export class Fold extends BaseComp {
  constructor(cards: SdkCard[]) {
    super(cards, true, CompType.Fold);
  }

  greaterThan(_other: CardComp): boolean {
    return false;
  }
}
