ARG GO_VERSION=1.24
ARG ALPINE_VERSION=3.21

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base

WORKDIR /src

# Build
FROM base AS build

COPY --link go.mod go.sum ./
RUN go mod download

COPY --link pkg ./pkg
COPY --link services/orchestrator ./services/orchestrator

WORKDIR /src/services/orchestrator
RUN CGO_ENABLED=0 go build -o ./api -ldflags="-s -w" ./cmd/api/

# Run
FROM alpine:${ALPINE_VERSION}
RUN apk add --no-cache tini ca-certificates headscale cdrkit

ENV PATH=/app:${PATH}

RUN addgroup -g 1000 -S api \
    && adduser -u 1000 -S -G api api

COPY --chown=1000:1000 --from=build \
    /src/services/orchestrator/api \
    /app/

USER 1000

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["api"]
