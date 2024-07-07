// sum.test.js
import { expect, test, describe, beforeAll, afterEach, vi, Mock } from 'vitest';
import { searchGamesQuery } from '.';

describe('searchGamesQuery', () => {
    let deps: Parameters<typeof searchGamesQuery>[0];
    let getGames: Mock;
    let requestEmbedding: Mock;

    const mockEmbedding = [0.12321, 0.3123, 0.312321];
    const mockGame = { name: 'game' };

    beforeAll(async () => {
        getGames = vi.fn().mockReturnValue([mockGame]);
        requestEmbedding = vi.fn().mockReturnValue(mockEmbedding);
        deps = {
            gamesRepo: { getGames },
            nats: { requestEmbedding },
        };
    });

    afterEach(() => {
        getGames.mockClear();
        requestEmbedding.mockClear();
    });

    test('Search by a term', async () => {
        const term = 'something';
        const games = await searchGamesQuery(deps, { term });

        expect(requestEmbedding).toHaveBeenCalledWith(term);
        expect(getGames.mock.lastCall[0]).toHaveProperty(
            'embedding',
            mockEmbedding
        );

        expect(games).to.deep.include(mockGame);
    });

    test('Search without a term', async () => {
        const games = await searchGamesQuery(deps, {});

        expect(requestEmbedding).not.toHaveBeenCalled();
        expect(getGames).toHaveBeenCalled();
        expect(games).to.deep.include(mockGame);
    });
});
