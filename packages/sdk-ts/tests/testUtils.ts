import { SdkCard } from '../src/card';
import { Suit } from '../src/types';
import type { CardComp } from '../src/compInterface';

export function cardsToShortStrings(comp: CardComp | null): string[] | null {
  if (!comp) return null;
  return comp.getCards().map(c => c.toShortString());
}

export function createStandardHandCards(level: number): SdkCard[] {
  return [
    new SdkCard(5, Suit.Heart, level, 0),
    new SdkCard(5, Suit.Heart, level, 1),
    new SdkCard(7, 0, level, 0),
    new SdkCard(7, 0, level, 1),
    new SdkCard(7, Suit.Heart, level, 0),
    new SdkCard(7, Suit.Club, level, 0),
    new SdkCard(13, 0, level, 0),
    new SdkCard(13, Suit.Heart, level, 0),
    new SdkCard(13, Suit.Club, level, 0),
    new SdkCard(13, Suit.Club, level, 1),
    new SdkCard(8, 0, level, 0),
    new SdkCard(8, Suit.Heart, level, 0),
    new SdkCard(8, Suit.Club, level, 0),
    new SdkCard(3, 0, level, 0),
    new SdkCard(3, 0, level, 1),
    new SdkCard(4, 0, level, 0),
    new SdkCard(4, 0, level, 1),
    new SdkCard(6, 0, level, 0),
    new SdkCard(6, 0, level, 1),
    new SdkCard(9, Suit.Heart, level, 0),
    new SdkCard(9, Suit.Club, level, 0),
    new SdkCard(10, Suit.Heart, level, 0),
    new SdkCard(10, Suit.Club, level, 0),
    new SdkCard(15, Suit.Joker, level, 0),
    new SdkCard(15, Suit.Joker, level, 1),
    new SdkCard(16, Suit.Joker, level, 0),
    new SdkCard(16, Suit.Joker, level, 1),
  ];
}
