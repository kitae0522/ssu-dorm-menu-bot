package menu

import (
	"path/filepath"
	"testing"
	"time"
)

func TestFindDayUsesKSTDate(t *testing.T) {
	store := Store{
		Days: []Day{
			{Date: "2026-06-01", Weekday: "월"},
			{Date: "2026-06-02", Weekday: "화"},
		},
	}

	got, ok := FindDay(store, time.Date(2026, 5, 31, 15, 30, 0, 0, time.UTC))
	if !ok {
		t.Fatal("expected to find 2026-06-01 after converting UTC to KST")
	}
	if got.Date != "2026-06-01" {
		t.Fatalf("got date %q, want 2026-06-01", got.Date)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "menus.json")
	want := Store{
		SourceURL: "https://example.com/food",
		FetchedAt: "2026-06-01T03:00:00+09:00",
		WeekStart: "2026-06-01",
		Days: []Day{{
			Date:    "2026-06-01",
			Weekday: "월",
			Meals:   Meals{Lunch: []string{"김치콩나물국"}},
		}},
	}

	if err := Save(path, want); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got.SourceURL != want.SourceURL || got.WeekStart != want.WeekStart {
		t.Fatalf("loaded metadata mismatch: got %+v want %+v", got, want)
	}
	if got.Days[0].Meals.Lunch[0] != "김치콩나물국" {
		t.Fatalf("loaded lunch = %q", got.Days[0].Meals.Lunch[0])
	}
}
