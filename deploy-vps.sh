#!/bin/sh
set -eu

binary=/home/juan/jsab/jsab-server
cert_dir=/opt/jsab-decouverte
secret_file=/opt/pac99-nextendo/nextendo_secret.key
nat_dir=/opt/mk8nex-nat

test -x "$binary"
test -r "$cert_dir/cert.pem"
test -e "$cert_dir/key.pem"
test -e "$secret_file"
test -r "$nat_dir/nat_endpoints.txt"
dash_token=$(docker inspect pac99nex --format '{{range .Config.Env}}{{println .}}{{end}}' | sed -n 's/^DASH_TOKEN=//p')
test -n "$dash_token"

if docker container inspect jsab >/dev/null 2>&1; then
    echo 'Container jsab already exists; inspect it before replacing it.' >&2
    exit 1
fi
if ss -ltnH '( sport = :60024 )' | grep -q .; then
    echo 'TCP port 60024 is in use.' >&2
    exit 1
fi

docker run -d \
    --name jsab \
    --network coolify \
    --memory 128m \
    --restart unless-stopped \
    -p 60024:60024 \
    -p 39000-39015:39000-39015/udp \
    -v /home/juan/jsab:/app:ro \
    -v "$cert_dir/cert.pem:/data/cert.pem:ro" \
    -v "$cert_dir/key.pem:/data/key.pem:ro" \
    -v "$secret_file:/data/nextendo_secret.key:ro" \
    -v "$nat_dir:/nat:ro" \
    -e AUTH_PORT=443 \
    -e SECURE_PORT=60024 \
    -e DASH_PORT=8113 \
    -e DASH_TOKEN="$dash_token" \
    -e NEXTENDO_PROXY_PROTOCOL=1 \
    -e NEXTENDO_HOST=g2a699600-lp1.s.n.srv.nintendo.net \
    -e NEXTENDO_REQUIRE_ACCOUNT=1 \
    -e NEXTENDO_REQUIRE_SIGNED_TOKEN=1 \
    -e NEXTENDO_SECRET_FILE=/data/nextendo_secret.key \
    -e NNCS_NAT_FILE=/nat/nat_endpoints.txt \
    -e CERT_FILE=/data/cert.pem \
    -e KEY_FILE=/data/key.pem \
    -l traefik.enable=true \
    -l traefik.tcp.routers.jsab.entrypoints=https \
    -l 'traefik.tcp.routers.jsab.rule=HostSNI(`g2a699600-lp1.s.n.srv.nintendo.net`) || HostSNI(`g2a699600-lp1.n.n.srv.nintendo.net`)' \
    -l traefik.tcp.routers.jsab.tls.passthrough=true \
    -l traefik.tcp.services.jsab.loadbalancer.proxyProtocol.version=1 \
    -l traefik.tcp.services.jsab.loadbalancer.server.port=443 \
    debian:stable-slim /app/jsab-server
