'use client';

import Button from '@/components/button';
import SearchIcon from '@/components/icons/searchIcon';
import { Game } from '@/lib/connectors/postgres/schema';
import { searchGames } from '@/lib/queries/searchGames';
import { ChangeEvent, useEffect, useReducer, useState } from 'react';
import Spinner from '../components/spinner';
import Filters, { FilterValues } from './filters';
import GamesList from './gamesList';
import { goToSearch } from './redirectToSearch';

type UrlQuery = {
    term?: string;
    maxPrice?: string;
    useSemantic?: string;
};

export default function Search({ searchParams }: { searchParams: UrlQuery }) {
    const [term, setTerm] = useState(searchParams.term ?? '');
    const onTermChange = (evt: ChangeEvent<HTMLInputElement>) =>
        setTerm(evt.target.value);

    const defaultMaxPrice = 20;
    const maxPrice = searchParams.maxPrice
        ? Number(searchParams.maxPrice)
        : defaultMaxPrice;
    const useSemantic = searchParams.useSemantic === 'true';

    const [filters, dispatchFilters] = useReducer(FiltersReducer, {
        maxPrice,
        useSemantic,
    });

    const onFiltersChange = (changes: Partial<FilterValues>) => {
        dispatchFilters(changes);
    };

    const [loading, setLoading] = useState(true);
    const [games, setGames] = useState<Game[]>([]);

    const search = () => {
        if (term.length === 0) {
            dispatchFilters({ useSemantic: false });
        }

        goToSearch({ term, ...filters });
    };

    useEffect(() => {
        const getGames = async () => {
            setLoading(true);
            const result = await searchGames({ term, ...filters });
            setGames(result);
            setLoading(false);
        };

        getGames();
    }, [searchParams.term]);

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
                            maxPrice={maxPrice}
                            useSemantic={useSemantic}
                            onChange={onFiltersChange}
                        />
                    </div>
                </form>
            </section>
            {loading && <Spinner />}
            <GamesList games={games} />
        </section>
    );
}

function FiltersReducer(
    filters: FilterValues,
    changes: Partial<FilterValues>
): FilterValues {
    return { ...filters, ...changes };
}
