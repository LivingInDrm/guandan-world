export interface ProtoCard {
  suit: number;
  rank: number;
  deckIndex: number;
}

export const Suit = {
  Spade: 0,
  Heart: 1,
  Club: 2,
  Diamond: 3,
  Joker: -1,
} as const;

export type SuitType = (typeof Suit)[keyof typeof Suit];
