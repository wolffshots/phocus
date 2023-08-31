// stringUtils.test.js

import { camelCaseToWords } from '$lib/stringUtils';
import { describe, expect, test } from 'vitest';

describe('camelCaseToWords function', () => {
    test('converts camelCase to words with spaces', () => {
        const camelCaseString = 'BatteryStateOfCharge';
        const expectedOutput = 'Battery State Of Charge';

        const converted = camelCaseToWords(camelCaseString);

        expect(converted).toBe(expectedOutput);
    });

    test('converts preceeding caps camelCase to words with spaces', () => {
        const camelCaseString = 'ACInputVoltage';
        const expectedOutput = 'AC Input Voltage';

        const converted = camelCaseToWords(camelCaseString);

        expect(converted).toBe(expectedOutput);
    });

    test('converts middle caps camelCase to words with spaces', () => {
        const camelCaseString = 'TotalACOutputActivePower';
        const expectedOutput = 'Total AC Output Active Power';

        const converted = camelCaseToWords(camelCaseString);

        expect(converted).toBe(expectedOutput);
    });
});