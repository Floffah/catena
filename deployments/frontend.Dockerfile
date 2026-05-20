FROM oven/bun:1-alpine AS builder
WORKDIR /app

COPY . .
RUN bun install --frozen-lockfile

ENV NEXT_TELEMETRY_DISABLED=1
RUN bun --bun run --cwd web build

FROM oven/bun:1-alpine AS runner
WORKDIR /app/web

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV HOSTNAME=0.0.0.0
ENV PORT=8080

#COPY --from=builder /app/web/public ./public
COPY --from=builder /app/web/.next ./.next
USER bun

EXPOSE 8080

CMD ["bun", "run", ".next/standalone/web/server.js"]