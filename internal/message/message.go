package message

import (
	"fmt"
	"strings"
	"time"

	"github.com/kitae0522/ssu-dorm-menu-bot/internal/menu"
)

func DailyMenu(day menu.Day) string {
	return fmt.Sprintf(
		"🍽️ %s %s 기숙사 식단\n\n🌤️ 중식\n%s\n\n🌙 석식\n%s",
		formatDate(day.Date),
		longWeekday(day.Weekday),
		bulletMeal(day.Meals.Lunch),
		bulletMeal(day.Meals.Dinner),
	)
}

func CrawlSuccess(at time.Time, store menu.Store, path string) string {
	return fmt.Sprintf(
		"✅ 식단 크롤링 완료\n⏰ timestamp: %s\n📅 period: %s\n🍱 days: %d\n🗂 file: %s",
		timestamp(at),
		period(store),
		len(store.Days),
		path,
	)
}

func CrawlFailure(at time.Time, sourceURL, reason string) string {
	return fmt.Sprintf("❌ 식단 크롤링 실패\n⏰ timestamp: %s\n🌐 source: %s\n🧾 reason: %s", timestamp(at), sourceURL, reason)
}

func formatDate(date string) string {
	parsed, err := menu.ParseKSTDate(date)
	if err != nil {
		return strings.ReplaceAll(date, "-", "/")
	}
	return parsed.Format("2006/01/02")
}

func timestamp(at time.Time) string {
	loc, err := time.LoadLocation(menu.KSTLocation)
	if err != nil {
		loc = time.FixedZone("KST", 9*60*60)
	}
	return at.In(loc).Format("2006/01/02 15:04:05 KST")
}

func longWeekday(short string) string {
	switch strings.TrimSpace(short) {
	case "월":
		return "월요일"
	case "화":
		return "화요일"
	case "수":
		return "수요일"
	case "목":
		return "목요일"
	case "금":
		return "금요일"
	case "토":
		return "토요일"
	case "일":
		return "일요일"
	default:
		return short
	}
}

func bulletMeal(items []string) string {
	var cleaned []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			cleaned = append(cleaned, item)
		}
	}
	if len(cleaned) == 0 {
		return "- 등록된 메뉴 없음"
	}
	for i, item := range cleaned {
		cleaned[i] = "- " + item
	}
	return strings.Join(cleaned, "\n")
}

func period(store menu.Store) string {
	if len(store.Days) == 0 {
		return "등록된 기간 없음"
	}
	first := store.Days[0]
	last := store.Days[len(store.Days)-1]
	return fmt.Sprintf("%s %s ~ %s %s", formatDate(first.Date), longWeekday(first.Weekday), formatDate(last.Date), longWeekday(last.Weekday))
}
