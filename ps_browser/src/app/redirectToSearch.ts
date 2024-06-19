'use server';
import { headers } from 'next/headers';
import { logger } from '@/lib/logger';
import { redirect } from 'next/navigation';
import { SearchGameParams } from '@/lib/queries/searchGames';

export async function goToSearch(filters: SearchGameParams) {
    logger.debug({ form: filters }, 'Searching');
    const nextUrl = new URL(headers().get('referer') ?? 'http://localhost');

    nextUrl.search = toQueryString(filters);

    logger.debug({ nextUrl }, 'Redirecting to new search');
    redirect(nextUrl.href);
}

function toQueryString(filters: SearchGameParams): string {
    const params = new URLSearchParams();

    for (const [attr, value] of Object.entries(filters)) {
        params.set(attr, value.toString());
    }

    return params.toString();
}
