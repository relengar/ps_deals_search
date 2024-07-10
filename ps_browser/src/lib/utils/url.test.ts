import { describe, expect, test } from 'vitest';
import { parseArrayParam, parseUrlNumber, toQueryString } from './url';

describe('utils:url', () => {
    describe('toQueryString', () => {
        test('single param', () => {
            const filters = { strProp: 'str' };
            const queryString = toQueryString(filters);

            expect(queryString).toEqual('strProp=str');
        });

        test('multiple params', () => {
            const filters = {
                strProp: 'str',
                numberProp: 2.32,
                boolProp: true,
            };
            const queryString = toQueryString(filters);

            expect(queryString).toEqual(
                'strProp=str&numberProp=2.32&boolProp=true'
            );
        });

        test('array param', () => {
            const filters = { arr: ['arr'] };
            const queryString = toQueryString(filters);

            expect(queryString).toEqual(`arr=${encodeURIComponent('["arr"]')}`);
        });
    });

    describe('parseArrayParam', () => {
        test('Simple string array', () => {
            const expectedArray = ['a', 'b', 'c'];
            const param = encodeURIComponent(JSON.stringify(expectedArray));

            const parsed = parseArrayParam(param);

            expect(parsed).toEqual(expectedArray);
        });

        test('Simple number array', () => {
            const expectedArray = [1, 2, 3];
            const param = encodeURIComponent(JSON.stringify(expectedArray));

            const parsed = parseArrayParam(param);

            expect(parsed).toEqual(expectedArray);
        });

        test('default', () => {
            const parsed = parseArrayParam(undefined);

            expect(parsed).toEqual([]);
        });

        test('custom default', () => {
            const expectedArray = ['meh'];
            const parsed = parseArrayParam(undefined, expectedArray);

            expect(parsed).toEqual(expectedArray);
        });
    });

    describe('parseUrlNumber', () => {
        test('parse valid number', () => {
            const number = parseUrlNumber('10');
            expect(number).toEqual(10);
        });

        test('use default on missing param', () => {
            const number = parseUrlNumber(undefined, 10);
            expect(number).toEqual(10);
        });

        test('use default on invalid param', () => {
            const number = parseUrlNumber('Not a number', 10);
            expect(number).toEqual(10);
        });
    });
});
