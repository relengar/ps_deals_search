import z from 'zod';

const platforms = ['PS4', 'PS5'] as const;
export const platformsSchema = z.enum(platforms);
export type Platform = z.infer<typeof platformsSchema>;
