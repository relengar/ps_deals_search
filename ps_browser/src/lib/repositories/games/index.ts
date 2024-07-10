import { PgClient, getPgClient } from '@/lib/connectors/postgres';
import { Database } from '@/lib/connectors/postgres/schema';
import { SelectQueryBuilder, sql } from 'kysely';
import { Pagination } from '../schema';
import { Platform } from './schema';

let gamesRepo: GamesRepository;

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

type GameFilterParams = {
    filters: GameFilters;
    embedding?: number[];
};

class GamesRepository {
    #pg: PgClient;
    constructor(pg: PgClient) {
        this.#pg = pg;
    }

    #query() {
        return this.#pg.db.selectFrom('games');
    }

    #applyGameFilters<T>({
        query,
        embedding,
        filters,
    }: GameFilterParams & {
        query: SelectQueryBuilder<Database, 'games', T>;
    }) {
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

        return query;
    }

    async getGames({
        embedding,
        filters,
    }: GameFilterParams): Promise<GameResponse[]> {
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

        query = this.#applyGameFilters({ query, embedding, filters });

        query = query.orderBy('ratingSum desc').orderBy('rating desc');

        query = query.limit(filters.limit).offset(filters.page * filters.limit);

        return query.execute();
    }

    async countGames({ embedding, filters }: GameFilterParams) {
        let query = this.#query().select((g) =>
            g.fn.count<number>('games.id').as('games_total')
        );

        query = this.#applyGameFilters({ query, embedding, filters });

        const { games_total } = await query.executeTakeFirstOrThrow();
        return games_total;
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
