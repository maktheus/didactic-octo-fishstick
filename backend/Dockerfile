# syntax=docker/dockerfile:1
ARG GO_VERSION=1.24
ARG TARGETOS=linux
ARG TARGETARCH=amd64

FROM golang:${GO_VERSION}-bullseye AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src

# Copy Go module definition and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the API gateway binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -o /out/api ./cmd/api

FROM gcr.io/distroless/base-debian12 AS runner
WORKDIR /app
COPY --from=builder /out/api /app/api

ENV HTTP_PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/api"]
