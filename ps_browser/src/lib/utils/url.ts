import { logger } from '../logger';

export function toQueryString(
    filters: Record<string, string | boolean | number | Array<unknown>>
): string {
    const params = new URLSearchParams();

    for (const [attr, value] of Object.entries(filters)) {
        if (Array.isArray(value)) {
            params.set(attr, JSON.stringify(value));
            continue;
        }
        params.set(attr, value.toString());
    }

    return params.toString();
}

export function parseArrayParam<T extends string>(
    param?: string,
    defaultValue: Array<T> = []
): Array<T> {
    try {
        if (param) {
            return JSON.parse(decodeURIComponent(param));
        }
    } catch (error) {
        logger.error(error, 'invalid url params');
    }

    return defaultValue;
}

export function parseUrlNumber(
    value?: string,
    defaultValue: number = 0
): number {
    if (!value) {
        return defaultValue;
    }

    const num = Number(value);
    return isNaN(num) ? defaultValue : num;
}
