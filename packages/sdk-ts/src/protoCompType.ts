import { CompType } from './compType';

export enum ProtoCompType {
  UNSPECIFIED = 0,
  FOLD = 1,
  ILLEGAL = 2,
  SINGLE = 3,
  PAIR = 4,
  TRIPLE = 5,
  FULL_HOUSE = 6,
  STRAIGHT = 7,
  PLATE = 8,
  TUBE = 9,
  JOKER_BOMB = 10,
  NAIVE_BOMB = 11,
  STRAIGHT_FLUSH = 12,
}

export function protoToSdkCompType(protoType: ProtoCompType | number): CompType {
  switch (protoType) {
    case ProtoCompType.FOLD:
      return CompType.Fold;
    case ProtoCompType.ILLEGAL:
      return CompType.Illegal;
    case ProtoCompType.SINGLE:
      return CompType.Single;
    case ProtoCompType.PAIR:
      return CompType.Pair;
    case ProtoCompType.TRIPLE:
      return CompType.Triple;
    case ProtoCompType.FULL_HOUSE:
      return CompType.FullHouse;
    case ProtoCompType.STRAIGHT:
      return CompType.Straight;
    case ProtoCompType.PLATE:
      return CompType.Plate;
    case ProtoCompType.TUBE:
      return CompType.Tube;
    case ProtoCompType.JOKER_BOMB:
      return CompType.JokerBomb;
    case ProtoCompType.NAIVE_BOMB:
      return CompType.NaiveBomb;
    case ProtoCompType.STRAIGHT_FLUSH:
      return CompType.StraightFlush;
    default:
      return CompType.Illegal;
  }
}
