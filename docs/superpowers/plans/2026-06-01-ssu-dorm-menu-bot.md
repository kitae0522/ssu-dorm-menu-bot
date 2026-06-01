# SSU Dorm Menu Bot Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go automation that crawls the SSU dorm cafeteria menu every Monday at 03:00 KST, commits `data/menus.json`, and sends the correct daily menu through Telegram at 06:00 KST.

**Architecture:** Two small commands call focused internal packages. `cmd/crawl-weekly` fetches and parses the SSU HTML page into a JSON store; `cmd/send-daily` loads that JSON, formats the KST date's menu, and posts it to Telegram. GitHub Actions supplies all scheduling and persistence; no server runs continuously.

**Tech Stack:** Go, `golang.org/x/net/html`, `golang.org/x/text`, GitHub Actions, Telegram Bot API.

---

### Task 1: Project Scaffold and Domain Tests

**Files:**
- Create: `go.mod`
- Create: `internal/menu/store_test.go`
- Create: `internal/ssudorm/parser_test.go`
- Create: `internal/telegram/client_test.go`

- [ ] Create Go module metadata.
- [ ] Write failing tests for menu lookup, message formatting, SSU table parsing, and Telegram payloads.
- [ ] Run `go test ./...` and confirm failures are caused by missing implementation packages.

### Task 2: Menu Store and Formatter

**Files:**
- Create: `internal/menu/store.go`

- [ ] Implement `Store`, `Day`, and `Meals` JSON types.
- [ ] Implement `Load`, `Save`, `FindDay`, `KSTNowDate`, and `FormatDailyMessage`.
- [ ] Run `go test ./internal/menu`.

### Task 3: SSU HTML Crawler and Parser

**Files:**
- Create: `internal/ssudorm/parser.go`

- [ ] Implement EUC-KR aware HTTP fetch with charset sniffing.
- [ ] Parse the `boxstyle02` table and extract date, weekday, breakfast, lunch, dinner, and late-night meals.
- [ ] Produce a `menu.Store` with `source_url`, `fetched_at`, `week_start`, and `days`.
- [ ] Run `go test ./internal/ssudorm`.

### Task 4: Telegram Sender

**Files:**
- Create: `internal/telegram/client.go`

- [ ] Implement Telegram `sendMessage` using form-encoded POST.
- [ ] Validate non-2xx responses and Telegram JSON `ok: false`.
- [ ] Run `go test ./internal/telegram`.

### Task 5: Commands and Configuration

**Files:**
- Create: `cmd/crawl-weekly/main.go`
- Create: `cmd/send-daily/main.go`
- Create: `.env.example`
- Create: `.gitignore`
- Create: `data/menus.json`

- [ ] Add CLI flags for paths, URL, timeouts, dates, and dry-run.
- [ ] Read `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`, and `MENU_JSON_PATH`.
- [ ] Run `go test ./...`.

### Task 6: GitHub Actions Workflows

**Files:**
- Create: `.github/workflows/crawl-weekly.yml`
- Create: `.github/workflows/send-daily.yml`
- Create: `scripts/wait-until-kst.sh`
- Create: `README.md`

- [ ] Add Monday 03:00 KST crawl workflow using UTC cron plus KST wait.
- [ ] Add daily 06:00 KST send workflow using UTC cron plus KST wait.
- [ ] Ensure crawl workflow commits and pushes `data/menus.json` only when changed.
- [ ] Document required secrets and manual dispatch.
- [ ] Run `go fmt ./...`, `go test ./...`, and `go vet ./...`.
