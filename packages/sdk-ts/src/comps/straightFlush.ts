import { SdkCard } from '../card';
import { CompType } from '../compType';
import { BaseComp } from '../compBase';
import type { CardComp } from '../compInterface';
import { sortCards, STRAIGHT_CARD_COUNT } from '../compUtil';
import { straightSatisfy } from '../satisfy/straight';

export class StraightFlush extends BaseComp {
  readonly comparisonKey: number;

  constructor(cards: SdkCard[]) {
    if (cards.length === STRAIGHT_CARD_COUNT) {
      const straightResult = straightSatisfy(cards);

      if (straightResult.isValid) {
        const colors = new Set<number>();
        for (const card of straightResult.normalizedCards) {
          if (!card.isWildcard()) {
            colors.add(card.suit);
          }
        }

        const valid = colors.size === 1;

        super(
          cards,
          valid,
          CompType.StraightFlush,
          valid ? straightResult.normalizedCards : undefined
        );
        this.comparisonKey = valid ? straightResult.comparisonKey : 0;
      } else {
        super(sortCards(cards), false, CompType.StraightFlush);
        this.comparisonKey = 0;
      }
    } else {
      super(sortCards(cards), false, CompType.StraightFlush);
      this.comparisonKey = 0;
    }
  }

  greaterThan(other: CardComp): boolean {
    if (!other.isBomb()) {
      return true;
    }

    if (other.getType() === CompType.JokerBomb) {
      return false;
    }

    if (other.getType() === CompType.StraightFlush) {
      const otherSF = other as StraightFlush;
      return this.comparisonKey > otherSF.comparisonKey;
    }

    if (other.getType() === CompType.NaiveBomb) {
      return other.getCards().length <= 5;
    }

    return false;
  }
}
