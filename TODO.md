# TODO

## Zrobione (example v1)

- [x] API + auth + JWT + weryfikacja e-mail (MailHog)
- [x] Frontend Nuxt 4 — panel, strony, filtry użytkownika
- [x] Scraper Wikipedia (`wikipedia.pl`) + Docker prod
- [x] Browser pool, proxy per-site, retry 403/rate-limit (sync z projektu prywatnego)
- [x] HTML mail templates (`TemplateBuilder`)
- [x] API refactor: `Container`, Redis idempotency, JSON logging
- [x] Subskrypcje Basic/Premium + example payment + UI `/panel/user/subscribe`
- [x] Notifier worker (`cmd/notify`) + kontener prod + dev manual start
- [x] README + AGENTS.md z opisem pochodzenia (private repo → AI-assisted port)

## Do zrobienia

- [ ] Testy integracyjne pod subskrypcje, orders i notifier
- [ ] Uporządkowanie deploy docs (prod checklist: migracje, config, redis, secrets)
- [ ] Opcjonalnie: sync kolejnych poprawek infrastruktury z prywatnego repozytorium źródłowego

## Świadomie poza zakresem tego example

- Grafana / Loki / observability
- Prawdziwe bramki płatności (Stripe, PayU, …)
