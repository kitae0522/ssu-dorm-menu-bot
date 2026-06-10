package ssudorm

import (
	"strings"
	"testing"
	"time"
)

const sampleFoodTable = `
<html>
  <body>
    <table class="boxstyle02">
      <thead>
        <tr><th>날짜</th><th>조식</th><th>중식</th><th>석식</th><th>야식</th></tr>
      </thead>
      <tbody>
        <tr>
          <th><a href="javascript:viewContent('2026-06-01');">2026-06-01 (월)</a></th>
          <td>미운영</td>
          <td>김치콩나물국<br />소불고기<br />요구르트</td>
          <td>토마토스파게티<br>마늘빵</td>
          <td class="end"></td>
        </tr>
        <tr>
          <th><a href="javascript:viewContent('2026-06-02');">2026-06-02 (화)</a></th>
          <td>미운영</td>
          <td>시래기된장국<br />닭갈비</td>
          <td>돼지고기김치찌개<br />계란말이</td>
          <td>컵라면</td>
        </tr>
      </tbody>
    </table>
  </body>
</html>`

func TestParseHTMLExtractsWeeklyMenuTable(t *testing.T) {
	fetchedAt := time.Date(2026, 6, 1, 3, 0, 0, 0, time.FixedZone("KST", 9*60*60))

	got, err := ParseHTML(strings.NewReader(sampleFoodTable), "https://example.com/source", fetchedAt)
	if err != nil {
		t.Fatalf("ParseHTML returned error: %v", err)
	}

	if got.SourceURL != "https://example.com/source" {
		t.Fatalf("SourceURL = %q", got.SourceURL)
	}
	if got.WeekStart != "2026-06-01" {
		t.Fatalf("WeekStart = %q, want 2026-06-01", got.WeekStart)
	}
	if len(got.Days) != 2 {
		t.Fatalf("parsed %d days, want 2", len(got.Days))
	}
	first := got.Days[0]
	if first.Date != "2026-06-01" || first.Weekday != "월" {
		t.Fatalf("first day metadata = %+v", first)
	}
	if strings.Join(first.Meals.Lunch, ",") != "김치콩나물국,소불고기,요구르트" {
		t.Fatalf("first lunch = %#v", first.Meals.Lunch)
	}
	if strings.Join(first.Meals.Dinner, ",") != "토마토스파게티,마늘빵" {
		t.Fatalf("first dinner = %#v", first.Meals.Dinner)
	}
	if len(first.Meals.LateNight) != 0 {
		t.Fatalf("first late-night = %#v, want empty", first.Meals.LateNight)
	}
}

func TestParseHTMLRejectsWeekWithNoLunchOrDinner(t *testing.T) {
	fetchedAt := time.Date(2026, 6, 8, 3, 0, 0, 0, time.FixedZone("KST", 9*60*60))
	html := `
<html>
  <body>
    <table class="boxstyle02">
      <tbody>
        <tr>
          <th><a href="javascript:viewContent('2026-06-08');">2026-06-08 (월)</a></th>
          <td>미운영</td>
          <td></td>
          <td></td>
          <td class="end"></td>
        </tr>
        <tr>
          <th><a href="javascript:viewContent('2026-06-09');">2026-06-09 (화)</a></th>
          <td>미운영</td>
          <td></td>
          <td></td>
          <td class="end"></td>
        </tr>
      </tbody>
    </table>
  </body>
</html>`

	_, err := ParseHTML(strings.NewReader(html), "https://example.com/source", fetchedAt)
	if err == nil {
		t.Fatal("ParseHTML returned nil error, want no meal items error")
	}
	if !strings.Contains(err.Error(), "no lunch or dinner items found") {
		t.Fatalf("error = %q, want no lunch or dinner items found", err)
	}
}
