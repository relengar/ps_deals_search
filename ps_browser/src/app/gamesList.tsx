import { GameResponse } from '@/lib/repositories/games';
import { Fragment } from 'react';
import GameItem from './gameItem';

export default function GamesList({ games }: { games: GameResponse[] }) {
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
