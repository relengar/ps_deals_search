import { getNatsClient } from '@/lib/connectors/nats';
import { getPgClient } from '@/lib/connectors/postgres';

type OrderBy = 'price' | 'rating';
type Order = 'ASC' | 'DESC';

type SearchGameParams = {
    term?: string;
    maxPrice?: number;
    minPrice?: number;
    orderBy?: OrderBy;
    order?: Order;
};

export const searchGames = async (params: SearchGameParams) => {
    const pg = getPgClient();
    let embedding: number[] | undefined = undefined;
    if (params.term) {
        const nats = await getNatsClient();
        embedding = await nats.requestEmbedding(params.term);
    }

    const games = await pg.getGame({ embedding, filters: { maxPrice: 20 } });

    // if (term) {
    //     redirect(`search?term=${term.toString()}`)
    // }

    return { games };
};
