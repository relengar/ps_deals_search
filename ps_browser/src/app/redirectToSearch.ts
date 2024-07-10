'use server';
import { logger } from '@/lib/logger';
import { SearchGameParams } from '@/lib/queries/searchGames';
import { Pagination } from '@/lib/repositories/schema';
import { toQueryString } from '@/lib/utils/url';
import { headers } from 'next/headers';
import { RedirectType, redirect } from 'next/navigation';

export async function goToSearch(filters: SearchGameParams & Pagination) {
    logger.debug({ form: filters }, 'Searching');
    const nextUrl = new URL(headers().get('referer') ?? 'http://localhost');

    nextUrl.search = toQueryString(filters);

    logger.debug({ nextUrl }, 'Redirecting to new search');
    redirect(nextUrl.href, RedirectType.replace);
}
