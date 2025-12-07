import { expect } from 'vitest';
import type { CardComp } from '../src/compInterface';
import { cardsToShortStrings, createStandardHandCards } from './testUtils';

export { cardsToShortStrings, createStandardHandCards };

export function expectCardsEqual(
  result: CardComp | null,
  expected: string[] | null
): void {
  if (expected === null) {
    expect(result).toBeNull();
    return;
  }
  expect(result).not.toBeNull();
  expect(cardsToShortStrings(result)).toEqual(expected);
}
