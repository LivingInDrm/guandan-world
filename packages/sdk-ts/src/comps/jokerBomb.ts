import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';
import { sortCards } from '../compUtil';

export class JokerBomb extends BaseComp {
  constructor(cards: SdkCard[]) {
    const sortedCards = sortCards(cards);
    let valid = false;

    if (cards.length === 4) {
      const numbers = sortedCards.map((c) => c.rank).sort((a, b) => a - b);
      valid =
        numbers.length === 4 &&
        numbers[0] === 15 &&
        numbers[1] === 15 &&
        numbers[2] === 16 &&
        numbers[3] === 16;
    }

    super(sortedCards, valid, CompType.JokerBomb);
  }

  greaterThan(other: CardComp): boolean {
    return other.getType() !== CompType.JokerBomb;
  }
}
