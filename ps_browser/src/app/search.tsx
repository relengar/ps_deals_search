import { Game } from '@/lib/connectors/postgres/schema';
import { searchGames } from '@/lib/queries/searchGames';
import GameItem from './gameItem';
import { Fragment } from 'react';

type SearchGamesParams = {
    term?: string;
};

export default async function SearchGames({ term }: SearchGamesParams) {
    const { games } = await searchGames({ term });

    return (
        <section className="container mx-auto space-y-10">
            {games.map((game) => (
                // @ts-ignore
                <Fragment key={game.id}>
                    <GameItem game={game} />
                </Fragment>
            ))}
        </section>
    );
}
