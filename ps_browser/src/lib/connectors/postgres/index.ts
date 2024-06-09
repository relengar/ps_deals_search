import { Kysely, PostgresDialect, sql } from 'kysely';
import { Database } from './schema';
import { Pool } from 'pg';
import { writeFileSync } from 'fs';

let client: PgClient;

type Pagination = {};

export type GameFilters = Pagination & {
    maxPrice?: number;
};

type PgConfig = {
    user?: string;
    password?: string;
    database?: string;
    host?: string;
    port?: number;
    maxConnections?: number;
};

class PgClient {
    #db: Kysely<Database>;

    constructor(cfg: PgConfig) {
        const { user, password, database, host, port, maxConnections } = cfg;

        const dialect = new PostgresDialect({
            pool: new Pool({
                database,
                host,
                user,
                password,
                port,
                max: maxConnections,
            }),
        });

        this.#db = new Kysely<Database>({
            dialect,
        });
    }

    async close() {
        await this.#db.destroy();
    }

    async getGame({
        embedding,
        filters,
    }: {
        embedding?: number[];
        filters: GameFilters;
    }) {
        let query = this.#db.selectFrom('games');

        if (filters.maxPrice) {
            query = query.where('price', '<=', filters.maxPrice);
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

        return query.selectAll().execute();
    }
}

export function getPgClient(): PgClient {
    if (client) {
        console.log('cached pg');
        return client;
    }

    console.log('new pg');

    client = new PgClient({
        user: process.env.POSTGRES_USER,
        password: process.env.POSTGRES_PASSWORD,
        database: process.env.POSTGRES_DATABASE,
        host: process.env.POSTGRES_HOST,
        port: Number(process.env.POSTGRES_PORT ?? 5432),
        maxConnections: Number(process.env.POSTGRES_MAX_CONNECTIONS ?? 10),
    });

    process.on('SIGTERM', client.close.bind(client));
    process.on('SIGINT', client.close.bind(client));

    return client;
}
