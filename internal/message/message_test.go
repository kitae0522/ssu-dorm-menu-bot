package message

import (
	"strings"
	"testing"
	"time"

	"github.com/kitae0522/ssu-dorm-menu-bot/internal/menu"
)

func TestDailyMenuUsesRichLunchDinnerFormat(t *testing.T) {
	day := menu.Day{
		Date:    "2026-06-01",
		Weekday: "월",
		Meals: menu.Meals{
			Lunch:  []string{"쌀밥", "소불고기", "김치"},
			Dinner: []string{"토마토스파게티", "요구르트"},
		},
	}

	got := DailyMenu(day)
	want := "🍽️ 2026/06/01 월요일 기숙사 식단\n\n🌤️ 중식\n- 쌀밥\n- 소불고기\n- 김치\n\n🌙 석식\n- 토마토스파게티\n- 요구르트"
	if got != want {
		t.Fatalf("DailyMenu() = %q, want %q", got, want)
	}
}

func TestDailyMenuUsesEmptyFallback(t *testing.T) {
	day := menu.Day{Date: "2026-06-01", Weekday: "월"}

	got := DailyMenu(day)
	if !strings.Contains(got, "🌤️ 중식\n- 등록된 메뉴 없음") {
		t.Fatalf("missing lunch fallback: %q", got)
	}
	if !strings.Contains(got, "🌙 석식\n- 등록된 메뉴 없음") {
		t.Fatalf("missing dinner fallback: %q", got)
	}
}

func TestCrawlStatusMessages(t *testing.T) {
	at := time.Date(2026, 6, 1, 3, 0, 5, 0, time.FixedZone("KST", 9*60*60))
	store := menu.Store{
		WeekStart: "2026-06-01",
		Days: []menu.Day{
			{Date: "2026-06-01", Weekday: "월"},
			{Date: "2026-06-02", Weekday: "화"},
			{Date: "2026-06-03", Weekday: "수"},
			{Date: "2026-06-04", Weekday: "목"},
			{Date: "2026-06-05", Weekday: "금"},
			{Date: "2026-06-06", Weekday: "토"},
			{Date: "2026-06-07", Weekday: "일"},
		},
	}

	success := CrawlSuccess(at, store, "data/menus.json")
	if success != "✅ 식단 크롤링 완료\n⏰ timestamp: 2026/06/01 03:00:05 KST\n📅 period: 2026/06/01 월요일 ~ 2026/06/07 일요일\n🍱 days: 7\n🗂 file: data/menus.json" {
		t.Fatalf("success message = %q", success)
	}

	failure := CrawlFailure(at, "https://ssudorm.ssu.ac.kr/...", "network timeout")
	if failure != "❌ 식단 크롤링 실패\n⏰ timestamp: 2026/06/01 03:00:05 KST\n🌐 source: https://ssudorm.ssu.ac.kr/...\n🧾 reason: network timeout" {
		t.Fatalf("failure message = %q", failure)
	}
}
