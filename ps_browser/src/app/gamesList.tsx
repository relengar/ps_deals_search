import { Game } from '@/lib/connectors/postgres/schema';
import GameItem from './gameItem';
import { Fragment } from 'react';

export default function GamesList({ games }: { games: Game[] }) {
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
