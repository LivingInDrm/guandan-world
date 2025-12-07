export enum CompType {
  Fold = 0,
  Illegal = 1,
  Single = 2,
  Pair = 3,
  Triple = 4,
  FullHouse = 5,
  Straight = 6,
  Plate = 7,
  Tube = 8,
  JokerBomb = 9,
  NaiveBomb = 10,
  StraightFlush = 11,
}

export function compTypeToString(ct: CompType): string {
  switch (ct) {
    case CompType.Fold:
      return 'Fold';
    case CompType.Illegal:
      return 'IllegalComp';
    case CompType.Single:
      return 'Single';
    case CompType.Pair:
      return 'Pair';
    case CompType.Triple:
      return 'Triple';
    case CompType.FullHouse:
      return 'FullHouse';
    case CompType.Straight:
      return 'Straight';
    case CompType.Plate:
      return 'Plate';
    case CompType.Tube:
      return 'Tube';
    case CompType.JokerBomb:
      return 'JokerBomb';
    case CompType.NaiveBomb:
      return 'NaiveBomb';
    case CompType.StraightFlush:
      return 'StraightFlush';
    default:
      return 'Unknown';
  }
}
