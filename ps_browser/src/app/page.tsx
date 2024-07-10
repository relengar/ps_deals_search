'use client';

import Button from '@/components/button';
import SearchIcon from '@/components/icons/searchIcon';
import { SearchGameParams, searchGames } from '@/lib/queries/searchGames';
import { GameResponse } from '@/lib/repositories/games';
import { Platform } from '@/lib/repositories/games/schema';
import { parseArrayParam, parseUrlNumber } from '@/lib/utils/url';
import { ChangeEvent, useEffect, useReducer, useState } from 'react';
import Pagination from '../components/pagination';
import Spinner from '../components/spinner';
import Filters from './filters';
import GamesList from './gamesList';
import { goToSearch } from './redirectToSearch';

type UrlQuery = {
    term?: string;
    maxPrice?: string;
    useSemantic?: string;
    platforms?: string;
    page?: string;
    limit?: string;
};

export default function Search({ searchParams }: { searchParams: UrlQuery }) {
    const [totalGames, setTotalGames] = useState<number | null>(null);
    const [term, setTerm] = useState(searchParams.term ?? '');
    const onTermChange = (evt: ChangeEvent<HTMLInputElement>) =>
        setTerm(evt.target.value);

    const maxPrice = parseUrlNumber(searchParams.maxPrice, 20);
    const useSemantic = searchParams.useSemantic === 'true';

    const platforms = parseArrayParam<Platform>(searchParams.platforms, [
        'PS4',
        'PS5',
    ]);

    let page = parseUrlNumber(searchParams.page, 0);
    const limit = parseUrlNumber(searchParams.limit, 20);

    const [filters, dispatchFilters] = useReducer(FiltersReducer, {
        maxPrice,
        useSemantic,
        platforms,
    });

    const onFiltersChange = (changes: Partial<SearchGameParams>) => {
        dispatchFilters(changes);
    };

    const [loading, setLoading] = useState(true);
    const [games, setGames] = useState<GameResponse[]>([]);

    const search = () => {
        if (term.length === 0) {
            dispatchFilters({ useSemantic: false });
        }

        goToSearch({ term, ...filters, page, limit });
    };

    const onPageSelect = (selected: number) => {
        page = selected;
        search();
    };

    useEffect(() => {
        const getGames = async () => {
            setLoading(true);
            const { games: result, total } = await searchGames({
                term,
                ...filters,
                page,
                limit,
            });
            setTotalGames(total);
            setGames(result);
            setLoading(false);
        };

        getGames();
    }, [searchParams]);

    return (
        <section className="container mx-auto flex-col content-center">
            <h2 className="text-center w-full text-3xl leading-loose p-10">
                Search for deals
            </h2>
            <section>
                <form action={search}>
                    <div className="flex-col w-full mb-14 place-items-center justify-center gap-2">
                        <div className="flex w-full mb-6 place-items-center justify-center gap-2">
                            <input
                                className="resize rounded-md center w-5/6 h-8"
                                placeholder="Search"
                                type="text"
                                name="term"
                                value={term}
                                onChange={onTermChange}
                            />
                            <Button className="p-0">
                                <SearchIcon className="h-8 w-8 p-2 m-0.5" />
                            </Button>
                        </div>
                        <Filters
                            term={term}
                            {...filters}
                            onChange={onFiltersChange}
                        />
                    </div>
                </form>
            </section>
            {loading && <Spinner />}
            <GamesList games={games} />
            {totalGames && totalGames - limit > 0 && (
                <Pagination
                    limit={limit}
                    total={totalGames}
                    page={page}
                    onChange={onPageSelect}
                />
            )}
        </section>
    );
}

function FiltersReducer(
    filters: SearchGameParams,
    changes: Partial<SearchGameParams>
): SearchGameParams {
    return { ...filters, ...changes };
}
