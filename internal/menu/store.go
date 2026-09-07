package menu

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const KSTLocation = "Asia/Seoul"

type Store struct {
	SourceURL string `json:"source_url"`
	FetchedAt string `json:"fetched_at"`
	WeekStart string `json:"week_start"`
	Days      []Day  `json:"days"`
}

type Day struct {
	Date    string `json:"date"`
	Weekday string `json:"weekday"`
	Meals   Meals  `json:"meals"`
}

type Meals struct {
	Breakfast []string `json:"breakfast"`
	Lunch     []string `json:"lunch"`
	Dinner    []string `json:"dinner"`
	LateNight []string `json:"late_night"`
}

func Load(path string) (Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Store{}, err
	}

	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		return Store{}, err
	}
	return store, nil
}

func Save(path string, store Store) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func SameMenu(a, b Store) bool {
	if a.WeekStart != b.WeekStart || len(a.Days) != len(b.Days) {
		return false
	}
	for i := range a.Days {
		if !sameDay(a.Days[i], b.Days[i]) {
			return false
		}
	}
	return true
}

func sameDay(a, b Day) bool {
	return a.Date == b.Date && a.Weekday == b.Weekday &&
		sameItems(a.Meals.Breakfast, b.Meals.Breakfast) &&
		sameItems(a.Meals.Lunch, b.Meals.Lunch) &&
		sameItems(a.Meals.Dinner, b.Meals.Dinner) &&
		sameItems(a.Meals.LateNight, b.Meals.LateNight)
}

func sameItems(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func FindDay(store Store, at time.Time) (Day, bool) {
	date := KSTDate(at)
	for _, day := range store.Days {
		if day.Date == date {
			return day, true
		}
	}
	return Day{}, false
}

func KSTDate(at time.Time) string {
	loc, err := time.LoadLocation(KSTLocation)
	if err != nil {
		loc = time.FixedZone("KST", 9*60*60)
	}
	return at.In(loc).Format("2006-01-02")
}

func ParseKSTDate(date string) (time.Time, error) {
	loc, err := time.LoadLocation(KSTLocation)
	if err != nil {
		loc = time.FixedZone("KST", 9*60*60)
	}
	return time.ParseInLocation("2006-01-02", date, loc)
}
