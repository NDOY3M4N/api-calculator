# syntax=docker/dockerfile:1

# ============================================
# 1st Stage: Build the application from source
# ============================================
FROM golang:1.23-alpine AS build
WORKDIR /app

# Install any runtime dependencies that are needed to run your application.
# Leverage a cache mount to /var/cache/apk/ to speed up subsequent builds.
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add gcc musl-dev

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Build the application
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux \
    go build -v \
    -ldflags '-linkmode external -extldflags -static' \
    -o /app/server

# =================================================
# 2nd Stage: Deploy app binary into a lean image
# =================================================
FROM gcr.io/distroless/base-debian11 AS final

WORKDIR /

COPY --from=build /app/server /server
COPY database.db database.db
COPY docs docs

EXPOSE 3000

ENTRYPOINT ["/server"]
