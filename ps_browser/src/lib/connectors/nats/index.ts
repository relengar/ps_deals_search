import { logger } from '@/lib/logger';
import { NatsConnection, connect } from 'nats';

let client: NatsClient;

type EmbeddingResponse = {
    ok: boolean;
    embeddings: number[][];
};

export class NatsClient {
    #nc: NatsConnection | undefined;
    #getEmbeddingSubject: string;

    constructor(getEmbeddingSubject: string) {
        this.#getEmbeddingSubject = getEmbeddingSubject;
    }

    async connect({
        host,
        port,
        token,
    }: {
        host: string;
        port: number;
        token?: string;
    }) {
        this.#nc = await connect({ servers: [`${host}:${port}`], token });
    }

    async close() {
        if (!this.#nc) {
            return;
        }

        logger.info('Closing nats connection');
        await this.#nc.close();
    }

    async requestEmbedding(text: string): Promise<number[]> {
        const msg = await this.#nc!.request(
            this.#getEmbeddingSubject,
            JSON.stringify([text]),
            {
                timeout: 8 * 1000,
            }
        );
        const resp: EmbeddingResponse = msg.json();

        if (!resp.ok) {
            throw new Error('Embedder service failed to create embedding');
        }

        return resp.embeddings[0];
    }
}

export async function getNatsClient(): Promise<NatsClient> {
    if ((global as any).client) {
        return client;
    }

    const host = process.env.NATS_HOST ?? 'localhost';
    const port = Number(process.env.NATS_PORT ?? 4222);
    const token = process.env.NATS_TOKEN;
    const getEmbeddingsSubject =
        process.env.NATS_GET_EMBEDDINGS_SUBJECT ?? 'psgames.embedding';

    client = new NatsClient(getEmbeddingsSubject);

    await client.connect({ host, port, token });

    process.on('SIGTERM', client.close.bind(client));
    process.on('SIGINT', client.close.bind(client));

    return client;
}
