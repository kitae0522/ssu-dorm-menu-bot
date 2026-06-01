package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/kitae0522/ssu-dorm-menu-bot/internal/menu"
	"github.com/kitae0522/ssu-dorm-menu-bot/internal/message"
	"github.com/kitae0522/ssu-dorm-menu-bot/internal/telegram"
)

const defaultMenuPath = "data/menus.json"

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	flags := flag.NewFlagSet("send-daily", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	input := flags.String("input", envOrDefault("MENU_JSON_PATH", defaultMenuPath), "path to menu JSON")
	date := flags.String("date", "", "target KST date in YYYY-MM-DD; defaults to today")
	dryRun := flags.Bool("dry-run", false, "print the message without calling Telegram")
	timeout := flags.Duration("timeout", 30*time.Second, "Telegram request timeout")
	baseURL := flags.String("telegram-base-url", telegram.DefaultBaseURL, "Telegram API base URL")
	if err := flags.Parse(args); err != nil {
		return err
	}

	store, err := menu.Load(*input)
	if err != nil {
		return err
	}

	target := time.Now()
	if *date != "" {
		target, err = menu.ParseKSTDate(*date)
		if err != nil {
			return fmt.Errorf("parse -date: %w", err)
		}
	}

	day, ok := menu.FindDay(store, target)
	if !ok {
		return fmt.Errorf("menu for %s not found in %s", menu.KSTDate(target), *input)
	}
	text := message.DailyMenu(day)
	if *dryRun {
		fmt.Fprintln(out, text)
		return nil
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client := telegram.Client{
		HTTPClient: &http.Client{Timeout: *timeout},
		BaseURL:    *baseURL,
	}
	if err := client.SendMessage(ctx, token, chatID, text); err != nil {
		return err
	}
	fmt.Fprintf(out, "sent menu for %s\n", day.Date)
	return nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
