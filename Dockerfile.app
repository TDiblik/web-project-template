# ============================
# Stage 1: Build frontend
# ============================
FROM oven/bun:alpine AS fe-build

ENV NODE_ENV=production

# Install dependencies
RUN apk add --no-cache git

# Build shared API client
WORKDIR /build/shared/fe/api-client
COPY shared/fe/api-client/package.json shared/fe/api-client/bun.lock ./
RUN bun install --frozen-lockfile
COPY shared/fe/api-client/ ./
RUN bun run build

# Build frontend
WORKDIR /build/fe
COPY fe/package.json fe/bun.lock ./
RUN bun install --frozen-lockfile
COPY fe/ ./

# Replace GIT_TAG in env file
COPY .git/ ./
RUN GIT_TAG=$(git rev-parse --verify HEAD) && \
    sed -i -e "s/__GIT_TAG__/${GIT_TAG}/g" .env.production && \
    rm -rf .git/ && \
    cp .env.production .env

RUN bun run build

# ============================
# Stage 2: Build backend
# ============================
FROM golang:alpine AS api-build

WORKDIR /build/api

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ ./

ENV CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOFLAGS=-buildvcs=false
RUN go build -trimpath -ldflags="-s -w" -o /build/api-executable
RUN chmod +x /build/api-executable

# ============================
# Stage 3: Final runtime image
# ============================
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies only
RUN apk add --no-cache dumb-init curl

# Create non-root user
RUN addgroup -S app_perms && adduser -S -G app_perms app_perms
RUN mkdir -p /app/images && chown -R app_perms:app_perms /app/images

# Copy built artifacts
COPY --chown=app_perms:app_perms --from=api-build /build/api-executable ./api
COPY --chown=app_perms:app_perms --from=api-build /build/api/database/migrations ./database/migrations
COPY --chown=app_perms:app_perms --from=fe-build /build/fe/dist ./public

USER app_perms:app_perms

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:35231/api/health || exit 1
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./api"]
