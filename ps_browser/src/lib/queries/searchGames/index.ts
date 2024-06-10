import { getNatsClient } from '@/lib/connectors/nats';
import { getPgClient } from '@/lib/connectors/postgres';
import { Game } from '../../connectors/postgres/schema';

type OrderBy = 'price' | 'rating';
type Order = 'ASC' | 'DESC';

type SearchGameParams = {
    term?: string;
    maxPrice?: number;
    minPrice?: number;
    orderBy?: OrderBy;
    order?: Order;
};

interface PgClient {
    getGame(params: {
        embedding?: number[];
        filters?: Record<string, string | number>;
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
    let embedding: number[] | undefined = undefined;

    if (params.term) {
        embedding = await nats.requestEmbedding(params.term);
    }

    const games = await pg.getGame({ embedding, filters: { maxPrice: 20 } });

    return { games };
}

export const searchGames = async (params: SearchGameParams) => {
    const pg = getPgClient();
    const nats = await getNatsClient();

    return searchGamesQuery({ pg, nats }, params);
};
