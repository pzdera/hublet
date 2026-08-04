FROM node:22-alpine AS web-build

WORKDIR /src/web

COPY web/package.json ./
RUN npm install

COPY web/ ./
RUN npm run build


FROM golang:1.24-alpine AS go-build

WORKDIR /src

COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

COPY --from=web-build \
  /src/cmd/hublet-v2/web-dist \
  ./cmd/hublet-v2/web-dist

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/hublet-v2 \
    ./cmd/hublet-v2


FROM alpine:3.21

RUN addgroup -S -g 1000 hublet-v2 \
    && adduser -S -D -H -u 1000 -G hublet-v2 hublet-v2 \
    && apk add --no-cache \
      ca-certificates \
      tzdata \
      wget \
    && mkdir -p \
      /app/data \
      /app/icons \
      /app/wallpapers \
    && chown -R hublet-v2:hublet-v2 /app

COPY --from=go-build \
  /out/hublet-v2 \
  /usr/local/bin/hublet-v2

USER hublet-v2
WORKDIR /app

ENV HUBLET_V2_ADDR=:3000
ENV HUBLET_V2_DATA_DIR=/app/data

EXPOSE 3000

HEALTHCHECK \
  --interval=30s \
  --timeout=3s \
  --start-period=5s \
  CMD wget -q -O /dev/null \
    http://127.0.0.1:3000/api/v2/health \
    || exit 1

ENTRYPOINT ["/usr/local/bin/hublet-v2"]
