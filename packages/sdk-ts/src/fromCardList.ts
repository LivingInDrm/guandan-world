import { SdkCard } from './card';
import { CompType } from './compType';
import type { CardComp } from './compInterface';
import { Fold } from './comps/fold';
import { Illegal } from './comps/illegal';
import { Single } from './comps/single';
import { Pair } from './comps/pair';
import { Triple } from './comps/triple';
import { FullHouse } from './comps/fullhouse';
import { Straight } from './comps/straight';
import { Plate } from './comps/plate';
import { Tube } from './comps/tube';
import { JokerBomb } from './comps/jokerBomb';
import { NaiveBomb } from './comps/naiveBomb';
import { StraightFlush } from './comps/straightFlush';

export function fromCardList(cards: SdkCard[], prev?: CardComp): CardComp {
  return fromCardListWithType(cards, prev?.getType());
}

export function fromCardListWithType(cards: SdkCard[], prevType?: CompType): CardComp {
  if (cards.length === 0) {
    return new Fold(cards);
  }

  switch (cards.length) {
    case 1: {
      const comp = new Single(cards);
      if (comp.isValid()) return comp;
      return new Illegal(cards);
    }

    case 2: {
      const comp = new Pair(cards);
      if (comp.isValid()) return comp;
      return new Illegal(cards);
    }

    case 3: {
      const comp = new Triple(cards);
      if (comp.isValid()) return comp;
      return new Illegal(cards);
    }

    case 4: {
      const jokerBomb = new JokerBomb(cards);
      if (jokerBomb.isValid()) return jokerBomb;

      const naiveBomb = new NaiveBomb(cards);
      if (naiveBomb.isValid()) return naiveBomb;

      return new Illegal(cards);
    }

    case 5: {
      const straightFlush = new StraightFlush(cards);
      if (straightFlush.isValid()) return straightFlush;

      const naiveBomb = new NaiveBomb(cards);
      if (naiveBomb.isValid()) return naiveBomb;

      if (prevType === undefined || prevType !== CompType.Straight) {
        const fullHouse = new FullHouse(cards);
        if (fullHouse.isValid()) return fullHouse;
      }

      const straight = new Straight(cards);
      if (straight.isValid()) return straight;

      const fullHouse = new FullHouse(cards);
      if (fullHouse.isValid()) return fullHouse;

      return new Illegal(cards);
    }

    case 6: {
      const naiveBomb = new NaiveBomb(cards);
      if (naiveBomb.isValid()) return naiveBomb;

      if (prevType === undefined || prevType !== CompType.Plate) {
        const tube = new Tube(cards);
        if (tube.isValid()) return tube;
      }

      const plate = new Plate(cards);
      if (plate.isValid()) return plate;

      const tube = new Tube(cards);
      if (tube.isValid()) return tube;

      return new Illegal(cards);
    }

    default: {
      const naiveBomb = new NaiveBomb(cards);
      if (naiveBomb.isValid()) return naiveBomb;

      return new Illegal(cards);
    }
  }
}
