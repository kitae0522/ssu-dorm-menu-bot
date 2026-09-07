package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kitae0522/ssu-dorm-menu-bot/internal/menu"
)

const sampleMenuPage = `
<html>
  <body>
    <table class="boxstyle02">
      <tbody>
        <tr>
          <th><a href="javascript:viewContent('2026-09-07');">2026-09-07 (월)</a></th>
          <td>미운영</td>
          <td>김치찌개<br />쌀밥</td>
          <td>된장찌개<br />계란말이</td>
          <td class="end"></td>
        </tr>
      </tbody>
    </table>
  </body>
</html>`

func TestRunSkipsWriteWhenMenuUnchanged(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(sampleMenuPage))
	}))
	t.Cleanup(server.Close)

	path := filepath.Join(t.TempDir(), "menus.json")
	existing := menu.Store{
		SourceURL: "https://example.com/old",
		FetchedAt: "2026-09-07T06:00:00+09:00",
		WeekStart: "2026-09-07",
		Days: []menu.Day{{
			Date:    "2026-09-07",
			Weekday: "월",
			Meals: menu.Meals{
				Breakfast: []string{"미운영"},
				Lunch:     []string{"김치찌개", "쌀밥"},
				Dinner:    []string{"된장찌개", "계란말이"},
				LateNight: nil,
			},
		}},
	}
	if err := menu.Save(path, existing); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	var out bytes.Buffer
	if err := run([]string{"-output", path, "-url", server.URL, "-timeout", "5s"}, &out); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !strings.Contains(out.String(), "no menu changes") {
		t.Fatalf("output = %q, want skip message", out.String())
	}

	got, err := menu.Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got.FetchedAt != existing.FetchedAt {
		t.Fatalf("FetchedAt = %q, want unchanged %q", got.FetchedAt, existing.FetchedAt)
	}
}

func TestRunWritesWhenMenuChanges(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(sampleMenuPage))
	}))
	t.Cleanup(server.Close)

	path := filepath.Join(t.TempDir(), "menus.json")
	existing := menu.Store{
		SourceURL: "https://example.com/old",
		FetchedAt: "2026-08-31T06:00:00+09:00",
		WeekStart: "2026-08-31",
		Days: []menu.Day{{
			Date:    "2026-08-31",
			Weekday: "월",
			Meals:   menu.Meals{Lunch: []string{"지난주메뉴"}},
		}},
	}
	if err := menu.Save(path, existing); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	var out bytes.Buffer
	if err := run([]string{"-output", path, "-url", server.URL, "-timeout", "5s"}, &out); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !strings.Contains(out.String(), "wrote") {
		t.Fatalf("output = %q, want write message", out.String())
	}

	got, err := menu.Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got.WeekStart != "2026-09-07" {
		t.Fatalf("WeekStart = %q, want 2026-09-07", got.WeekStart)
	}
	if got.FetchedAt == existing.FetchedAt {
		t.Fatal("expected fetched_at to be updated")
	}
	if _, err := time.Parse(time.RFC3339, got.FetchedAt); err != nil {
		t.Fatalf("fetched_at %q is not RFC3339: %v", got.FetchedAt, err)
	}
}
