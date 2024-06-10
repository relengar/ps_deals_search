import tsconfigPaths from 'vite-tsconfig-paths';

export default {
    plugins: [tsconfigPaths()],
    resolve: {
        alias: {
            '@/': new URL('./src/app', import.meta.url).pathname,
        },
    },
};
