# SSU Dorm Menu Bot

GitHub Actions based Go automation for Soongsil University dorm cafeteria menus.

## What It Does

- Every Monday at 06:00 KST, crawl the dorm cafeteria menu page, commit `data/menus.json`, and notify Telegram with crawl success or failure.
- Every day at 07:00 KST, read `data/menus.json` and send that date's lunch and dinner through Telegram.
- Run without any external server or always-on process.

Source URL format:

```text
https://ssudorm.ssu.ac.kr/SShostel/mall_main.php?viewform=B0001_foodboard_list&gyear=YYYY&gmonth=MM&gday=DD
```

## Local Commands

```bash
go test ./...
go run ./cmd/crawl-weekly
go run ./cmd/send-daily -date 2026-06-01 -dry-run
```

## Message Examples

```text
✅ 식단 크롤링 완료
⏰ timestamp: 2026/06/01 03:00:05 KST
📅 period: 2026/06/01 월요일 ~ 2026/06/07 일요일
🍱 days: 7
🗂 file: data/menus.json
```

```text
❌ 식단 크롤링 실패
⏰ timestamp: 2026/06/01 03:00:05 KST
🌐 source: https://ssudorm.ssu.ac.kr/...
🧾 reason: network timeout
```

```text
🍽️ 2026/06/01 월요일 기숙사 식단

🌤️ 중식
- 쌀밥
- 소불고기
- 김치

🌙 석식
- 토마토스파게티
- 요구르트
```

## GitHub Secrets

Set these repository secrets before enabling the daily send workflow:

- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_CHAT_ID`

## Schedule Notes

GitHub Actions cron uses UTC. Workflows start five minutes before the target KST time, then `scripts/wait-until-kst.sh` waits until the exact KST target.
