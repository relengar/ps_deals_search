ARG NODE_VERSION=20
# Install dependencies only when needed
FROM node:${NODE_VERSION}-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json* ./
RUN apk add --no-cache libc6-compat &&\
    npm ci

# Rebuild the source code only when needed
FROM node:${NODE_VERSION}-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NEXT_TELEMETRY_DISABLED 1

RUN npm run build

# Production image, copy all the files and run next
FROM node:${NODE_VERSION}-alpine AS runner
WORKDIR /app

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

RUN adduser --system --uid 1001 browser

COPY --from=builder /app/public ./public

# Set the correct permission for prerender cache
RUN mkdir .next
RUN chown browser .next

# Automatically leverage output traces to reduce image size
# https://browser.org/docs/advanced-features/output-file-tracing
COPY --from=builder --chown=browser /app/.next/standalone ./
COPY --from=builder --chown=browser /app/.next/static ./.next/static

USER browser

EXPOSE ${PORT}

CMD HOSTNAME="0.0.0.0" node server.js
