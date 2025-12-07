import { describe, it, expect, beforeAll } from 'vitest';
import * as fs from 'fs';
import * as path from 'path';
import { fromCardList } from '../src/fromCardList';
import {
  convertJSONToCards,
  createCompByType,
  normalizeComp,
  formatCompForLog,
  ComparisonTestData,
  ComparisonTestCase,
} from './testHelper';

describe('ComparisonBatch', () => {
  let testData: ComparisonTestData;

  beforeAll(() => {
    const testDataPath = path.join(__dirname, '..', '..', '..', 'test-data', 'comparison_test_data.json');
    const data = fs.readFileSync(testDataPath, 'utf-8');
    testData = JSON.parse(data) as ComparisonTestData;
  });

  it('should load test data', () => {
    expect(testData).toBeDefined();
    expect(testData.level).toBe(5);
    expect(testData.comparisons.length).toBeGreaterThan(0);
    console.log(`Loaded ${testData.comparisons.length} test cases`);
  });

  it('should pass all comparison tests', () => {
    const level = testData.level;
    let passCount = 0;
    let failCount = 0;
    const failures: string[] = [];

    for (const testCase of testData.comparisons) {
      const comp2Cards = convertJSONToCards(testCase.comp2.cards, level);
      const comp2 = createCompByType(comp2Cards, testCase.comp2.type);
      const normalizedComp2 = normalizeComp(comp2, level);

      const comp1Cards = convertJSONToCards(testCase.comp1.cards, level);
      const comp1 = fromCardList(comp1Cards, normalizedComp2);

      const actualComp1Greater = comp1.greaterThan(normalizedComp2);
      const actualComp2Greater = normalizedComp2.greaterThan(comp1);

      const comp1Matches = actualComp1Greater === testCase.comp1_greater_than_comp2;
      const comp2Matches = actualComp2Greater === testCase.comp2_greater_than_comp1;

      if (comp1Matches && comp2Matches) {
        passCount++;
      } else {
        failCount++;
        if (failures.length < 20) {
          failures.push(
            `[TestID:${testCase.test_id}] ${testCase.comparison_type}/${testCase.comp_type}\n` +
              `  Comp1: ${formatCompForLog(comp1)}\n` +
              `  Comp2: ${formatCompForLog(normalizedComp2)}\n` +
              `  Expected: comp1>comp2=${testCase.comp1_greater_than_comp2}, comp2>comp1=${testCase.comp2_greater_than_comp1}\n` +
              `  Actual:   comp1>comp2=${actualComp1Greater}, comp2>comp1=${actualComp2Greater}`
          );
        }
      }
    }

    console.log(`\n=== Comparison Test Results ===`);
    console.log(`Total: ${testData.comparisons.length}`);
    console.log(`Passed: ${passCount} (${((passCount / testData.comparisons.length) * 100).toFixed(1)}%)`);
    console.log(`Failed: ${failCount} (${((failCount / testData.comparisons.length) * 100).toFixed(1)}%)`);

    if (failures.length > 0) {
      console.log(`\n=== First ${failures.length} Failures ===`);
      failures.forEach((f) => console.log(f));
    }

    expect(failCount).toBe(0);
  });
});

describe('ComparisonSpecific', () => {
  let testData: ComparisonTestData;

  beforeAll(() => {
    const testDataPath = path.join(__dirname, '..', '..', '..', 'test-data', 'comparison_test_data.json');
    const data = fs.readFileSync(testDataPath, 'utf-8');
    testData = JSON.parse(data) as ComparisonTestData;
  });

  const specificTestIDs = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];

  specificTestIDs.forEach((targetTestID) => {
    it(`TestID_${targetTestID}`, () => {
      const testCase = testData.comparisons.find((tc) => tc.test_id === targetTestID);
      if (!testCase) {
        console.log(`TestID ${targetTestID} not found, skipping`);
        return;
      }

      const level = testData.level;

      const comp1Cards = convertJSONToCards(testCase.comp1.cards, level);
      const comp1 = fromCardList(comp1Cards);

      const comp2Cards = convertJSONToCards(testCase.comp2.cards, level);
      const comp2 = fromCardList(comp2Cards);

      const actualComp1Greater = comp1.greaterThan(comp2);
      const actualComp2Greater = comp2.greaterThan(comp1);

      console.log(`[TestID:${testCase.test_id}] ${testCase.comparison_type}:`);
      console.log(`  Comp1: ${formatCompForLog(comp1)}`);
      console.log(`  Comp2: ${formatCompForLog(comp2)}`);
      console.log(
        `  Expected: comp1>comp2=${testCase.comp1_greater_than_comp2}, comp2>comp1=${testCase.comp2_greater_than_comp1}`
      );
      console.log(`  Actual:   comp1>comp2=${actualComp1Greater}, comp2>comp1=${actualComp2Greater}`);

      expect(actualComp1Greater).toBe(testCase.comp1_greater_than_comp2);
      expect(actualComp2Greater).toBe(testCase.comp2_greater_than_comp1);
    });
  });
});
