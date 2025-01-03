# syntax=docker/dockerfile:1

FROM golang:1.23-alpine

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /user/local/bin/api-calculator
# RUN go build -v -o /usr/local/bin/api-calculator ./...

CMD ["api-calculator"]
