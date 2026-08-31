# Repository Guidelines

## Project Structure & Module Organization

This Go project runs only on GitHub Actions; do not add an external server or daemon. It crawls `https://ssudorm.ssu.ac.kr/SShostel/mall_main.php?viewform=B0001_foodboard_list&gyear=YYYY&gmonth=MM&gday=DD` every Monday at 06:00 KST, saves weekly menu JSON, commits and pushes `data/menus.json`, then sends the matching daily menu through Telegram at 07:00 KST. Keep executables under `cmd/`: `cmd/crawl-weekly` for refresh and `cmd/send-daily` for delivery. Put reusable code under `internal/`, such as `internal/ssudorm`, `internal/menu`, and `internal/telegram`. Keep fixtures under `testdata/`.

## Build, Test, and Development Commands

- `go test ./...`: run all unit tests.
- `go test -race ./...`: run race checks before release or scheduler changes.
- `go fmt ./...`: format Go code.
- `go vet ./...`: catch suspicious Go patterns.
- `go run ./cmd/crawl-weekly`: fetch weekly cafeteria menus into JSON.
- `go run ./cmd/send-daily`: send today's JSON-backed menu to Telegram.

## Coding Style & Naming Conventions

Use standard Go formatting and small packages. Prefer lowercase names such as `ssudorm`, `menu`, and `telegram`; avoid `utils`. Treat KST (`Asia/Seoul`) as the source of truth for all menu dates and scheduler targets. Use `context.Context` and explicit HTTP timeouts. Store JSON dates as `YYYY-MM-DD`.

## Testing Guidelines

Use Go's `testing` package, `httptest` for crawler and Telegram API behavior, and `testdata/` fixtures for sample cafeteria HTML or JSON. Do not call live Soongsil or Telegram endpoints in normal tests. Cover parsing, missing dates, Monday refresh logic, Korean message formatting, and retries.

## GitHub Actions Operations

Keep workflows under `.github/workflows/`. Run jobs on GitHub-hosted `ubuntu-latest`, not self-hosted lab runners. Target times are KST: weekly crawl Monday 06:00 and daily send 07:00. Because GitHub cron is UTC-only, use `55 20 * * 0` for the crawl and `55 21 * * *` for delivery, then wait in-job until the exact KST time. The crawl commits/pushes `data/menus.json` only when changed. Use `permissions: contents: write`, a bot commit identity, and `workflow_dispatch` for manual recovery.

## Commit & Pull Request Guidelines

Use short imperative commit subjects. Conventional Commit prefixes are fine: `feat: add weekly crawler`, `fix: handle empty menu day`.

Pull requests should include changed behavior, test results such as `go test ./...`, linked issues if any, and sample Telegram output when formatting changes. Mention workflow, cron, timezone, or env-var changes.

## Security & Configuration Tips

Never commit Telegram tokens, chat IDs, cookies, or production `.env` files. Store production values in GitHub Actions secrets, and document required variables in `.env.example`, including `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`, and `MENU_JSON_PATH`.
