# Repository Guidelines

## Project Structure & Module Organization
Back-End-Tcc is a Go microservice workspace. Keep runnable binaries in `cmd/` (e.g., `cmd/api`, `cmd/orchestrator`), shared infrastructure under `pkg/` (`config`, `logger`, `queue`, `observability`), and domain logic inside `services/<module>` pairing handlers, services, and repositories. Tests live in `tests/unit` for focused coverage, `tests/integration` for end-to-end flows, and `tests/fixtures` for reusable payloads. When adding a feature, align new packages with this layout so that runners, queues, and storage adapters remain reusable across services.

## Build, Test, and Development Commands
- `go mod tidy`: pins module dependencies before committing.
- `go run ./cmd/api`: boots the in-memory API gateway wiring every service.
- `go run ./cmd/<service>` (e.g., `go run ./cmd/runner`): runs a single microservice when isolating behavior.
- `go test ./...`: executes all unit and integration tests; use `APP_ENV=test` to mimic CI settings.
- `go test ./tests/integration -run E2E`: runs the benchmark flow regression referenced in `tests/integration/e2e_benchmark_flow_test.go`.

## Coding Style & Naming Conventions
Format Go code with `gofmt` (or `go fmt ./...` before pushing). Keep packages lower_snake_case, exported types in `CamelCase`, and unexported helpers in `camelCase`. Favor constructor functions (e.g., `NewQueue()`), context-aware method signatures, and structured logging via `pkg/logger`. Place HTTP routes under `pkg/http`, keep DTOs in `pkg/models`, and store configuration defaults in `pkg/config` so services stay thin.

## Testing Guidelines
Every new handler or worker needs unit tests alongside the package they touch plus scenario tests in `tests/unit`. Integration specs should focus on cross-service behavior via the in-memory queue; mirror existing patterns in `e2e_benchmark_flow_test.go`. Use descriptive test names like `TestOrchestrator_SubmitBatch`. When a change affects scoring or persistence, add fixtures to `tests/fixtures` and document any new environment flags referenced by the tests.

## Commit & Pull Request Guidelines
Commits in this repo follow short, imperative statements (`Add service observability instrumentation`, `feat/…`). Keep related changes squashed, include the surface area touched (service, package, or feature), and note testing commands in the message trailer if non-trivial. Pull requests should: describe the user-facing impact, list new configs or migrations, link any tracking issues, and attach logs or curl traces when touching HTTP flows. Highlight rollback considerations (queue messages, storage keys) before requesting review.

## Security & Configuration Tips
Configuration is environment-driven; document new vars in `README.md` and provide safe defaults in `pkg/config`. Never commit real secrets—use placeholder strings like `memory://default`. When running locally, export `JWT_SIGNING_SECRET` and `QUEUE_BUFFER_SIZE` via a `.env` the reviewer can reproduce.
