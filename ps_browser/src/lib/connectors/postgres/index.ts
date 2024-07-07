import { logger } from '@/lib/logger';
import { Kysely, LogEvent, PostgresDialect } from 'kysely';
import { Pool } from 'pg';
import { Database } from './schema';

let client: PgClient;

type PgConfig = {
    user: string;
    password: string;
    database: string;
    host: string;
    port?: number;
    maxConnections?: number;
    logQueries: boolean;
};

export class PgClient {
    #db: Kysely<Database>;
    #logQueries: boolean;

    constructor(cfg: PgConfig) {
        const {
            user,
            password,
            database,
            host,
            port,
            maxConnections,
            logQueries,
        } = cfg;

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
            log: this.logHandler.bind(this),
        });

        this.#logQueries = logQueries;
    }

    logHandler(event: LogEvent) {
        if (event.level === 'error') {
            logger.error(event.error, 'Postgres error');
        }
        if (this.#logQueries) {
            logger.info(
                { query: event.query.sql, time: event.queryDurationMillis },
                'Postgres query'
            );
        }
    }

    async close() {
        logger.info('Closing postgres connection');
        await this.#db.destroy();
    }

    get db() {
        return this.#db;
    }
}

export function getPgClient(): PgClient {
    if (client) {
        return client;
    }

    logger.info('Initializing pg client');

    client = new PgClient({
        user: process.env.POSTGRES_USER ?? 'postgres',
        password: process.env.POSTGRES_PASSWORD ?? 'postgres',
        database: process.env.POSTGRES_DATABASE ?? 'postgres',
        host: process.env.POSTGRES_HOST ?? '0.0.0.0',
        port: Number(process.env.POSTGRES_PORT ?? 5432),
        maxConnections: Number(process.env.POSTGRES_MAX_CONNECTIONS ?? 10),
        logQueries: process.env.LOG_QUERIES === 'true',
    });

    process.on('SIGTERM', client.close.bind(client));
    process.on('SIGINT', client.close.bind(client));

    return client;
}
