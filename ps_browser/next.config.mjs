/** @type {import('next').NextConfig} */
const nextConfig = {
    experimental: {
        serverComponentsExternalPackages: ['pino'],
    },
    output: 'standalone',
};

export default nextConfig;
