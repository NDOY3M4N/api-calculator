# syntax=docker/dockerfile:1

# ============================================
# 1st Stage: Build the application from source
# ============================================
FROM golang:1.23-alpine AS build-stage

# Install build tools
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

COPY . ./

# Install dependencies and build
RUN go mod download && go mod verify
RUN CGO_ENABLED=1 GOOS=linux go build -v -ldflags '-extldflags "-static"' -o /app/api-calculator

# =================================================
# 2nd Stage: Deploy app binary into a lean image
# =================================================
FROM gcr.io/distroless/base-debian11 AS build-release

WORKDIR /

COPY --from=build-stage /app/api-calculator /api-calculator
COPY --from=build-stage /app/docs ./docs

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/api-calculator"]
