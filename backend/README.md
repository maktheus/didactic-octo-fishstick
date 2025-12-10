# Back-End-Tcc

This repository provides a reference implementation of a benchmark orchestration backend written in Go. It wires multiple domain services (auth, agents, benchmarks, orchestrator, runner, scoring, trace, leaderboard) behind a single API gateway, backed by in-memory queue, storage, and observability components.

## Core capabilities

- **Authentication**: `POST /auth` issues tokens for seeded admin users stored under `services/auth`.
- **Agent registry**: `GET/POST /agents` allows registering workers that will submit benchmark results.
- **Benchmark catalog**: `GET/POST /benchmarks` lets admins maintain runnable scenarios.
- **Submission pipeline**: `POST /submissions` persists payloads and publishes jobs; `GET /submissions` lists queued work.
- **Runner & scoring**: `GET /results` exposes runner outputs, `GET /scores` aggregates scoring summaries with async workers consuming the queue.
- **Telemetry**: `/traces` records execution events and `/leaderboard` lists aggregated benchmark winners, all instrumented with `pkg/observability/metrics`.

## Prerequisites

- Go 1.22+ (for local development)
- Docker 24+ (for containerized runs)

## Project layout

```
├── cmd/              # Entrypoints for the API gateway and individual microservices
├── pkg/              # Shared packages (configuration, logging, queue, storage, models)
├── services/         # Domain modules broken into handlers, services, and repositories
├── tests/            # Unit and integration tests with fixtures
└── README.md
```

## Running locally (Go)

1. Install dependencies and vendor module metadata:

   ```bash
   go mod tidy
   ```

2. Execute all unit and integration tests:

   ```bash
   go test ./...
   ```

3. Start the API gateway, which wires together every in-memory service:

   ```bash
   go run ./cmd/api
   ```

   Override defaults with environment variables such as `HTTP_PORT=9090 go run ./cmd/api`.

4. You can also run a single service for targeted debugging, e.g.:

   ```bash
   go run ./cmd/orchestrator
   ```

Structured logging and in-memory metrics (`pkg/observability/metrics`) remain enabled in all entrypoints so you can inspect queue throughput and handler latencies during development.

## Running with Docker

1. Build the container image (override `GO_VERSION`, `TARGETOS`, or `TARGETARCH` if needed):

   ```bash
   docker build -t back-end-tcc .
   ```

2. Run the API gateway container, exposing the HTTP port and supplying config overrides as needed:

   ```bash
   docker run --rm -p 8080:8080 \
     -e APP_ENV=development \
     -e HTTP_PORT=8080 \
     -e JWT_SIGNING_SECRET=dev-secret \
     back-end-tcc
   ```

   The binary inside the container is compiled with `CGO_ENABLED=0`, so it has no external dependencies; logs and metrics are emitted to stdout/stderr, making it suitable for orchestration with Docker Compose or Kubernetes.

## API exploration (Postman & Swagger)

- Postman: import `Back-End-Tcc.postman_collection.json` plus the companion environment `Back-End-Tcc.postman_environment.json`, then start the API (Go or Docker) and hit **Runner & Scores** → **Submissions** flow to see the in-memory pipeline in action.
- Swagger/OpenAPI: the API gateway now serves the spec at `GET /swagger.json`. To spin up Swagger UI locally, run:

  ```bash
  docker run --rm -p 8081:8080 \
    -e SWAGGER_JSON=/docs/swagger.json \
    -v "$(pwd)"/docs/openapi.json:/docs/swagger.json \
    swaggerapi/swagger-ui
  ```

  Then browse http://localhost:8081 to inspect and execute requests against your running gateway.

## Configuration

Configuration values are loaded from environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `APP_ENV` | `development` | Environment indicator logged by the services |
| `HTTP_PORT` | `8080` | Port used by each HTTP service |
| `QUEUE_BUFFER_SIZE` | `100` | Capacity hint for the in-memory queue |
| `STORAGE_DSN` | `memory://default` | Placeholder storage connection string |
| `JWT_SIGNING_SECRET` | `dev-secret` | Secret used for signing authentication tokens |

## Testing

Unit tests cover individual services such as orchestrator submission handling and scoring aggregation. Integration tests (`tests/integration/e2e_benchmark_flow_test.go`) exercise the full submission-to-scoring flow using the in-memory queue.

To run the full suite locally:
```bash
go test ./...
```

To run tests using Docker (useful if you don't have Go installed):
```bash
docker run --rm -v $(pwd):/src -v /var/run/docker.sock:/var/run/docker.sock -w /src golang:1.24-bullseye sh -c "go mod download && go test -v ./..."
```
Note: The E2E tests require access to the Docker daemon to spin up sandbox containers, so mounting `/var/run/docker.sock` is necessary when running tests inside a container.

## Continuous integration

GitHub Actions (`.github/workflows/ci.yml`) runs on every push/pull request to `main`. The pipeline enforces `gofmt`, `go vet`, `go test ./...`, a Linux build of `./cmd/api`, and a Docker image build to ensure the container stays healthy. Use `make test` and `make docker-build` locally to match the CI checks before opening a PR.

## Sample data

Fixtures located in `tests/fixtures/` provide sample payloads that can be used with tools such as `curl` or Postman to simulate submissions when experimenting with the API gateway.

## Next steps

This skeleton is intended to be extended with persistent storage backends, authentication tokens, real scoring logic, production-ready queue implementations, and real telemetry exporters (Prometheus/OpenTelemetry) connected to the existing logging and metrics hooks. The modular structure and clean interfaces make it straightforward to add these capabilities.
