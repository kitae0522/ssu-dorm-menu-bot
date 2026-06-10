package ssudorm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kitae0522/ssu-dorm-menu-bot/internal/menu"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

const DefaultSourceURL = "https://ssudorm.ssu.ac.kr/SShostel/mall_main.php?viewform=B0001_foodboard_list&board_no=1"

var datePattern = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})\s*\(([^)]+)\)`)

func FetchAndParse(ctx context.Context, client *http.Client, sourceURL string, fetchedAt time.Time) (menu.Store, error) {
	if client == nil {
		client = http.DefaultClient
	}
	if sourceURL == "" {
		sourceURL = DefaultSourceURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return menu.Store{}, err
	}
	req.Header.Set("User-Agent", "ssu-dorm-menu-bot/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return menu.Store{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return menu.Store{}, fmt.Errorf("fetch menu page: unexpected status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return menu.Store{}, err
	}
	encoding, _, _ := charset.DetermineEncoding(body, resp.Header.Get("Content-Type"))
	decoded, err := io.ReadAll(encoding.NewDecoder().Reader(bytes.NewReader(body)))
	if err != nil {
		return menu.Store{}, err
	}

	return ParseHTML(bytes.NewReader(decoded), sourceURL, fetchedAt)
}

func ParseHTML(r io.Reader, sourceURL string, fetchedAt time.Time) (menu.Store, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return menu.Store{}, err
	}

	table := findFoodTable(doc)
	if table == nil {
		return menu.Store{}, fmt.Errorf("food table with class boxstyle02 not found")
	}

	var days []menu.Day
	for _, row := range rows(table) {
		cells := directCells(row)
		if len(cells) < 5 {
			continue
		}

		dateText := strings.Join(cellLines(cells[0]), " ")
		match := datePattern.FindStringSubmatch(dateText)
		if match == nil {
			continue
		}

		days = append(days, menu.Day{
			Date:    match[1],
			Weekday: match[2],
			Meals: menu.Meals{
				Breakfast: cellLines(cells[1]),
				Lunch:     cellLines(cells[2]),
				Dinner:    cellLines(cells[3]),
				LateNight: cellLines(cells[4]),
			},
		})
	}
	if len(days) == 0 {
		return menu.Store{}, fmt.Errorf("no menu rows found")
	}
	if !hasLunchOrDinnerItems(days) {
		return menu.Store{}, fmt.Errorf("no lunch or dinner items found in menu rows")
	}

	return menu.Store{
		SourceURL: sourceURL,
		FetchedAt: formatKST(fetchedAt),
		WeekStart: days[0].Date,
		Days:      days,
	}, nil
}

func hasLunchOrDinnerItems(days []menu.Day) bool {
	for _, day := range days {
		if len(day.Meals.Lunch) > 0 || len(day.Meals.Dinner) > 0 {
			return true
		}
	}
	return false
}

func findFoodTable(n *html.Node) *html.Node {
	if isElement(n, "table") && classContains(n, "boxstyle02") {
		return n
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if found := findFoodTable(child); found != nil {
			return found
		}
	}
	return nil
}

func rows(n *html.Node) []*html.Node {
	var out []*html.Node
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if isElement(node, "tr") {
			out = append(out, node)
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(n)
	return out
}

func directCells(row *html.Node) []*html.Node {
	var cells []*html.Node
	for child := row.FirstChild; child != nil; child = child.NextSibling {
		if isElement(child, "th") || isElement(child, "td") {
			cells = append(cells, child)
		}
	}
	return cells
}

func cellLines(cell *html.Node) []string {
	var lines []string
	var current strings.Builder

	var flush = func() {
		line := strings.Join(strings.Fields(current.String()), " ")
		current.Reset()
		if line != "" {
			lines = append(lines, line)
		}
	}

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			current.WriteString(node.Data)
			return
		}
		if isElement(node, "br") {
			flush()
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
		if isElement(node, "p") || isElement(node, "div") {
			flush()
		}
	}

	walk(cell)
	flush()
	if lines == nil {
		return []string{}
	}
	return lines
}

func isElement(n *html.Node, name string) bool {
	return n != nil && n.Type == html.ElementNode && strings.EqualFold(n.Data, name)
}

func classContains(n *html.Node, want string) bool {
	for _, attr := range n.Attr {
		if strings.EqualFold(attr.Key, "class") {
			for _, class := range strings.Fields(attr.Val) {
				if class == want {
					return true
				}
			}
		}
	}
	return false
}

func formatKST(t time.Time) string {
	loc, err := time.LoadLocation(menu.KSTLocation)
	if err != nil {
		loc = time.FixedZone("KST", 9*60*60)
	}
	return t.In(loc).Format(time.RFC3339)
}
