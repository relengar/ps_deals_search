'use server';

import { getNatsClient } from '@/lib/connectors/nats';
import { GameFilters, Platform, getPgClient } from '@/lib/connectors/postgres';
import { Game } from '../../connectors/postgres/schema';
import { logger } from '@/lib/logger';

type OrderBy = 'price' | 'rating';
type Order = 'ASC' | 'DESC';

export type SearchGameParams = {
    term?: string;
    maxPrice?: number;
    orderBy?: OrderBy;
    order?: Order;
    useSemantic?: boolean;
    platforms?: Platform[];
};

interface PgClient {
    getGame(params: {
        embedding?: number[];
        filters?: GameFilters;
    }): Promise<Game[]>;
}

interface NatsClient {
    requestEmbedding(text: string): Promise<number[]>;
}

type Dependencies = {
    pg: PgClient;
    nats: NatsClient;
};

export async function searchGamesQuery(
    deps: Dependencies,
    params: SearchGameParams
) {
    const { pg, nats } = deps;
    const { term } = params;
    let embedding: number[] | undefined = undefined;

    if (term) {
        logger.info({ term }, 'Requesting embedding for term');
        embedding = await nats.requestEmbedding(term);
    }

    logger.info({ params }, 'Retrieving game from db');
    const games = await pg.getGame({ embedding, filters: params });

    return games;
}

export const searchGames = async (params: SearchGameParams) => {
    const pg = getPgClient();
    const nats = await getNatsClient();

    return searchGamesQuery({ pg, nats }, params);
};
