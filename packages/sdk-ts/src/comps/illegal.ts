import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';

export class Illegal extends BaseComp {
  constructor(cards: SdkCard[]) {
    super(cards, false, CompType.Illegal);
  }

  greaterThan(_other: CardComp): boolean {
    return false;
  }
}
