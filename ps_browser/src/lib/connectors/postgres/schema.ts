import {Selectable} from 'kysely'

interface GamesTable {
    id: number;
    name: string;
    description: string;
    price: number;
    original_price: number;
    url: string;
    rating: number;
    rating_sum: number;
    expiration: Date;
}

export type Game = Selectable<GamesTable>

interface GameEcmbeddingsTable {
    id: number;
    game_id: number;
    property_name: string;
    embedding: number[];
}

export type GameEmbedding = Selectable<GameEcmbeddingsTable>;

export interface Database {
    games: GamesTable
    game_embeddings: GameEcmbeddingsTable
  }