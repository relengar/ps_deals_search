import { Game } from '@/lib/connectors/postgres/schema';
import { searchGames } from '@/lib/queries/searchGames';
import GameItem from './gameItem';
import { Fragment } from 'react';

type SearchGamesParams = {
    term?: string;
};

export default async function SearchGames({ term }: SearchGamesParams) {
    let data: Game[] | null = null;
    if (term && term.length > 0) {
        const { games } = await searchGames({ term });
        data = games;
    }

    return (
        <section className="container mx-auto space-y-10">
            {data &&
                data.map((game) => (
                    // @ts-ignore
                    <Fragment key={game.id}>
                        <GameItem game={game} />
                    </Fragment>
                ))}
        </section>
    );
}
