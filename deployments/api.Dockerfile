# syntax=docker/dockerfile:1

FROM golang:1.26.3-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY build.go ./
COPY api ./api
COPY cmd ./cmd
COPY data ./data
COPY internal ./internal
COPY tern.conf ./

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG RAILWAY_GIT_COMMIT_SHA=unknown
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -X github.com/floffah/catena.Commit=$RAILWAY_GIT_COMMIT_SHA" -o /out/catena-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/tern github.com/jackc/tern/v2

FROM alpine:3.22

RUN apk add --no-cache ca-certificates git git-daemon su-exec \
    && addgroup -S catena \
    && adduser -S -G catena catena \
    && mkdir -p /var/lib/catena/git \
    && chown -R catena:catena /var/lib/catena

COPY --from=build /out/catena-api /usr/local/bin/catena-api
COPY --from=build /out/tern /usr/local/bin/tern
COPY --from=build /src/data/migrations /app/data/migrations
COPY --from=build /src/tern.conf /app/tern.conf
COPY deployments/api-entrypoint.sh /usr/local/bin/api-entrypoint
RUN chmod +x /usr/local/bin/api-entrypoint

WORKDIR /app

ENV ENVIRONMENT=production \
    PORT=8080 \
    CATENA_GIT_ROOT=/var/lib/catena/git

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/api-entrypoint"]
