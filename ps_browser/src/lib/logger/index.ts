import pino, { LoggerOptions } from 'pino';

const isDev = process.env.NODE_ENV === 'development';

const devEnvConfig: LoggerOptions = {
    transport: {
        target: 'pino-pretty',
        options: {
            colorize: true,
        },
    },
};

const config: LoggerOptions = {
    ...(isDev ? devEnvConfig : {}),
};

export const logger = pino(config);
