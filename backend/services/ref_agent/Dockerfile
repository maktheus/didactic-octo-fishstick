FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY pkg ./pkg
COPY services/ref_agent ./services/ref_agent

# Build the binary
RUN go build -o /ref-agent ./services/ref_agent

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /ref-agent .

EXPOSE 8080
CMD ["./ref-agent"]
