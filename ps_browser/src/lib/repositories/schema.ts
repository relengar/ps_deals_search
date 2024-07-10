import z from 'zod';

const order = ['ASC', 'DESC'] as const;
export const orderSchema = z.enum(order);
export type Order = z.infer<typeof orderSchema>;

export const paginationSchema = z.object({
    limit: z.number(),
    page: z.number().max(1000),
});
export type Pagination = z.infer<typeof paginationSchema>;
