import { SdkCard } from './card';
import { CompType, compTypeToString } from './compType';
import type { CardComp } from './compInterface';

export class BaseComp implements CardComp {
  readonly cards: SdkCard[];
  readonly normalizedCards: SdkCard[];
  readonly valid: boolean;
  readonly type: CompType;

  constructor(
    cards: SdkCard[],
    valid: boolean,
    type: CompType,
    normalizedCards?: SdkCard[]
  ) {
    this.cards = cards;
    this.valid = valid;
    this.type = type;
    this.normalizedCards = normalizedCards ?? cards;
  }

  getCards(): SdkCard[] {
    return this.cards;
  }

  isValid(): boolean {
    return this.valid;
  }

  getType(): CompType {
    return this.type;
  }

  toString(): string {
    return `${compTypeToString(this.type)}: ${this.cards.map((c) => c.toShortString()).join(', ')}`;
  }

  isBomb(): boolean {
    return (
      this.type === CompType.JokerBomb ||
      this.type === CompType.NaiveBomb ||
      this.type === CompType.StraightFlush
    );
  }

  greaterThan(_other: CardComp): boolean {
    return false;
  }
}
