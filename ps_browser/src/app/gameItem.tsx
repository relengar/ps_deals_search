import { GameResponse } from '@/lib/repositories/games';

export default function GameItem({ game }: { game: GameResponse }) {
    return (
        <div className="shadow-md w-full text-center rounded-sm border-slate-500 rounded-md p-3">
            <div>
                <a href={game.url} target="_blank">
                    <span className="text-xl">{game.name}</span>
                </a>
            </div>
            <div className="flex justify-between">
                <span>
                    Price: <span>{game.price}</span>&nbsp;/&nbsp;
                    <span className="line-through text-slate-400">
                        {game.originalPrice}
                    </span>
                </span>
                <span>
                    Rating: {game.rating}&nbsp;
                    <span className="text-slate-400">
                        (of {game.ratingSum})
                    </span>
                </span>
            </div>
            <div className="flex justify-between">
                <span>Expires: {game.expiration.toDateString()}</span>
            </div>
            <div>
                <span className="font-thin">{game.description}</span>
            </div>
        </div>
    );
}
