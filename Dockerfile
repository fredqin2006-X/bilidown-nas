FROM node:22-bookworm-slim AS frontend
WORKDIR /src/client
RUN npm install --global pnpm@9.15.4
COPY client/package.json client/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY client/ ./
RUN pnpm build

FROM golang:1.23-bookworm AS backend
WORKDIR /src/server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=0 go build -tags headless -trimpath -ldflags="-s -w" -o /out/bilidown .

FROM debian:bookworm-slim
ARG DEBIAN_MIRROR=http://deb.debian.org/debian
ARG DEBIAN_SECURITY_MIRROR=http://deb.debian.org/debian-security
RUN sed -i "s|http://deb.debian.org/debian-security|${DEBIAN_SECURITY_MIRROR}|g; s|http://deb.debian.org/debian|${DEBIAN_MIRROR}|g" /etc/apt/sources.list.d/debian.sources \
    && apt-get -o Acquire::ForceIPv4=true -o Acquire::http::Pipeline-Depth=0 -o Acquire::http::Timeout=30 -o Acquire::Retries=2 update \
    && apt-get -o Acquire::ForceIPv4=true -o Acquire::http::Timeout=30 -o Acquire::Retries=2 install -y --no-install-recommends ca-certificates ffmpeg \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /app /data /downloads && chown -R 1000:1000 /data /downloads
WORKDIR /app
COPY LICENSE NOTICE /app/
COPY --from=backend /out/bilidown /app/bilidown
COPY --from=frontend /src/server/static/ /app/static/
RUN find /app/static -type d -exec chmod 755 {} + \
    && find /app/static -type f -exec chmod 644 {} +
ENV BILIDOWN_SERVER=1 \
    BILIDOWN_LISTEN=:8098 \
    BILIDOWN_DB_PATH=/data/data.db \
    BILIDOWN_DOWNLOAD_DIR=/downloads \
    BILIDOWN_STATIC_DIR=/app/static
USER 1000:1000
EXPOSE 8098
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 CMD ["/app/bilidown", "--healthcheck"]
ENTRYPOINT ["/app/bilidown"]
