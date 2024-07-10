'use server';

import z from 'zod';

import { getNatsClient } from '@/lib/connectors/nats';
import { logger } from '@/lib/logger';
import {
    GameFilters,
    GameResponse,
    getGamesRepo,
} from '@/lib/repositories/games';
import { platformsSchema } from '@/lib/repositories/games/schema';
import { orderSchema, paginationSchema } from '@/lib/repositories/schema';

const orderBy = ['price', 'rating'] as const;
const orderBySchema = z.enum(orderBy);

const searchParamsSchema = z.object({
    term: z.string().optional(),
    maxPrice: z.number().optional(),
    orderBy: orderBySchema.optional(),
    order: orderSchema.optional(),
    useSemantic: z.boolean().optional(),
    platforms: platformsSchema.array().optional(),
});
export type SearchGameParams = z.infer<typeof searchParamsSchema>;

const allSearchParamsSchema = searchParamsSchema.merge(paginationSchema);
type SearchParams = z.infer<typeof allSearchParamsSchema>;

interface GamesRepo {
    getGames(params: {
        embedding?: number[];
        filters?: GameFilters;
    }): Promise<GameResponse[]>;
    countGames(params: {
        embedding?: number[];
        filters?: GameFilters;
    }): Promise<number>;
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
    params: SearchParams
) {
    const { gamesRepo, nats } = deps;
    const { term } = params;
    let embedding: number[] | undefined = undefined;

    if (term) {
        logger.info({ term }, 'Requesting embedding for term');
        embedding = await nats.requestEmbedding(term);
    }

    const filters = allSearchParamsSchema.parse(params);
    logger.info({ params, filters }, 'Retrieving game from db');
    const searchParams = { embedding, filters };
    const [games, total] = await Promise.all([
        gamesRepo.getGames(searchParams),
        gamesRepo.countGames(searchParams),
    ]);

    return { games, total };
}

export const searchGames = async (params: SearchParams) => {
    const gamesRepo = getGamesRepo();
    const nats = await getNatsClient();

    return searchGamesQuery({ gamesRepo, nats }, params);
};
