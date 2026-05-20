FROM caddy:builder AS builder

ARG CADDY_PLUGINS=

COPY --chmod=755 deployments/proxy-build.sh /proxy-build.sh
RUN sh /proxy-build.sh

FROM caddy:alpine

COPY --from=builder /usr/local/bin/caddy-custom /usr/bin/caddy
COPY ./www /www
COPY deployments/Caddyfile /etc/caddy/Caddyfile