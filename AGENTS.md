# AGENTS.md

## Cursor Cloud specific instructions

This project is a Dockerized Wikipedia-scraper app: a Go API + scraper (`app/`), a
Nuxt 4 dashboard (`frontend/`), PostgreSQL, RabbitMQ, MailHog, Adminer, all behind an
Nginx reverse proxy on host port **8080**. Everything runs via Docker Compose; there is
no host-level Go/Node build in normal operation. See `README.md` for the product overview.

### Docker daemon
- Docker (with `fuse-overlayfs` storage driver and `iptables-legacy`) is installed in the
  VM image. There is no systemd, so `sudo service docker start` does not work.
- If `docker info` fails with a socket error, the daemon is not running. Start it in the
  background (e.g. in a tmux session): `sudo dockerd`. Wait a few seconds, then verify
  with `sudo docker info`. Use `sudo docker ...` (the `ubuntu` user is in the `docker`
  group, but the group only applies to newly-started login shells).

### Bringing up the stack (dev mode)
- Dev stack: `sudo docker compose -f compose.yaml -f compose.local.yaml up -d`
  (this is what `bin/run-local` runs, minus the sudo). The `--build` flag is only needed
  the first time or after Dockerfile changes.
- If `up` fails with a container-name conflict from a previous partial run, run
  `sudo docker compose -f compose.yaml -f compose.local.yaml down` first, then `up` again.
- Access points (all through Nginx on `http://localhost:8080`):
  - Dashboard: `/dashboard/`   API: `/api/`   MailHog: `/mailhog/`
  - Adminer: `/adminer/`   RabbitMQ mgmt: `/rabbitmq/`

### Backend (Go API) — non-obvious startup
- The `app` container does **not** auto-start the API; its command is `sleep infinity`.
  You must exec in and run it yourself. `air`, `dlv`, `migrate`, `make` are preinstalled.
- First-run inside the container (`sudo docker compose -f compose.yaml -f compose.local.yaml exec app sh`, then `cd /app`):
  1. `go mod tidy` (required — `go.sum` is gitignored, so it must be regenerated).
  2. `go run ./cmd/migrate/main.go` (runs GORM AutoMigrate; the API's own AutoMigrate is
     commented out, so migrations must be run separately).
  3. `air` to run the API with hot reload (listens on `:8080` in-container, proxied at `/api/`).
- The frontend container auto-runs `npm install && npm run dev` on startup (see
  `compose.local.yaml`); just wait for the Nuxt "Vite server warmed up" log.

### Config files (gitignored — required to start)
- `app/config.json` (copy from `app/config.example.json`) and `frontend/.env` (copy from
  `frontend/.env.example`) must exist or the Go apps `log.Fatal` and the frontend has no
  API base. The update script recreates them if missing.

### Auth / email flow
- New users must verify their email before they can log in (`Login` returns
  "Email not verified" otherwise). Verification emails are delivered to **MailHog**; open
  `http://localhost:8080/mailhog/` and use the `verify-email?token=...` link. This is the
  path to exercise register → verify → login end to end.

### Lint / test / build (per service)
- Backend build: `make api` inside the `app` container. NOTE: `make build` (and the
  production `docker/golang/.Dockerfile`) **fail** because `cmd/notify/main.go` does not
  exist yet (see `TODO.md`); the scraper/notify workers are also not wired into compose.
- Backend vet: `go vet ./...` (no golangci-lint configured).
- Backend tests: `go test ./internal/...`. Some unit tests currently fail on `main`
  (`internal/api`, `internal/service/browser`, `internal/service/scraper`) — these are
  pre-existing failures, not environment issues. Integration tests under
  `app/test/integration/` expect a separate `test_db` database and a RabbitMQ `test_user`
  with vhost `test` (see `config.test.json`), which are not provisioned by default.
- Frontend: there are **no** `lint`/`test` npm scripts. Run `npx eslint .`,
  `npx vitest run`, and `npm run build` directly inside the `frontend` container. Several
  vitest specs and eslint rules currently fail/error on `main` (pre-existing).
