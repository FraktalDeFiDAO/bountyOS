# Repository Guidelines

## Project Structure & Module Organization
- `cmd/obsidian/main.go`: application entrypoint.
- `internal/core`: domain entities, interfaces, and scoring logic.
- `internal/adapters`: platform scanners, SQLite storage, and web adapter.
- `internal/security`, `internal/config`, `internal/notify`: cross-cutting security, config loading, and notifications.
- `config/config.yaml`: default runtime configuration.
- `web/`: Vue 3 frontend (`src/components`, `src/views`, `src/stores`, `src/router`).
- `e2e-tests/`: Playwright end-to-end tests.
- `scripts/`: maintenance utilities (Go formatting, security audit).

## Build, Test, and Development Commands
- `go mod tidy`: sync Go dependencies.
- `go build -o obsidian ./cmd/obsidian`: build backend binary.
- `./obsidian -no-ui`: run API/WebSocket server without terminal UI.
- `go test ./...`: run all Go unit tests.
- `BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP=1 go test ./...`: faster tests for rate-limited paths.
- `npx playwright test`: run E2E tests in `e2e-tests/` (expects backend on `http://localhost:12496`).
- `cd web && npm run dev`: run Vite dev server.
- `cd web && npm run build`: produce frontend assets for `web/dist`.
- `podman compose -f docker-compose.dev.yml up --build`: full local dev stack.

## Coding Style & Naming Conventions
- Go: format with `goimports` (use `./scripts/format-go.sh`), keep package names lowercase, use `PascalCase` for exported symbols.
- Vue/JS: use existing 2-space indentation, `PascalCase` component filenames (for example, `TopNav.vue`), and camelCase for store variables/functions.
- Keep scanner, storage, and UI logic inside their adapter modules; avoid leaking infra concerns into `internal/core`.

## Testing Guidelines
- Place Go tests beside code as `*_test.go`; name tests `TestXxx`.
- Prefer table-driven tests for scoring/parser logic.
- E2E specs live in `e2e-tests/*.spec.js`; keep user-facing flows deterministic and assert visible UI state.

## Commit & Pull Request Guidelines
- Follow Conventional Commit prefixes seen in history: `feat:`, `fix:`, `chore:`, `docs:`.
- Use imperative subjects and keep each commit focused.
- PRs should include: concise description, linked issue/task, test evidence (`go test ./...`, `npx playwright test` when relevant), and UI screenshots for frontend-visible changes.

## Security & Configuration Tips
- Never commit secrets; pass `GITHUB_TOKEN`/`DISCORD_WEBHOOK_URL` via environment variables.
- Keep operational settings in `config/config.yaml` and document any new keys in `README.md`.
