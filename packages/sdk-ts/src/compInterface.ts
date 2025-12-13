import { SdkCard } from './card';
import { CompType } from './compType';

export interface CardComp {
  greaterThan(other: CardComp): boolean;
  isBomb(): boolean;
  getCards(): SdkCard[];
  getNormalizedCards(): SdkCard[];
  toString(): string;
  isValid(): boolean;
  getType(): CompType;
}
