import { sql } from 'kysely';
import { PgClient, getPgClient } from '@/lib/connectors/postgres';

let gamesRepo: GamesRepository;

type Pagination = {};

export type Platform = 'PS4' | 'PS5';

export type GameFilters = Pagination & {
    maxPrice?: number;
    platforms?: Platform[];
};

export type GameResponse = {
    id: number;
    price: number;
    rating: number;
    name: string;
    description: string;
    url: string;
    expiration: Date;
    ratingSum: number;
    originalPrice: number;
};

class GamesRepository {
    #pg: PgClient;
    constructor(pg: PgClient) {
        this.#pg = pg;
    }

    #query() {
        return this.#pg.db.selectFrom('games');
    }

    async getGames({
        embedding,
        filters,
    }: {
        embedding?: number[];
        filters: GameFilters;
    }): Promise<GameResponse[]> {
        let query = this.#query().select([
            'games.id',
            'games.name',
            'games.description',
            'games.expiration',
            'games.rating',
            'games.rating_sum as ratingSum',
            'games.price',
            'games.original_price as originalPrice',
            'games.url',
        ]);

        if (filters.maxPrice) {
            query = query.where('price', '<=', filters.maxPrice);
        }

        if (filters.platforms?.length) {
            query = query.where(
                'platforms',
                '<@',
                sql<string>`ARRAY[${sql.join(filters.platforms)}]`
            );
        }

        if (embedding) {
            query = query
                .innerJoin('game_embeddings', 'game_id', 'games.id')
                .orderBy(
                    sql`${sql.ref(
                        'game_embeddings.embedding'
                    )} <=> ${JSON.stringify(embedding)}`
                );
        }

        query = query.orderBy('ratingSum desc').orderBy('rating desc');

        return query.execute();
    }
}

export function getGamesRepo() {
    if (gamesRepo) {
        return gamesRepo;
    }

    const pg = getPgClient();
    const gamesRespo = new GamesRepository(pg);

    return gamesRespo;
}
