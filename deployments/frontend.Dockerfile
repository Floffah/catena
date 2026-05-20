FROM oven/bun:1-alpine AS builder
WORKDIR /app

COPY . .
RUN bun install --frozen-lockfile

ARG CATENA_DIRECT_INSTANCE_URL
ARG CLERK_SECRET_KEY
ARG NEXT_PUBLIC_CATENA_INSTANCE_URL
ARG NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY

ENV CATENA_DIRECT_INSTANCE_URL=$CATENA_DIRECT_INSTANCE_URL
ENV CLERK_SECRET_KEY=$CLERK_SECRET_KEY
ENV NEXT_PUBLIC_CATENA_INSTANCE_URL=$NEXT_PUBLIC_CATENA_INSTANCE_URL
ENV NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=$NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY
ENV NEXT_TELEMETRY_DISABLED=1

RUN bun --bun run --cwd web build

FROM oven/bun:1-alpine AS runner
WORKDIR /app/web

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV HOSTNAME=0.0.0.0
ENV PORT=8080

#COPY --from=builder /app/web/public ./public
COPY --from=builder --link /app/web/.next ./.next
COPY --from=builder --link /app/web/.next/static ./.next/standalone/web/.next/static

USER bun

EXPOSE 8080

CMD ["bun", "run", ".next/standalone/web/server.js"]
