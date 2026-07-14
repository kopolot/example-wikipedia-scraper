# AGENTS.md

Instrukcje dla agentów AI (Cursor Cloud, Codex itd.) pracujących nad tym repozytorium.

## Kontekst projektu

**Wiki Scraper Newsletter** (`example-wikipedia-scraper`) to **publiczny example** zbudowany na bazie **prywatnego** repozytorium autora o podobnej architekturze (newsletter-scraper: Docker, Go API, scraper, Nuxt). Ten fork:

- używa domeny **Wikipedia `Page`** i filtrów stron,
- zachowuje wspólną infrastrukturę (RabbitMQ, Redis, subskrypcje, example payment),
- **nie** zawiera Grafany/Loki ani produkcyjnych integracji płatności.

Właściciel synchronizował zmiany z projektu prywatnego z pomocą **Cursor AI**: diff między repozytoriami, kopiowanie modułów generycznych z zamianą ścieżek modułu Go, ręczna adaptacja modeli pod strony Wikipedia. Przy dalszym syncu portuj tylko infrastrukturę uniwersalną — nie kod z innej domeny biznesowej bez wyraźnej prośby.

Zobacz też [`README.md`](README.md) — sekcja „Pochodzenie kodu”.

---

## Cursor Cloud — środowisko

### Docker daemon
- Docker (z `fuse-overlayfs` i `iptables-legacy`) jest w obrazie VM. Brak systemd — `sudo service docker start` nie działa.
- Przy błędzie socketa: `sudo dockerd` w tle (np. tmux), potem `sudo docker info`.
- Używaj `sudo docker ...` jeśli grupa `docker` nie jest aktywna w shellu.

### Stack dev
```bash
sudo docker compose -f compose.yaml -f compose.local.yaml up -d
```
(jak `bin/run-local`, bez sudo lokalnie jeśli użytkownik ma grupę docker)

- Porty zajęte → `.env` z `.env.example` (`NGINX_HTTP_PORT`, `FRONTEND_DEV_PORT`, …).
- Konflikt nazw kontenerów → `docker compose ... down`, potem `up`.
- **Nginx** (domyślnie `:8080`): `/dashboard/`, `/api/`, `/mailhog/`, `/adminer/`, `/rabbitmq/`.

### Usługi w `compose.yaml`
`app`, `db`, `nginx`, `frontend`, `rabbitmq`, `mailhog`, **`redis`**. API **wymaga** Redis (idempotency) — `cmd/api/main.go` kończy się błędem bez połączenia.

---

## Backend (Go API) — startup dev

Kontenery `app` i `scraper`: **`sleep infinity`** — procesy uruchamiasz ręcznie.

W kontenerze `app` (`exec app sh`, `cd /app`):

1. `go mod tidy` — `go.sum` jest gitignored
2. `go run ./cmd/migrate/main.go` — migracje SQL **000001–000009** (nie GORM AutoMigrate w API)
3. `air` — API na `:8080`, proxowane jako `/api/`

Scraper: `exec scraper sh` → `go mod tidy` → `air -c .air.scraper.toml` lub `go run ./cmd/scraper/main.go`.

Frontend (compose.local): auto `npm install && npm run dev` — czekaj na log Vite.

---

## Pliki konfiguracyjne (gitignored)

| Plik | Źródło |
|------|--------|
| `app/config.json` | `app/config.example.json` |
| `frontend/.env` | `frontend/.env.example` |

**Wymagane w `config.json` dla pełnego API:**
- `redis` (host `redis`, port 6379)
- `payment_methods` — co najmniej `{ "name": "example", "enabled": true }`
- `rabbitmq`, `db`, `api`, `mailer`

Brak `config.json` → `log.Fatal` przy starcie.

---

## Architektura API (po refactorze)

- **`internal/api/container.go`** — DI repozytoriów i serwisów, `LoadModules()`
- **Moduły:** `UserApiModule`, `PageApiModule`, `UserWantedFiltersApiModule`, `OrderApiModule`
- **Idempotency:** `middleware/idempotent_middleware.go` + Redis (`internal/cache/`)
- **Routing:** `routing.go` (`joinRoutePath`), `response_writer.go` (redirecty 3xx)
- Frontend wysyła `X-Idempotent-Token` (`useApi.ts`, `uuid`)

### Endpointy subskrypcji / płatności
- `GET /user/subscription_levels` — plany + produkty (auth)
- `GET /user/subscription_levels/:level`
- `POST /order/` — utworzenie zamówienia + example payment (auth, idempotent)
- `GET /order/payment_methods`
- `PATCH /order/payment/:method/notify` — webhook PSP
- `GET /order/example_payment/:paymentId` — stub akceptacji

Kolejka `order_payment_notfied` → `SubscriptionService.AddSubscriptionTime`.

### Filtry stron
- Policy w `policy/http/user_wanted_pages_filters_policy.go` — **wymaga aktywnej subskrypcji** (Basic: 3 filtry, Premium: 10).
- Serwis: `PageFilterService`, repo: `UserWantedPagesFilterRepository`.

---

## Scraper

- **`BrowserPool`** + opcjonalne `proxy_url` per site w config
- Registry: tylko `wikipedia.pl` (+ stub `example`)
- JS: `app/resource/scraper/js/wikipedia.pl/`
- Prod: `compose.prod.yaml` — osobny serwis `scraper`, `xvfb`, mount `./app/resource`

---

## Auth / e-mail

- Login bez weryfikacji → `"Email not verified"`
- Maile w MailHog: `http://localhost:8080/mailhog/`
- Szablony HTML: `internal/service/mailer/template_builder.go`

Flow testowy: register → mailhog → verify-email → login → subscribe (example) → filtry.

---

## Lint / test / build

| Akcja | Komenda (w kontenerze `app`) |
|-------|------------------------------|
| Build API | `make api` |
| Build all | `make build` — **może failować** bez `cmd/notify/main.go` (TODO) |
| Vet | `go vet ./...` |
| Testy | `go test ./internal/...` — część testów browser/mailer wymaga sieci lub MailHog; znane flaky na `main` |

Frontend: brak skryptów `lint`/`test` w package.json — `npx eslint .`, `npx vitest run`, `npm run build` w kontenerze frontend.

---

## Czego NIE dodawać bez prośby

- Kod z **innej domeny biznesowej** prywatnego repozytorium źródłowego
- Grafana, Loki, observability stack
- Produkcja PSP (Stripe, PayU, …) — jest tylko **example payment**
- `cmd/notify` — planowany, jeszcze nie zaimplementowany

---

## Sync z projektem prywatnym

Przy aktualizacjach z repozytorium źródłowego autora:

1. Identyfikuj moduły **generyczne** (browser, cache, API, orders, mail) vs **domenowe**
2. Kopiuj generyczne z zamianą module path i nazw modeli Wikipedia
3. Uruchom `go build ./...` po każdej większej porcji
4. Osobne commity: infra / domena Wikipedia / frontend
5. Nie commituj `config.json`, `.env`, `go.sum` (gitignored)

Ostatni duży sync (2026-07): browser pool, proxy, mail templates, API container + Redis, subskrypcje + example payment.
