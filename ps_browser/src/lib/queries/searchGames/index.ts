'use server';

import { getNatsClient } from '@/lib/connectors/nats';
import { logger } from '@/lib/logger';
import {
    GameFilters,
    GameResponse,
    Platform,
    getGamesRepo,
} from '@/lib/repositories/games';

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

interface GamesRepo {
    getGames(params: {
        embedding?: number[];
        filters?: GameFilters;
    }): Promise<GameResponse[]>;
}

interface NatsClient {
    requestEmbedding(text: string): Promise<number[]>;
}

type Dependencies = {
    gamesRepo: GamesRepo;
    nats: NatsClient;
};

export async function searchGamesQuery(
    deps: Dependencies,
    params: SearchGameParams
) {
    const { gamesRepo, nats } = deps;
    const { term } = params;
    let embedding: number[] | undefined = undefined;

    if (term) {
        logger.info({ term }, 'Requesting embedding for term');
        embedding = await nats.requestEmbedding(term);
    }

    logger.info({ params }, 'Retrieving game from db');
    const games = await gamesRepo.getGames({ embedding, filters: params });

    return games;
}

export const searchGames = async (params: SearchGameParams) => {
    const gamesRepo = getGamesRepo();
    const nats = await getNatsClient();

    return searchGamesQuery({ gamesRepo, nats }, params);
};
