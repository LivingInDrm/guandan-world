import { SdkCard } from '../src/card';
import { CompType } from '../src/compType';
import type { CardComp } from '../src/compInterface';
import { BaseComp } from '../src/compBase';
import { Suit } from '../src/types';
import { Single } from '../src/comps/single';
import { Pair } from '../src/comps/pair';
import { Triple } from '../src/comps/triple';
import { FullHouse } from '../src/comps/fullhouse';
import { Straight } from '../src/comps/straight';
import { Plate } from '../src/comps/plate';
import { Tube } from '../src/comps/tube';
import { JokerBomb } from '../src/comps/jokerBomb';
import { NaiveBomb } from '../src/comps/naiveBomb';
import { StraightFlush } from '../src/comps/straightFlush';
import { Fold } from '../src/comps/fold';
import { Illegal } from '../src/comps/illegal';
import { fromCardList } from '../src/fromCardList';

type CardJsonData = [number, string];

interface CompData {
  cards: CardJsonData[];
  type: string;
}

interface ComparisonTestCase {
  test_id: number;
  comparison_type: string;
  comp_type: string;
  wildcard_count?: number;
  wildcard_count_1?: number;
  wildcard_count_2?: number;
  comp1: CompData;
  comp2: CompData;
  comp1_greater_than_comp2: boolean;
  comp2_greater_than_comp1: boolean;
}

interface ComparisonTestData {
  level: number;
  description: string;
  total_comparisons: number;
  intra_type_comparisons: number;
  intra_type_cross_wildcard_comparisons: number;
  inter_type_comparisons: number;
  comparisons: ComparisonTestCase[];
}

function colorToSuit(color: string): number {
  switch (color) {
    case 'Spade':
      return Suit.Spade;
    case 'Heart':
      return Suit.Heart;
    case 'Club':
      return Suit.Club;
    case 'Diamond':
      return Suit.Diamond;
    case 'Joker':
      return Suit.Joker;
    default:
      return Suit.Spade;
  }
}

export function convertJSONToCards(cardDataList: CardJsonData[], level: number): SdkCard[] {
  const cards: SdkCard[] = [];
  for (const cardData of cardDataList) {
    const [rank, color] = cardData;
    const suit = colorToSuit(color);
    const card = new SdkCard(rank, suit, level);
    cards.push(card);
  }
  return cards;
}

export function createCompByType(cards: SdkCard[], typeName: string): CardComp {
  switch (typeName) {
    case 'Single':
      return new Single(cards);
    case 'Pair':
      return new Pair(cards);
    case 'Triple':
      return new Triple(cards);
    case 'FullHouse':
      return new FullHouse(cards);
    case 'Straight':
      return new Straight(cards);
    case 'Plate':
      return new Plate(cards);
    case 'Tube':
      return new Tube(cards);
    case 'JokerBomb':
      return new JokerBomb(cards);
    case 'NaiveBomb':
      return new NaiveBomb(cards);
    case 'StraightFlush':
      return new StraightFlush(cards);
    case 'Fold':
      return new Fold(cards);
    case 'IllegalComp':
      return new Illegal(cards);
    default:
      return fromCardList(cards);
  }
}

export function normalizeComp(comp: CardComp, _level: number): CardComp {
  if (!comp.isValid()) {
    return comp;
  }

  const baseComp = comp as BaseComp;
  if (baseComp.normalizedCards && baseComp.normalizedCards.length > 0) {
    return createCompByType(baseComp.normalizedCards, compTypeToTypeName(comp.getType()));
  }

  return comp;
}

function compTypeToTypeName(ct: CompType): string {
  switch (ct) {
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

export function formatCompForLog(comp: CardComp): string {
  const cards = comp.getCards();
  if (cards.length === 0) {
    return `${compTypeToTypeName(comp.getType())}: Empty`;
  }

  const cardStrs = cards.map((card) => card.toShortString());
  return `${compTypeToTypeName(comp.getType())}: [${cardStrs.join(',')}]`;
}

export { ComparisonTestData, ComparisonTestCase };
